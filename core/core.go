package core

import (
	"errors"
	"plugin"

	"github.com/fsnotify/fsnotify"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/args"
	"github.com/itfantasy/gonode/utils/ini"
)

func Run() error {
	parser := core().configParser()
	if err := core().initialize(parser); err != nil {
		return err
	}
	go core().watchingDirectory()
	if err := core().autoRun(); err != nil {
		return err
	}
	return nil
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

var _core *Grid = nil

func core() *Grid {
	if _core == nil {
		_core = new(Grid)
	}
	return _core
}

func (this *Grid) initialize(parser *args.ArgParser) error {
	proj, exist := parser.Get("-d")
	if !exist {
		return errors.New("(-d) please set the target dir of the runtime!")
	}
	this.proj = proj

	runtime, exist := parser.Get("-l")
	this.runtime = runtime

	nodeId, exist := parser.Get("-i")
	this.nodeId = nodeId

	nodeUrl, exist := parser.Get("-u")
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
		AddArg("-d", "", "set the target dir of the runtime").
		AddArg("-l", "", "set the latest runtime.so").
		AddArg("-i", "", "dynamic set the id of the node").
		AddArg("-u", "", "dynamic set the urls")
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

func (this *Grid) saveRuntimeName() error {
	// save the runtimename to .runtime as a default name for next time
	return nil
}
