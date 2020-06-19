package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/itfantasy/gonode/utils/args"
	"github.com/itfantasy/gonode/utils/ts"
)

func run() error {
	parser := args.Parser().
		AddArg("proj", "proj-template", "set the project name of the runtime")

	projName, b := parser.Get("proj")
	if !b {
		return errors.New("the project name (-p) is necessary!")
	}
	return buildTheRunTime(projName)
}

func buildTheRunTime(projName string) error {

	fileName := projName
	fileName += "_" + ts.TimeToStr(ts.Time(), ts.FORMAT_NOW_C)
	fileName += ".so"

	projPath := "../" + projName + "/"

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", projPath+fileName, projPath+"runtime.go", projPath+"node.go")
	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			fmt.Println(err2)
			break
		}
		fmt.Println(line)
	}
	cmd.Wait()
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
	}
}
