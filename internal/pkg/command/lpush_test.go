package command

import (
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/pkg/errors"

	"github.com/namreg/godown-v2/internal/pkg/storage"
	"github.com/namreg/godown-v2/internal/pkg/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestLpush_Name(t *testing.T) {
	cmd := new(Lpush)
	assert.Equal(t, "LPUSH", cmd.Name())
}

func TestLpush_Help(t *testing.T) {
	cmd := new(Lpush)
	expected := `Usage: LPUSH key value [value ...]
Prepend one or multiple values to a list.`
	assert.Equal(t, expected, cmd.Help())
}

func TestLpush_Execute(t *testing.T) {
	strg := memory.New(map[storage.Key]*storage.Value{
		"string": storage.NewStringValue("string"),
		"list":   storage.NewListValue([]string{"val1", "val2"}),
	})

	tests := []struct {
		name string
		args []string
		want Result
	}{
		{"ok", []string{"key", "field", "value"}, OkResult{}},
		{"wrong_type_op", []string{"string", "value"}, ErrResult{Value: ErrWrongTypeOp}},
		{"wrong_args_number/1", []string{}, ErrResult{Value: ErrWrongArgsNumber}},
		{"wrong_args_number/2", []string{"key"}, ErrResult{Value: ErrWrongArgsNumber}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Lpush{strg: strg}
			res := cmd.Execute(tt.args...)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestLpush_Execute_WhiteBox(t *testing.T) {
	expired := storage.NewListValue([]string{"val"})
	expired.SetTTL(time.Now().Add(-1 * time.Second))

	strg := memory.New(map[storage.Key]*storage.Value{
		"list":    storage.NewListValue([]string{"val1"}),
		"expired": expired,
	})
	tests := []struct {
		name   string
		args   []string
		verify func(t *testing.T, items map[storage.Key]*storage.Value)
	}{
		{
			"add_new_value_to_existing_key",
			[]string{"list", "val2"},
			func(t *testing.T, items map[storage.Key]*storage.Value) {
				val, ok := items["list"]
				assert.True(t, ok)

				expected := []string{"val2", "val1"}
				actual := val.Data().([]string)

				assert.Equal(t, expected, actual)
			},
		},
		{
			"add_new_value_to_not_existing_key",
			[]string{"list2", "val1"},
			func(t *testing.T, items map[storage.Key]*storage.Value) {
				val, ok := items["list2"]
				assert.True(t, ok)

				expected := []string{"val1"}
				actual := val.Data().([]string)

				assert.Equal(t, expected, actual)
			},
		},
		{
			"add_new_value_to_expired_key",
			[]string{"expired", "val2"},
			func(t *testing.T, items map[storage.Key]*storage.Value) {
				val, ok := items["expired"]
				assert.True(t, ok)

				expected := []string{"val2"}
				actual := val.Data().([]string)

				assert.Equal(t, expected, actual)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Lpush{strg: strg}
			res := cmd.Execute(tt.args...)
			assert.Equal(t, OkResult{}, res)

			items, err := strg.All()
			assert.NoError(t, err)

			tt.verify(t, items)
		})
	}
}

func TestLpush_Execute_StorageErr(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	err := errors.New("error")

	strg1 := NewStorageMock(t)
	strg1.GetMock.Return(nil, err)
	strg1.LockMock.Return()
	strg1.UnlockMock.Return()

	strg2 := NewStorageMock(t)
	strg2.GetMock.Return(storage.NewListValue([]string{"val"}), nil)
	strg2.PutMock.Return(err)
	strg2.LockMock.Return()
	strg2.UnlockMock.Return()

	cmd1 := Lpush{strg: strg1}
	cmd2 := Lpush{strg: strg2}

	res1 := cmd1.Execute("key", "val")
	assert.Equal(t, ErrResult{Value: err}, res1)

	res2 := cmd2.Execute("key", "val")
	assert.Equal(t, ErrResult{Value: err}, res2)
}
