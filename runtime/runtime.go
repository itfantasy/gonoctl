package main

import (
	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
)

var node *GridNode

func OnHotUpdate() {
	gonode.Node().Bind(node)
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
