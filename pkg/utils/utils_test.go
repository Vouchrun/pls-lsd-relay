package utils_test

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ipfs/go-cid"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
	"github.com/web3-storage/go-w3s-client"
)

func TestAppendFile(t *testing.T) {
	path := "../../log_data/append_test2.txt"
	lastLine, err := utils.ReadLastLine(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lastLine)
	err = utils.AppendToFile(path, "\ntest1")
	if err != nil {
		t.Fatal(err)
	}
	err = utils.AppendToFile(path, "\ntest1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecodeInputData(t *testing.T) {
	client, err := ethclient.Dial("https://rpc.zhejiang.ethpandaops.io")
	if err != nil {
		t.Fatal(err)
	}
	type params struct {
		DealedHeight              *big.Int `json:"_dealedHeight"`
		UserAmount                *big.Int `json:"_userAmount"`
		NodeAmount                *big.Int `json:"_nodeAmount"`
		PlatformAmount            *big.Int `json:"_platformAmount"`
		MaxClaimableWithdrawIndex *big.Int `json:"_maxClaimableWithdrawIndex"`
	}

	abiBts, err := os.ReadFile("../../bindings/Withdraw/withdraw_abi.json")
	if err != nil {
		t.Log(err)
	}
	withdrawal, err := abi.JSON(bytes.NewReader(abiBts))
	if err != nil {
		t.Fatal(err)
	}
	hashs := []string{
		"0x68e297c767a3f65c44a97c464618b064a851b0c18b79c2da44d1f3723fdf35e3",
		"0x07a03487c382f31b34431ffb1ec8a5ee19998bb3a938ea62921d541edd18ec43",
		"0xa97d54212df21092bfc95375be4eee72cd148e514943f3f875a9eb48e9d0cc9b",
		"0x34490e517c0ce6345937b0c594165b1943296e64bfdec073006227b3a16319e3",
		"0x53bbcec2c9efdb8a0067d1336f21253a9874c8ee0c739a2f0461e8e29a6f8470",
		"0xfdf97b3502f4456f2712c29589639ea93233dd01736b8f40c0a1ee787d2867ff",
		"0x1ee51390948a7ee05510011112064864f1cc6fe8c2efd6ce43ec08b229f049d5",
		"0x8b99a8c3ee20d040576ba3f3bb403a4c9eae37198ad69a90176c2cdaad98b96c",
		"0x0aeefd268dfe4b601ff47adc7aa4d1f262391662e638c6bfdfdba8f6ffd3d323",
		"0xf6068990323225c704a87d059da0813ca2b9673b4fecf49ecf01d1f5aa36baef",
		"0x2a2b98b7203240e0289e5e8622b0adde445d9f5a6e526efc147010516831170b",
		"0x8ff158c187e933774471d8612faf96f2ee9cf12e244705adaf3ec700c029824e",
		"0xb4bc1dac092ccc9d3885f79823db28749346d0087c442980895f4c82a5c0e3c1",
		"0xc930066c98dad05903057cfaf7b3dd0eda3fbab47cf3257364754d5c4492e171",
		"0xab513904267fc3e78aa6fdcdd09e6e93051fadb3dd6b1733022a40997f3dc7b5",
		"0xa80f1fa39b67a0485bbc69a74309a70f7063ee845882bf20d63c9d1b3cc37366",
		"0x366ce8a336427d7609c231de4dc076e92419018626c28ff18aa70982a2a96a8f",
		"0x7ea63e6c41cc1ed2447c69e09f409b893d8680856ded4175533354a7164fe4eb",
		"0x236b5c8327b3328de76aa903fb8a2d7db41627f5ddf3e5c4eda2a74858b1fe4b",
		"0x7347fc4aeda8b5064234979f1271df6c33e832c535b2442bd6210704dbf9ac7d",
		"0x509d551269a0c9e85d38a7ec3add46aa21921f28fdf124de2a92db31684ef0ef",
		"0xa654a80cf78983a75e8f168df4d3686aaad960c30415d6574cb8cfbc25e6b348",
		"0x4683eb073ba24faf781ed85eefad74955c2d9dcc578d04b8194b4523aacff404",
		"0xf61d4b0d65336bf9c746a28e05a20fff55d545363b5d9459f54bda59a8def92b",
		"0x88705401927ff936229b153cd3516c8cad9e16173ed51b06f16876658babfc89",
		"0x7a43d8ac64d1435239400be138f96f5b476477941dfe40d52de3b9767f8722ea",
		"0x8554b6785c99f207eb4691da08327f158ed398e925f5e340b5178da238ae3002",
		"0xe32ecdbba7b37499aa9fa9a59989a0be20d0b96cc785eb81a273248a1269348c",
	}

	totalUserAmount := big.NewInt(0)
	totalNodeAmount := big.NewInt(0)
	totalPlatformAmount := big.NewInt(0)

	totalUserAmount1 := big.NewInt(0)
	totalNodeAmount1 := big.NewInt(0)
	totalPlatformAmount1 := big.NewInt(0)

	for _, hash := range hashs {

		tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(hash))
		if err != nil {
			t.Fatal(err)
		}
		receipient, err := client.TransactionReceipt(context.Background(), common.HexToHash(hash))
		if err != nil {
			t.Fatal(err)
		}

		t.Log("iiiiiii")
		for _, log := range receipient.Logs {
			t.Log("log ", log.Data)
		}
		method, err := withdrawal.MethodById(tx.Data()[:4])
		if err != nil {
			t.Fatal(err)
		}

		p := params{}

		inputMap := make(map[string]interface{}, 0)
		err = method.Inputs.UnpackIntoMap(inputMap, tx.Data()[4:])
		if err != nil {
			t.Log(err)
		}
		bts, _ := json.Marshal(inputMap)
		json.Unmarshal(bts, &p)
		total := new(big.Int).Add(new(big.Int).Add(p.UserAmount, p.NodeAmount), p.PlatformAmount)
		t.Log(p, total)
		totalUserAmount1 = new(big.Int).Add(totalUserAmount1, p.UserAmount)
		totalNodeAmount1 = new(big.Int).Add(totalNodeAmount1, p.NodeAmount)
		totalPlatformAmount1 = new(big.Int).Add(totalPlatformAmount1, p.PlatformAmount)

		// totalUserAmount = new(big.Int).Add(totalUserAmount, user.BigInt())
		// totalNodeAmount = new(big.Int).Add(totalNodeAmount, node.BigInt())
		// totalPlatformAmount = new(big.Int).Add(totalPlatformAmount, platform.BigInt())
	}
	t.Log(totalUserAmount, totalNodeAmount, totalPlatformAmount, new(big.Int).Add(totalNodeAmount, totalPlatformAmount))
	t.Log(totalUserAmount1, totalNodeAmount1, totalPlatformAmount1, new(big.Int).Add(totalNodeAmount1, totalPlatformAmount1))
	// {255680 4215416512500000 1328222037500000 291770450000000 3} 5835409000000000
	// utils_test.go:222: 4595384587500000 948253962500000 291770450000000
	// utils_test.go:220: {257920 3525303037500000 1477145212500000 263286750000000 3} 5265735000000000
	// utils_test.go:222: 4146766312500000 855681937500000 263286750000000
}

