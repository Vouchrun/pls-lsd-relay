package utils

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
)

func ErrToLogStr(err error) string {
	if err == nil {
		return ""
	}

	var rpcHttpErr rpc.HTTPError
	if errors.As(err, &rpcHttpErr) {
		// strip rpc response body from error messages
		str := err.Error()
		if len(rpcHttpErr.Body) > 0 {
			return strings.Replace(str, ": "+string(rpcHttpErr.Body), "", -1)
		}
	}

	return err.Error()
}
