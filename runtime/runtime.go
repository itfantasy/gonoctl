package main

import (
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
)

var node *GridNode

func OnHotUpdate() {
	gonode.Node().Bind(node)
}

func OnLaunch(proj string, nodeId string, nodeUrl string) {

	conf, err := ini.Load(proj + ".conf")
	if err != nil {
		fmt.Println("[OnLaunch]::" + err.Error())
		return
	}

	nodeInfo := gen_server.NewNodeInfo()

	if nodeId == "" {
		nodeInfo.Id = conf.Get("node", "id")
	} else {
		nodeInfo.Id = nodeId
	}
	if nodeUrl == "" {
		nodeInfo.Url = conf.Get("node", "url")
	} else {
		nodeInfo.Url = nodeUrl
	}

	nodeInfo.AutoDetect = conf.GetInt("node", "autodetect", 0) > 0
	nodeInfo.Public = conf.GetInt("node", "public", 0) > 0

	nodeInfo.RedUrl = conf.Get("redis", "url")
	nodeInfo.RedPool = conf.GetInt("redis", "pool", 0)
	nodeInfo.RedDB = conf.GetInt("redis", "db", 0)
	nodeInfo.RedAuth = conf.Get("redis", "auth")

	Launch(nodeInfo)
}

func Launch(nodeInfo *gen_server.NodeInfo) {
	node = NewGridNode(nodeInfo)
	gonode.Node().Initialize(node)
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
