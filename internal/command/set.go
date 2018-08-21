package command

import (
	"strings"

	"github.com/namreg/godown-v2/internal/storage"
)

//Set is the SET command
type Set struct {
	strg commandStorage
}

//Name implements Name of Command interface
func (c *Set) Name() string {
	return "SET"
}

//Help implements Help of Command interface
func (c *Set) Help() string {
	return `Usage: SET key value
Set key to hold the string value.
If key already holds a value, it is overwritten.`
}

//Execute implements Execute of Command interface
func (c *Set) Execute(args ...string) Result {
	if len(args) != 2 {
		return ErrResult{Value: ErrWrongArgsNumber}
	}

	value := strings.Join(args[1:], " ")

	setter := func(old *storage.Value) (*storage.Value, error) {
		return storage.NewString(value), nil
	}

	if err := c.strg.Put(storage.Key(args[0]), setter); err != nil {
		return ErrResult{Value: err}
	}
	return OkResult{}
}
