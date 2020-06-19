package main

import (
	"fmt"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/io"
	"github.com/itfantasy/gonode/utils/yaml"
)

func main() {
	if nodeInfo, err := setupConfig(); err != nil {
		fmt.Println("Launch Error:" + err.Error())
	} else {
		Launch(nodeInfo)
	}
}

func setupConfig() (*gen_server.NodeInfo, error) {
	conf, err := io.LoadFile(io.CurrentDir() + "conf.yaml")
	if err != nil {
		return nil, err
	}
	info := gen_server.NewNodeInfo()
	if err := yaml.Unmarshal(conf, info); err != nil {
		return nil, err
	}
	return info, nil
}
