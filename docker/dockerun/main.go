package main

import (
	//	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/itfantasy/gonode/utils/args"
)

func run() error {
	parser := configParser()
	return runGridDockerImage(parser)
}

func configParser() *args.ArgParser {
	parser := args.Parser().
		AddArg("d", "", "set the runtime dir that you will mount for the grid image").
		AddArg("p", "", "the ports that you will expose for the grid image, of course which must be compatible with the grid config").
		AddArg("l", "runtime_0.so", "set the latest runtime.so").
		AddArg("i", "", "dynamic set the id of the node").
		AddArg("u", "", "dynamic set the urls")
	return parser
}

func runGridDockerImage(parser *args.ArgParser) error {

	dir, exist := parser.Get("d")
	if !exist || dir == "" {
		return errors.New("the runtime dir (-d) is necessary!")
	}

	runtime, exist := parser.Get("l")
	nodeid, exist := parser.Get("i")
	nodeurl, exist := parser.Get("u")

	gridCmd := "etc/grid/grid -d=/etc/grid/runtime -l=" + runtime
	if nodeid != "" {
		gridCmd += " -i=" + nodeid
	}
	if nodeurl != "" {
		gridCmd += " -u=" + nodeurl
	}

	baseCmd := "docker run -v " + dir + ":/etc/grid/runtime "
	ports, exist := parser.Get("p")
	if exist && ports != "" {
		baseCmd += "-p " + ports + " "
	} else {
		baseCmd += "-P "
	}
	baseCmd += " itfantasy/grid " + gridCmd
	err := save(".sh", baseCmd)
	if err != nil {
		return err
	}
	fmt.Println(baseCmd)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("chmod", "777", ".sh")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err2 := cmd.Run()
	if err2 != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err2
	}

	cmd2 := exec.Command("bash", "./.sh")
	cmd2.Stdout = &out
	cmd2.Stderr = &stderr
	err3 := cmd2.Run()
	if err3 != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err3
	}

	return nil
}

func save(file string, content string) error {
	data := []byte(content)
	return ioutil.WriteFile(file, data, 0644)
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
