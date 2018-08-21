package command

import (
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/gojuno/minimock"

	"github.com/namreg/godown-v2/internal/storage"
	"github.com/namreg/godown-v2/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestHkeys_Name(t *testing.T) {
	cmd := new(Hkeys)
	assert.Equal(t, "HKEYS", cmd.Name())
}

func TestHkeys_Help(t *testing.T) {
	cmd := new(Hkeys)
	expected := `Usage: HKEYS key
Returns all field names in the hash stored at key. Order of fields is not guaranteed`
	assert.Equal(t, expected, cmd.Help())
}

func TestHkeys_Execute(t *testing.T) {
	expired := storage.NewMap(map[string]string{"field": "value", "field2": "value2"})
	expired.SetTTL(time.Now().Add(-1 * time.Second))

	strg := memory.New(map[storage.Key]*storage.Value{
		"string":       storage.NewString("value"),
		"hash":         storage.NewMap(map[string]string{"field": "value", "field2": "value2"}),
		"expired_hash": expired,
	})

	tests := []struct {
		name string
		args []string
		want Result
	}{
		{"ok", []string{"hash"}, SliceResult{Value: []string{"field", "field2"}}},
		{"not_existing_key", []string{"not_existing_key"}, NilResult{}},
		{"expired_key", []string{"expired_hash"}, NilResult{}},
		{"wront_type_op", []string{"string"}, ErrResult{Value: ErrWrongTypeOp}},
		{"wrong_number_of_args/1", []string{}, ErrResult{Value: ErrWrongArgsNumber}},
		{"wrong_number_of_args/2", []string{"key", "arg1"}, ErrResult{Value: ErrWrongArgsNumber}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Hkeys{strg: strg}
			res := cmd.Execute(tt.args...)
			if sr, ok := res.(SliceResult); ok {
				expected := tt.want.(SliceResult).Value
				sort.Strings(expected)

				actual := sr.Value
				sort.Strings(actual)

				assert.Equal(t, tt.want, res)
			} else {
				assert.Equal(t, tt.want, res)
			}
		})
	}
}

func TestHkeys_Execute_StorageErr(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	err := errors.New("error")

	strg := NewcommandStorageMock(mc)
	strg.GetMock.Return(nil, err)

	cmd := Hkeys{strg: strg}
	res := cmd.Execute("key")

	assert.Equal(t, ErrResult{Value: err}, res)
}
