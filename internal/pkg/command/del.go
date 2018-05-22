package command

import (
	"github.com/namreg/godown-v2/internal/pkg/storage"
)

func init() {
	commands["DEL"] = new(Del)
}

//Del is the DEL command
type Del struct{}

//Name implements Name of Command interface
func (c *Del) Name() string {
	return "DEL"
}

//Help implements Help of Command interface
func (c *Del) Help() string {
	return `Usage: DEL key
Del the given key.`
}

//ArgsNumber implements ArgsNumber of Command interface
func (c *Del) ArgsNumber() int {
	return 1
}

//Execute implements Execute of Command interface
func (c *Del) Execute(strg storage.Storage, args ...string) Result {
	err := strg.Del(storage.Key(args[0]))
	if err != nil {
		return ErrResult{err}
	}
	return OkResult{}
}