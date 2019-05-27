package main

import (
	"fmt"

	//	"github.com/golang/protobuf/proto"
	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
)

type GridNode struct {
	nodeInfo *gen_server.NodeInfo
}

func NewGridNode(nodeInfo *gen_server.NodeInfo) *GridNode {
	g := new(GridNode)
	g.nodeInfo = nodeInfo
	return g
}

func (g *GridNode) Setup() *gen_server.NodeInfo {
	return g.nodeInfo
}
func (g *GridNode) Start() {

}
func (g *GridNode) OnConn(id string) {

}
func (g *GridNode) OnMsg(id string, msg []byte) {
	gonode.Send(id, msg)

}
func (g *GridNode) OnClose(id string) {
	fmt.Println("conn has been closed! " + id)
}
