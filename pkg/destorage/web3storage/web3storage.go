package web3storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-unixfsnode/data/builder"
	ipldCar "github.com/ipld/go-car/v2"
	"github.com/ipld/go-car/v2/blockstore"
	dagpb "github.com/ipld/go-codec-dagpb"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/destorage"
	ucanto_client "github.com/web3-storage/go-ucanto/client"
	"github.com/web3-storage/go-ucanto/core/car"
	ucanto_delegation "github.com/web3-storage/go-ucanto/core/delegation"
	"github.com/web3-storage/go-ucanto/did"
	"github.com/web3-storage/go-ucanto/principal"
	"github.com/web3-storage/go-ucanto/principal/ed25519/signer"
	"github.com/web3-storage/go-w3up/capability/storeadd"
	"github.com/web3-storage/go-w3up/capability/uploadadd"
	"github.com/web3-storage/go-w3up/client"
	"github.com/web3-storage/go-w3up/cmd/util"
	"github.com/web3-storage/go-w3up/delegation"
)

var _ destorage.DeStorage = &Storage{}
var fileUrlFormatter string = "https://w3s.link/ipfs/%s/%s"

type Storage struct {
	proofs []ucanto_delegation.Delegation
	space  did.DID
	signer principal.Signer
	conn   ucanto_client.Connection
}

func NewStorage(proofFilePath, spaceDid, privateKey string) (*Storage, error) {
	signer, err := signer.Parse(privateKey)
	if err != nil {
		return nil, fmt.Errorf("fail to parse private key: %w", err)
	}

	prfbytes, err := os.ReadFile(proofFilePath)
	if err != nil {
		return nil, fmt.Errorf("fail to read proof file: %w", err)
	}
	proof, err := delegation.ExtractProof(prfbytes)
	if err != nil {
		return nil, fmt.Errorf("fail to extract proof: %w", err)
	}

	space, err := did.Parse(spaceDid)
	if err != nil {
		return nil, fmt.Errorf("fail to parse space did: %w", err)
	}

	return &Storage{
		proofs: []ucanto_delegation.Delegation{proof},
		space:  space,
		signer: signer,
		conn:   util.MustGetConnection(),
	}, nil
}

func (c *Storage) DownloadFile(cid, fileName string) (content []byte, err error) {
	url := fmt.Sprintf(fileUrlFormatter, cid, fileName)
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rsp status err %d", rsp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if len(bodyBytes) == 0 {
		return nil, fmt.Errorf("bodyBytes zero err")
	}
	return bodyBytes, nil
}

func (c *Storage) UploadFile(content []byte, path string) (cid string, err error) {
	tempDir, err := os.MkdirTemp("", "ws3storage-upload-file")
	if err != nil {
		return "", fmt.Errorf("create temp dir error: %w", err)
	}

	filepath := joinPath(tempDir, filepath.Base(path))
	if err = os.WriteFile(filepath, content, 0600); err != nil {
		return "", fmt.Errorf("write to file error: %w", err)
	}
	defer os.RemoveAll(tempDir)

	carFileName := joinPath(tempDir, "carfile")
	_, err = c.createCar(context.TODO(), carFileName, filepath)
	if err != nil {
		return "", err
	}

	return c.uploadCar(carFileName)
}

func (c *Storage) uploadCar(carFile string) (string, error) {
	var shdlnks []ipld.Link
	carFs, err := os.Open(carFile)
	if err != nil {
		return "", fmt.Errorf("open carFile: %s err: %w", carFile, err)
	}
	defer carFs.Close()

	link, err := c.storeShard(carFs)
	if err != nil {
		return "", fmt.Errorf("store shard err: %w", err)
	}

	// fmt.Println(link.String())
	shdlnks = append(shdlnks, link)

	if _, err = carFs.Seek(0, 0); err != nil {
		return "", err
	}
	roots, _, err := car.Decode(carFs)
	if err != nil {
		return "", fmt.Errorf("reading roots: %w", err)
	}
	if len(roots) == 0 {
		return "", fmt.Errorf("roots is empty")
	}

	rcpt, err := client.UploadAdd(
		c.signer,
		c.space,
		&uploadadd.Caveat{
			Root:   roots[0],
			Shards: shdlnks,
		},
		client.WithProofs(c.proofs),
	)
	if err != nil {
		return "", fmt.Errorf("upload err: %w", err)
	}
	if rcpt.Out().Error() != nil {
		return "", fmt.Errorf("upload err: %s", rcpt.Out().Error().Message)
	}

	return roots[0].String(), nil
}

