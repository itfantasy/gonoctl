package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"io/ioutil"

	"github.com/itfantasy/grid/utils/args"
)

func run() error {
	parser := configParser()
	return runGridDockerImage(parser)
}

func configParser() *args.ArgParser {
	parser := args.Parser().
		AddArg("d", "", "set the runtime dir that you will mount for the grid image").
		AddArg("i", "", "dynamic set the id of the node").
		AddArg("l", "", "dynamic set the local url").
		AddArg("p", "n", "dynamic set if the node is public(n or y)")
	return parser
}

func runGridDockerImage(parser *args.ArgParser) error {

	dir, exist := parser.Get("d")
	if !exist || dir == "" {
		return errors.New("the runtime dir (-d) is necessary!")
	}

	nodeid, exist := parser.Get("i")
	if !exist || nodeid == "" {
		return errors.New("the runtime nodeid (-i) is necessary!")
	}

	nodeurl, exist := parser.Get("l")
	if !exist || nodeurl == "" {
		return errors.New("the runtime nodeurl (-l) is necessary!")
	}

	pub, _ := parser.Get("p")

	gridCmd := "etc/grid/grid -d=/etc/grid/runtime"
	gridCmd += " -i=" + nodeid
	gridCmd += " -l=" + nodeurl
	if puburl != "" {
		gridCmd += " -p=" + pub
	}

	baseCmd := "docker run -v " + dir + ":/etc/grid/runtime "

	port, udpOrNot, err := extractPort(nodeurl)
	if err != nil {
		return err
	}
	baseCmd += "-p " + port + ":" + port
	if udpOrNot != "" {
		baseCmd += "/" + udpOrNot
	}
	baseCmd += " "

	if puburl != "" {
		port, udpOrNot, err := extractPort(puburl)
		if err != nil {
			return err
		}
		baseCmd += "-p " + port + ":" + port
		if udpOrNot != "" {
			baseCmd += "/" + udpOrNot
		}
		baseCmd += " "
	}

	baseCmd += " itfantasy/grid " + gridCmd
	err1 := ioutil.WriteFile(".sh", []byte(baseCmd), 0644)
	if err1 != nil {
		return err1
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

func extractPort(url string) (string, string, error) {
	urlInfos := strings.Split(url, "://")
	if len(urlInfos) != 2 {
		return "", "", errors.New("illegal url!!" + url)
	}
	proto := urlInfos[0]
	tempInfos := strings.Split(urlInfos[1], "/")
	ipAndPort := strings.Split(tempInfos[0], ":")
	if len(ipAndPort) != 2 {
		return "", "", errors.New("illegal url!!" + url)
	}
	port, err := strconv.Atoi(ipAndPort[1])
	if err != nil {
		return "", "", errors.New("illegal port!!" + ipAndPort[1])
	}

	udpOrNot := ""
	if proto == "kcp" {
		udpOrNot = "udp"
	}

	return strconv.Itoa(port), udpOrNot, nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
