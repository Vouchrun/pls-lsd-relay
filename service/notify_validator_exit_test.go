package service

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidatorSelection(t *testing.T) {
	for i := 0; i < 100; i++ {
		validatorsMap := map[string]*Validator{
			"0x011": {
				ValidatorIndex: 11,
				ActiveEpoch:    100,
			},
			"0x02": {
				ValidatorIndex: 2,
				ActiveEpoch:    100,
			},
			"0x03": {
				ValidatorIndex: 3,
				ActiveEpoch:    100,
			},
			"0x05": {
				ValidatorIndex: 5,
				ActiveEpoch:    901,
			},
		}
		vals := make([]*Validator, 0, len(validatorsMap))
		for _, val := range validatorsMap {
			vals = append(vals, val)
		}
		sort.SliceStable(vals, func(i, j int) bool {
			return vals[i].ActiveEpoch < vals[j].ActiveEpoch ||
				(vals[i].ActiveEpoch == vals[j].ActiveEpoch && vals[i].ValidatorIndex < vals[j].ValidatorIndex)
		})
		assert.Equal(t, uint64(2), vals[0].ValidatorIndex)
	}
}
