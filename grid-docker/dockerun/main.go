package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"github.com/itfantasy/grid/utils/args"
)

func run() error {
	parser := configParser()
	return runGridDockerImage(parser)
}

func configParser() *args.ArgParser {
	parser := args.Parser().
		AddArg("proj", "aa", "the project dir which will mount for grid").
		AddArg("namespace", "", "set the namespace, such as 'itfantasy'").
		AddArg("nodeid", "", "set the nodeid, such as 'game_1024'").
		AddArg("endpoints", "", "set the endpoint list, such as 'tcp://yourserver:30005,kcp://yourserver:30006,ws://yourserver:30007/game_1024'").
		AddArg("etc", "", "extra configs")
	return parser
}

func runGridDockerImage(parser *args.ArgParser) error {

	proj, exist := parser.Get("proj")
	if !exist || proj == "" {
		return errors.New("project dir (-proj) is necessary!")
	}

	nodeid, exist := parser.Get("nodeid")
	if !exist || nodeid == "" {
		return errors.New("nodeid (-nodeid) is necessary!")
	}

	endpoints, exist := parser.Get("endpoints")
	if !exist || endpoints == "" {
		return errors.New("endpoint list (-endpoints) is necessary!")
	}

	namespace, _ := parser.Get("namespace")
	etc, _ := parser.Get("etc")

	gridCmd := "etc/grid/grid-core -proj=/etc/grid/runtime"
	gridCmd += " -nodeid=" + nodeid
	gridCmd += " -endpoints=" + endpoints
	if namespace != "" {
		gridCmd += " -namespace=" + namespace
	}
	if etc != "" {
		gridCmd += " -etc=" + etc
	}

	baseCmd := "docker run -v " + proj + ":/etc/grid/runtime "
	var endpointList []string = strings.Split(endpoints, ",")
	for _, endpoint := range endpointList {
		port, udpOrNot, err := extractPort(endpoint)
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
	if err := ioutil.WriteFile(".sh", []byte(baseCmd), 0644); err != nil {
		return err
	}
	fmt.Println(baseCmd)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("chmod", "777", ".sh")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	cmd2 := exec.Command("bash", "./.sh")
	cmd2.Stdout = &out
	cmd2.Stderr = &stderr
	if err := cmd2.Run(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
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
