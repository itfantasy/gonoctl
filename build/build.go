package build

import (
	"github.com/itfantasy/gonode/utils/args"
)

func Run() error {
	parser := args.Parser().
		AddArg("-p", "proj", "set the project name of the runtime").
		AddArg("-v", "0", "set the runtime version").
		AddArg("-f", "run.go", "select the gofile of the runtime entrance")
}

func buildTheRunTime(projName string, ver int) error {

}
