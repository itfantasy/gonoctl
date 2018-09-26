package main

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
)

func main() {
	nodeInfo := gen_server.NewNodeInfo()
	OnHotUpdate()
	Launch(nodeInfo)
}

func OnHotUpdate() {

}

func Launch(nodeInfo *gen_server.NodeInfo) {

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
