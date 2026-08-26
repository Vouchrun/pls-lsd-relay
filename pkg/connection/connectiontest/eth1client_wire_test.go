// Package connectiontest holds wire-format tests for pkg/connection that only
// depend on pkg/connection itself. It lives in a subpackage so these tests stay
// hermetic and CGO-free: the sibling connection_test package imports the
// service package, which transitively requires CGO (prysm BLS deps), so tests
// in that package cannot compile on machines without a C toolchain.
package connectiontest

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection"
)

// Regression test for the production incident where debug_traceBlockByNumber
// was called with a raw *big.Int, which json.Marshal encodes as a DECIMAL
// string (e.g. "10394400"). Every geth-compatible RPC server rejects that with:
//
//	invalid argument 0: hex string without 0x prefix
//
// and the submitBalances/distributePriorityFee handlers then retry it forever.
// The block number must be hex-encoded (hexutil.EncodeBig) on the wire.
func TestDebugTraceBlockByNumberHexEncodesBlockNumber(t *testing.T) {
	var capturedParams json.RawMessage

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string          `json:"method"`
			Params json.RawMessage `json:"params"`
			ID     json.RawMessage `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		switch req.Method {
		case "eth_getBlockByNumber":
			// Health check in NewEth1Client: must return a valid block with a
			// recent timestamp or the client is marked unhealthy/outOfSync.
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":%s}`, req.ID, mockBlockJSON())
		case "debug_traceBlockByNumber":
			capturedParams = req.Params
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":[]}`, req.ID)
		default:
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"error":{"code":-32601,"message":"the method %s does not exist/is not available"}}`, req.ID, req.Method)
		}
	}))
	defer srv.Close()

	eth1Client, err := connection.NewEth1Client([]string{srv.URL})
	if err != nil {
		t.Fatal(err)
	}

	blockNumber := uint64(10394400)
	_, err = eth1Client.Debug_TraceBlockByNumber(context.Background(), new(big.Int).SetUint64(blockNumber), connection.Tracer{Tracer: "callTracer"})
	if err != nil {
		t.Fatal(err)
	}

	if len(capturedParams) == 0 {
		t.Fatal("mock server never received debug_traceBlockByNumber")
	}
	var params []json.RawMessage
	if err := json.Unmarshal(capturedParams, &params); err != nil {
		t.Fatalf("unmarshal params: %v", err)
	}
	if len(params) != 2 {
		t.Fatalf("expected 2 params (block number, tracer), got %d: %s", len(params), capturedParams)
	}

	var blockParam string
	if err := json.Unmarshal(params[0], &blockParam); err != nil {
		t.Fatalf("block number param must be a JSON string, got: %s", params[0])
	}
	if !strings.HasPrefix(blockParam, "0x") {
		t.Fatalf("block number param %q has no 0x prefix; RPC servers reject this with 'invalid argument 0: hex string without 0x prefix'", blockParam)
	}
	decoded, err := hexutil.DecodeUint64(blockParam)
	if err != nil {
		t.Fatalf("block number param %q is not valid hex: %v", blockParam, err)
	}
	if decoded != blockNumber {
		t.Fatalf("block number mismatch: got %d, want %d", decoded, blockNumber)
	}
}

// Live smoke test against a real execution endpoint (requires the debug
// namespace). Skipped unless ETH1_ENDPOINT is set, e.g.:
//
//	ETH1_ENDPOINT=https://rpc.vouch.run go test ./pkg/connection/connectiontest/ -run TestDebugTraceBlockByNumberLive -v
func TestDebugTraceBlockByNumberLive(t *testing.T) {
	endpoint := os.Getenv("ETH1_ENDPOINT")
	if endpoint == "" {
		t.Skip("ETH1_ENDPOINT not set")
	}

	eth1Client, err := connection.NewEth1Client([]string{endpoint})
	if err != nil {
		t.Fatal(err)
	}

	// A known PulseChain mainnet block with plenty of transactions.
	result, err := eth1Client.Debug_TraceBlockByNumber(context.Background(), big.NewInt(10394400), connection.Tracer{Tracer: "callTracer"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("expected non-empty trace result for block 10394400")
	}
	t.Logf("traced %d transactions", len(result))
}

func mockBlockJSON() string {
	return fmt.Sprintf(`{
		"number": "0x1a2b3c4",
		"hash": "0x0000000000000000000000000000000000000000000000000000000000000001",
		"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"miner": "0x0000000000000000000000000000000000000000",
		"stateRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"receiptsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"logsBloom": "0x%s",
		"difficulty": "0x0",
		"gasLimit": "0x2fefd80",
		"gasUsed": "0x0",
		"timestamp": "%s",
		"extraData": "0x",
		"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"nonce": "0x0000000000000000",
		"baseFeePerGas": "0x3b9aca00",
		"transactions": [],
		"uncles": []
	}`, strings.Repeat("00", 256), hexutil.EncodeUint64(uint64(time.Now().Unix())))
}
