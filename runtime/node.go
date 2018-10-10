package main

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
)

type GridNode struct {
	nodeInfo *gen_server.NodeInfo
}

func NewGridNode(nodeInfo *gen_server.NodeInfo) *GridNode {
	this := new(GridNode)
	this.nodeInfo = nodeInfo
	return this
}

func (this *GridNode) SelfInfo() (*gen_server.NodeInfo, error) {
	return this.nodeInfo, nil
}

func (this *GridNode) Start() {

}
func (this *GridNode) OnDetect(string) bool {
	return false
}
func (this *GridNode) OnConn(string) {

}
func (this *GridNode) OnMsg(string, []byte) {

}
func (this *GridNode) OnClose(string) {

}
func (this *GridNode) OnShell(string, string) {

}
func (this *GridNode) OnRanId() string {
	return ""
}
