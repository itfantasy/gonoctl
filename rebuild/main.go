package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"

	"github.com/itfantasy/gonode/utils/args"
)

func run() error {
	parser := args.Parser().
		AddArg("p", "runtime", "set the project name of the runtime").
		AddArg("v", "0", "set the runtime version")

	projName, b := parser.Get("p")
	if !b {
		return errors.New("the project name (-p) is necessary!")
	}
	ver, b := parser.GetInt("v")
	if !b {
		return errors.New("the runtime version (-v) is necessary!")
	}
	return buildTheRunTime(projName, ver)
}

func buildTheRunTime(projName string, ver int) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", projName+"_"+strconv.Itoa(ver)+".so", "runtime.go", "node.go")
	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err3 := cmd.Start()
	if err3 != nil {
		return err3
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
