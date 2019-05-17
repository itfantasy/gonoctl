package main

import (
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
)

var node *GridNode

func OnHotUpdate() {
	gonode.Bind(node)
}

func OnLaunch(proj string, nodeId string, nodeUrl string, pub bool) {

	conf, err := ini.Load(proj + ".conf")
	if err != nil {
		fmt.Println("[OnLaunch]::" + err.Error())
		return
	}

	nodeInfo := gen_server.NewNodeInfo()
	nodeInfo.Id = nodeId
	nodeInfo.Url = nodeUrl
	nodeInfo.Pub = pub
	nodeInfo.BackEnds = conf.Get("node", "backends")
	nodeInfo.LogLevel = conf.Get("log", "loglevel")
	nodeInfo.LogComp = conf.Get("log", "logcomp")
	nodeInfo.RegComp = conf.Get("reg", "regcomp")
	Launch(nodeInfo)
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