// createCar creates a car
func (c *Storage) createCar(ctx context.Context, carFile string, files ...string) (cid.Cid, error) {
	// make a cid with the right length that we eventually will patch with the root.
	hasher, err := multihash.GetHasher(multihash.SHA2_256)
	if err != nil {
		return cid.Undef, err
	}
	digest := hasher.Sum([]byte{})
	hash, err := multihash.Encode(digest, multihash.SHA2_256)
	if err != nil {
		return cid.Undef, err
	}
	proxyRoot := cid.NewCidV1(uint64(multicodec.DagPb), hash)

	options := []ipldCar.Option{blockstore.WriteAsCarV1(true)}

	cdest, err := blockstore.OpenReadWrite(carFile, []cid.Cid{proxyRoot}, options...)
	if err != nil {
		return cid.Undef, err
	}

	root, err := c.writeFiles(ctx, false, cdest, files...)
	if err != nil {
		return cid.Undef, err
	}

	if err := cdest.Finalize(); err != nil {
		return cid.Undef, err
	}
	// fmt.Println("cid: ", root.String())
	// re-open/finalize with the final root.
	if err = ipldCar.ReplaceRootsInFile(carFile, []cid.Cid{root}); err != nil {
		return cid.Undef, fmt.Errorf("ipldCar.ReplaceRootsInFile err: %w", err)
	}

	return root, nil
}

func (c *Storage) storeShard(shard io.Reader) (ipld.Link, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(shard)
	if err != nil {
		return nil, fmt.Errorf("reading CAR: %w", err)
	}

	mh, err := multihash.Sum(buf.Bytes(), multihash.SHA2_256, -1)
	if err != nil {
		return nil, fmt.Errorf("hashing CAR: %w", err)
	}

	link := cidlink.Link{Cid: cid.NewCidV1(0x0202, mh)}

	rcpt, err := client.StoreAdd(
		c.signer,
		c.space,
		&storeadd.Caveat{
			Link: link,
			Size: uint64(buf.Len()),
		},
		client.WithProofs(c.proofs),
	)

	if err != nil {
		return nil, fmt.Errorf("store/add %s: %w", link, err)
	}

	if rcptErr := rcpt.Out().Error(); rcptErr != nil {
		return nil, fmt.Errorf("rcpt out err: %s", rcptErr.Message)
	}

	if rcpt.Out().Ok().Status == "upload" {
		hr, err := http.NewRequest("PUT", *rcpt.Out().Ok().Url, bytes.NewReader(buf.Bytes()))
		if err != nil {
			return nil, fmt.Errorf("creating HTTP request: %w", err)
		}

		hdr := map[string][]string{}
		for k, v := range rcpt.Out().Ok().Headers.Values {
			hdr[k] = []string{v}
		}

		hr.Header = hdr
		hr.ContentLength = int64(buf.Len())
		httpClient := http.Client{}
		res, err := httpClient.Do(hr)
		if err != nil {
			return nil, fmt.Errorf("doing HTTP request: %w", err)
		}
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("non-200 status code while uploading file: %d", res.StatusCode)
		}
		err = res.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("closing request body: %w", err)
		}
	}

	return link, nil
}

// createCar creates a car
func (c *Storage) writeFiles(ctx context.Context, noWrap bool, bs *blockstore.ReadWrite, paths ...string) (cid.Cid, error) {
	if len(paths) == 0 {
		return cid.Undef, fmt.Errorf("you must specify at least one file")
	}

	ls := cidlink.DefaultLinkSystem()
	ls.TrustedStorage = true
	ls.StorageReadOpener = func(_ ipld.LinkContext, l ipld.Link) (io.Reader, error) {
		cl, ok := l.(cidlink.Link)
		if !ok {
			return nil, fmt.Errorf("not a cidlink")
		}
		blk, err := bs.Get(ctx, cl.Cid)
		if err != nil {
			return nil, err
		}
		return bytes.NewBuffer(blk.RawData()), nil
	}
	ls.StorageWriteOpener = func(_ ipld.LinkContext) (io.Writer, ipld.BlockWriteCommitter, error) {
		buf := bytes.NewBuffer(nil)
		return buf, func(l ipld.Link) error {
			cl, ok := l.(cidlink.Link)
			if !ok {
				return fmt.Errorf("not a cidlink")
			}
			blk, err := blocks.NewBlockWithCid(buf.Bytes(), cl.Cid)
			if err != nil {
				return err
			}
			bs.Put(ctx, blk)
			return nil
		}, nil
	}

	topLevel := make([]dagpb.PBLink, 0, len(paths))
	for _, p := range paths {
		l, size, err := builder.BuildUnixFSRecursive(p, &ls)
		if err != nil {
			return cid.Undef, err
		}
		if noWrap {
			rcl, ok := l.(cidlink.Link)
			if !ok {
				return cid.Undef, fmt.Errorf("could not interpret %s", l)
			}
			return rcl.Cid, nil
		}
		name := path.Base(p)
		entry, err := builder.BuildUnixFSDirectoryEntry(name, int64(size), l)
		if err != nil {
			return cid.Undef, err
		}
		topLevel = append(topLevel, entry)
	}

	// make a directory for the file(s).

	root, _, err := builder.BuildUnixFSDirectory(topLevel, &ls)
	if err != nil {
		return cid.Undef, nil
	}
	rcl, ok := root.(cidlink.Link)
	if !ok {
		return cid.Undef, fmt.Errorf("could not interpret %s", root)
	}

	return rcl.Cid, nil
}

func joinPath(dir, name string) string {
	if len(dir) > 0 && os.IsPathSeparator(dir[len(dir)-1]) {
		return dir + name
	}
	return dir + string(os.PathSeparator) + name
}