// 0x04df80 319360
// 0x039fc6d02bbbc0 1020101175000000

func TestGetGas(t *testing.T) {
	base, err := utils.GetGaspriceFromBeacon()
	if err != nil {
		t.Log(err)
	}
	t.Log(base)
	client, err := ethclient.Dial("https://mainnet-rpc.wetez.io/eth/v1/601083a01bf2f40729c5f75e62042208")
	if err != nil {
		t.Fatal(err)
	}
	gasTip, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		t.Log(err)
	}

	t.Log(gasTip)
}

func TestTx(t *testing.T) {
	client, err := ethclient.Dial("https://goerli.infura.io/v3/b3611f564322439ab2491e04ddd55b39")
	if err != nil {
		t.Fatal(err)
	}
	txr, err := client.TransactionReceipt(context.Background(), common.HexToHash("0x82c151b3d2ca12f59a3264b40fe60af85ef02680d5649e23802dcb532b937d45"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txr.Logs)

	tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash("0x82c151b3d2ca12f59a3264b40fe60af85ef02680d5649e23802dcb532b937d45"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tx)

}

func TestUploadFileToWeb3Storage(t *testing.T) {
	client, err := w3s.NewClient(w3s.WithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkaWQ6ZXRocjoweDNhZGRkOEUyZENjN2E3NDhmMzQ5NUYwYmNDQTUyYzgzQjVEOTQxMDYiLCJpc3MiOiJ3ZWIzLXN0b3JhZ2UiLCJpYXQiOjE2OTM3MDc2NzcyNTcsIm5hbWUiOiJ0ZXN0In0.2oOoiCvqjRGiVYNFH6UpRqVWtvE2NQQvzaJdjhC5jSQ"))
	if err != nil {
		t.Fatal(err)
	}

	// f, err := os.Open("/Users/tpkeeper/gowork/stafi/eth-lsd-relay/bindings/Erc20/erc20_abi.json")
	f, err := os.Open("/Users/tpkeeper/gowork/stafi/eth-lsd-relay/bindings/UserDeposit/userdeposit_abi.json")
	if err != nil {
		t.Fatal(err)
	}
	cid, err := client.Put(context.Background(), f)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cid.String())
}

func TestGetFileFromWeb3Storage(t *testing.T) {
	client, err := w3s.NewClient(w3s.WithToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkaWQ6ZXRocjoweDNhZGRkOEUyZENjN2E3NDhmMzQ5NUYwYmNDQTUyYzgzQjVEOTQxMDYiLCJpc3MiOiJ3ZWIzLXN0b3JhZ2UiLCJpYXQiOjE2OTM3MDc2NzcyNTcsIm5hbWUiOiJ0ZXN0In0.2oOoiCvqjRGiVYNFH6UpRqVWtvE2NQQvzaJdjhC5jSQ"))
	if err != nil {
		t.Fatal(err)
	}
	c, err := cid.Decode("bafybeia3y3ch5fykhmmesfiyfx3vn7jfzg4dz7lki7olwt2ydiksaqkqne")
	if err != nil {
		t.Fatal(err)
	}
	s, err := client.Status(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", s)

}

type NodeReward struct {
	Address                string          `json:"address"` // hex with 0x
	Index                  uint32          `json:"index"`
	TotalRewardAmount      decimal.Decimal `json:"totalRewardAmount"`
	TotalExitDepositAmount decimal.Decimal `json:"totalExitDepositAmount"`
	Proof                  string          `json:"proof"`
}

func TestMarshal(t *testing.T) {
	nodes := make([]*NodeReward, 0)
	bts, err := json.Marshal(nodes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bts))
}
