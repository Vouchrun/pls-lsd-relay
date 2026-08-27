package connectiontest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/config"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection"
)

// Startup tolerance: a single offline eth2 (beacon) endpoint must not prevent
// the relay from starting while other endpoints are usable — the runtime health
// checks already tolerate dead endpoints, so startup should too. Regression for
// the production incident where the dead g4mm4.io beacon endpoint (HTTP 522)
// made NewServiceManager fail and the relay exited before doing any work.
//
// Pre-fix behavior: NewConnection returned the first dead endpoint's error.
// Post-fix: the dead endpoint is skipped with a Warn log; NewConnection only
// fails when NO eth2 endpoint is usable, with an aggregated meaningful error.

const (
	testBeaconGenesisTime = uint64(1683785555) // PulseChain mainnet genesis time
	testBeaconSecPerSlot  = uint64(10)
	testBeaconSlotsPerEp  = uint64(32)
)

// deadEndpoint refuses connections instantly (nothing listens on port 1).
const deadEndpoint = "http://127.0.0.1:1"

func TestNewConnectionSkipsDeadEth2Endpoint(t *testing.T) {
	eth1 := mockEth1Server(t)
	defer eth1.Close()
	eth2 := mockEth2Server(t)
	defer eth2.Close()

	var logBuf bytes.Buffer
	logrus.SetOutput(&logBuf)
	defer logrus.SetOutput(os.Stderr)

	endpoints := []config.Endpoint{
		{Eth1: eth1.URL, Eth2: eth2.URL},
		{Eth1: eth1.URL, Eth2: deadEndpoint}, // offline beacon endpoint
	}
	conn, err := connection.NewConnection(endpoints, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("NewConnection must tolerate one dead eth2 endpoint, got: %v", err)
	}

	// the dead endpoint must be called out in a Warn log
	if !strings.Contains(logBuf.String(), deadEndpoint) {
		t.Fatalf("expected a startup Warn log naming the dead endpoint %s, logs: %s", deadEndpoint, logBuf.String())
	}

	// and the surviving endpoint must serve requests
	head, err := conn.GetBeaconHead()
	if err != nil {
		t.Fatalf("GetBeaconHead via surviving endpoint: %v", err)
	}
	if head.FinalizedEpoch == 0 {
		t.Fatal("expected non-zero finalized epoch from surviving endpoint")
	}
}

func TestNewConnectionFailsMeaningfullyWhenAllEth2EndpointsDead(t *testing.T) {
	eth1 := mockEth1Server(t)
	defer eth1.Close()

	endpoints := []config.Endpoint{
		{Eth1: eth1.URL, Eth2: deadEndpoint},
	}
	_, err := connection.NewConnection(endpoints, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected NewConnection to fail when no eth2 endpoint is usable")
	}
	if !strings.Contains(err.Error(), "no usable eth2 endpoint") {
		t.Fatalf("error should explain that no eth2 endpoint is usable (and name the offenders), got: %v", err)
	}
	if !strings.Contains(err.Error(), deadEndpoint) {
		t.Fatalf("error should name the failing endpoint %s, got: %v", deadEndpoint, err)
	}
}

// --- mocks ---

// mockEth1Server answers the two JSON-RPC calls NewConnection makes:
// eth_chainId (connect) and eth_getBlockByNumber (health check).
func mockEth1Server(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string          `json:"method"`
			ID     json.RawMessage `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch req.Method {
		case "eth_chainId":
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":"0x171"}`, req.ID)
		case "eth_getBlockByNumber":
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":%s}`, req.ID, mockBlockJSON())
		default:
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"error":{"code":-32601,"message":"the method %s does not exist/is not available"}}`, req.ID, req.Method)
		}
	}))
}

// mockEth2Server answers the Beacon REST API calls made at startup:
// /eth/v1/config/spec + /eth/v1/beacon/genesis (client construction) and
// /eth/v1/beacon/states/head/finality_checkpoints (health check), with a
// finalized epoch fresh enough that the client is not marked outOfSync.
func mockEth2Server(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/eth/v1/config/spec":
			fmt.Fprintf(w, `{"data":{"SECONDS_PER_SLOT":"%d","SLOTS_PER_EPOCH":"%d","EPOCHS_PER_SYNC_COMMITTEE_PERIOD":"256"}}`, testBeaconSecPerSlot, testBeaconSlotsPerEp)
		case "/eth/v1/beacon/genesis":
			fmt.Fprintf(w, `{"data":{"genesis_time":"%d","genesis_fork_version":"0x00000369","genesis_validators_root":"0x0000000000000000000000000000000000000000000000000000000000000000"}}`, testBeaconGenesisTime)
		case "/eth/v1/beacon/states/head/finality_checkpoints":
			secondsPerEpoch := testBeaconSecPerSlot * testBeaconSlotsPerEp
			currentEpoch := (uint64(time.Now().Unix()) - testBeaconGenesisTime) / secondsPerEpoch
			finalized := currentEpoch - 2 // realistic finality lag, fresh enough for the 20min outOfSync window
			epoch := strconv.FormatUint(finalized, 10)
			fmt.Fprintf(w, `{"data":{"previous_justified":{"epoch":%q},"current_justified":{"epoch":%q},"finalized":{"epoch":%q}}}`, epoch, epoch, epoch)
		default:
			http.Error(w, `{"code":404,"message":"not found"}`, http.StatusNotFound)
		}
	}))
}
