package main

import (
	"fmt"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"
)

func main() {
	if nodeInfo, err := setupConfig(); err != nil {
		fmt.Println("Launch Error:" + err.Error())
	} else {
		Launch(nodeInfo)
	}
}

func setupConfig() (*gen_server.NodeInfo, error) {
	conf, err := ini.Load(io.CurDir() + ".conf")
	if err != nil {
		return nil, err
	}

	nodeInfo := gen_server.NewNodeInfo()

	nodeInfo.Id = conf.Get("node", "id")
	nodeInfo.Url = conf.Get("node", "url")
	nodeInfo.AutoDetect = conf.GetInt("node", "autodetect", 0) > 0
	nodeInfo.Public = conf.GetInt("node", "public", 0) > 0

	nodeInfo.RedUrl = conf.Get("redis", "url")
	nodeInfo.RedPool = conf.GetInt("redis", "pool", 0)
	nodeInfo.RedDB = conf.GetInt("redis", "db", 0)
	nodeInfo.RedAuth = conf.Get("redis", "auth")

	return nodeInfo, nil
}
