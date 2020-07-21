package main

import (
	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
)

var node *GridNode

func OnHotUpdate() {
	gonode.Bind(node)
}

func OnLaunch(proj string, namespace string, regdc string, nodeid string, endpoints []string, backends string, ispub int, etc string) {
	info := gen_server.NewNodeInfo()
	info.NameSpace = namespace
	info.RegDC = regdc
	info.NodeId = nodeid
	info.EndPoints = endpoints
	info.BackEnds = backends
	info.IsPub = ispub > 0
	Launch(info)
}

func Launch(nodeInfo *gen_server.NodeInfo) {
	node = NewGridNode(nodeInfo)
	gonode.Bind(node)
	gonode.Launch()
}

// @Meta.VersionName
func VersionName() string {
	return "1.0.0.0"
}

// @Meta.VersionCode
func VersionCode() int {
	return 1
}

// @Meta.VersionInfo
func VersionInfo() string {
	return "info...."
}
