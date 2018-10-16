package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/itfantasy/gonode/utils/args"
)

func run() error {

	//return runGridDockerImage()
	return nil
}

func configParser() *args.ArgParser {
	parser := args.Parser().
		AddArg("d", "runtime", "set the runtime dir that you will mount for the grid image").
		AddArg("p", "", "the ports that you will expose for the grid image").
		AddArg("l", "runtime_0.so", "set the latest runtime.so").
		AddArg("i", "", "dynamic set the id of the node").
		AddArg("u", "", "dynamic set the urls")
	return parser
}

func runGridDockerImage() error {
	cmd := exec.Command("docker", "run", "-it", "itfantasy/grid", ".............")
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
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
