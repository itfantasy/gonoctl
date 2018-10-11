package main

import (
	"errors"
	"fmt"
	"plugin"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/args"
	"github.com/itfantasy/gonode/utils/ini"
)

func main() {
	grid := NewGrid()
	parser := grid.configParser()
	if err := grid.initialize(parser); err != nil {
		fmt.Println(err)
	}
	go grid.watchingDirectory()
	if err := grid.autoRun(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("finished!!")
}

type Grid struct {
	watcher *fsnotify.Watcher
	proj    string
	runtime string
	version int
	nodeId  string
	nodeUrl string
	vername string
	verinfo string
}

func NewGrid() *Grid {
	this := new(Grid)
	return this
}

func (this *Grid) initialize(parser *args.ArgParser) error {
	proj, exist := parser.Get("d")
	if !exist {
		return errors.New("(d) please set the target dir of the runtime!")
	}
	this.proj = proj + "/"

	runtime, exist := parser.Get("l")
	this.runtime = runtime

	nodeId, exist := parser.Get("i")
	this.nodeId = nodeId

	nodeUrl, exist := parser.Get("u")
	this.nodeUrl = nodeUrl

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if watcher.Add(this.proj); err != nil {
		return err
	}
	this.watcher = watcher

	return nil
}

func (this *Grid) configParser() *args.ArgParser {
	parser := args.Parser().
		AddArg("d", "runtime", "set the target dir of the runtime").
		AddArg("l", "runtime_0.so", "set the latest runtime.so").
		AddArg("i", "", "dynamic set the id of the node").
		AddArg("u", "", "dynamic set the urls")
	return parser
}

func (this *Grid) setupConfig() (*gen_server.NodeInfo, error) {
	conf, err := ini.Load(this.proj + ".conf")
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

func (this *Grid) autoHotUpdate() error {
	so, err := plugin.Open(this.proj + this.runtime)
	if err != nil {
		return err
	}
	err2 := this.autoVersion(so)
	if err2 != nil {
		return err2
	}
	funcUpdate, err := so.Lookup("OnHotUpdate")
	if err != nil {
		return err
	}
	funcUpdate.(func())()
	return nil
}

func (this *Grid) autoRun() error {
	so, err := plugin.Open(this.proj + this.runtime)
	if err != nil {
		return err
	}
	err2 := this.autoVersion(so)
	if err2 != nil {
		return err2
	}
	this.printVersionInfo()
	funcLaunch, err := so.Lookup("Launch")
	if err != nil {
		return err
	}
	config, err := this.setupConfig()
	if err != nil {
		return err
	}
	funcLaunch.(func(*gen_server.NodeInfo))(config)
	return nil
}

func (this *Grid) autoVersion(so *plugin.Plugin) error {
	funcVersionName, err := so.Lookup("VersionName")
	if err != nil {
		return err
	}
	this.vername = funcVersionName.(func() string)()
	funcVersionInfo, err := so.Lookup("VersionInfo")
	if err != nil {
		return err
	}
	this.verinfo = funcVersionInfo.(func() string)()
	return nil
}

func (this *Grid) saveRuntimeName() error {
	// save the runtimename to .runtime as a default name for next time
	return nil
}

func (this Grid) printVersionInfo() {
	fmt.Println("--------" + this.runtime + "--------")
	fmt.Println(" ver:	" + this.vername + "|" + strconv.Itoa(this.version))
	fmt.Println(" info:	" + this.verinfo)
	fmt.Println("----------------------------")
}
