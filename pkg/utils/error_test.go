package utils_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

func Test_ErrToLogStr(t *testing.T) {
	{
		err := rpc.HTTPError{
			StatusCode: 520,
			Status:     "520",
			Body:       []byte("here is an unknown connection issue between Cloudflare and the origin web server. As a result, the web page can not be displayed"),
		}
		err1 := fmt.Errorf("call method err: %w", err)
		fmt.Println(utils.ErrToLogStr(err1))
	}

	{
		err := rpc.HTTPError{
			StatusCode: 520,
			Status:     "520",
			Body:       []byte("here is an unknown connection issue between Cloudflare and the origin web server. As a result, the web page can not be displayed"),
		}
		fmt.Println(utils.ErrToLogStr(errors.Join(fmt.Errorf("call method error"), err)))
	}

	{
		err := fmt.Errorf("fail to parse json")
		err1 := fmt.Errorf("call method err: %w", err)
		fmt.Println(utils.ErrToLogStr(err1))
	}
}
