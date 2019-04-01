package main

import (
	"github.com/itfantasy/gonode"
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

func (this *GridNode) Setup() (*gen_server.NodeInfo, error) {
	return this.nodeInfo, nil
}
func (this *GridNode) Start() {

}
func (this *GridNode) OnConn(id string) {

}
func (this *GridNode) OnMsg(id string, msg []byte) {
	gonode.Send(id, msg)
}
func (this *GridNode) OnClose(id string) {

}
