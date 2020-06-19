package main

import (
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/io"
	"github.com/itfantasy/gonode/utils/yaml"
)

var node *GridNode

func OnHotUpdate() {
	gonode.Bind(node)
}

func OnLaunch(proj string, namespace string, nodeid string, endpoints []string, etc string) {

	conf, err := io.LoadFile(proj + "conf.yaml")
	if err != nil {
		fmt.Println("[OnLaunch]::" + err.Error())
		return
	}
	info := gen_server.NewNodeInfo()
	if err := yaml.Unmarshal(conf, info); err != nil {
		fmt.Println("[OnLaunch]::" + err.Error())
		return
	}

	info.NameSpace = namespace
	info.NodeId = nodeid
	info.EndPoints = endpoints
	Launch(info)
}

func Launch(nodeInfo *gen_server.NodeInfo) {
	node = NewGridNode(nodeInfo)
	gonode.Bind(node)
	gonode.Launch()
}

func VersionName() string {
	return "1.0.0.0"
}

func VersionCode() int {
	return 1
}

func VersionInfo() string {
	return "info...."
}
