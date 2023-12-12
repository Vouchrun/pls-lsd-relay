package local_store_test

import (
	"os"
	"testing"

	"github.com/stafiprotocol/eth-lsd-relay/pkg/local_store"
	"github.com/stretchr/testify/assert"
)

func TestReadAndUpdate(t *testing.T) {
	testFileName, err := os.CreateTemp(os.TempDir(), "eth-lsd-relay-")
	assert.Nil(t, err)
	defer os.Remove(testFileName.Name())
	s, err := local_store.NewLocalStore(testFileName.Name())
	assert.Nil(t, err)

	{
		val, err := s.Read("non-exist-key")
		assert.Nil(t, err)
		assert.Nil(t, val)
	}

	{
		addr := "0x179386303fC2B51c306Ae9D961C73Ea9a9EA0C8d"
		err := s.Update(local_store.Info{
			Address:      addr,
			SyncedHeight: 100,
		})
		assert.Nil(t, err)

		val, err := s.Read(addr)
		assert.Nil(t, err)
		assert.NotNil(t, val)
		assert.Equal(t, addr, val.Address)
		assert.Equal(t, uint64(100), val.SyncedHeight)

		err = s.Update(local_store.Info{
			Address:      addr,
			SyncedHeight: 299,
		})
		assert.Nil(t, err)

		val, err = s.Read(addr)
		assert.Nil(t, err)
		assert.NotNil(t, val)
		assert.Equal(t, addr, val.Address)
		assert.Equal(t, uint64(299), val.SyncedHeight)
	}
}
