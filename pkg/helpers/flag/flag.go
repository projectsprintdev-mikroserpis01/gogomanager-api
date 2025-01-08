package flag

import (
	"flag"
)

type Flag struct {
}

var FlagVars = getFlags()

func getFlags() *Flag {

	flag.Parse()

	flag := &Flag{}

	return flag
}
