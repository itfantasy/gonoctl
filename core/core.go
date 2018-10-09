package core

import (
	"errors"
	"plugin"

	"github.com/fsnotify/fsnotify"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/args"
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
	soPath  string
	soFile  string
	soVer   int
	nodeId  string
	nodeUrl string
}

var _core *Grid = nil

func core() *Grid {
	if _core == nil {
		_core = new(Grid)
	}
	return _core
}

func (this *Grid) initialize(parser *args.ArgParser) error {
	soPath, exist := parser.Get("-t")
	if !exist {
		return errors.New("(-t) please set the target dir of the runtime!")
	}
	this.soPath = soPath

	soFile, exist := parser.Get("-l")
	this.soFile = soFile

	nodeId, exist := parser.Get("-i")
	this.nodeId = nodeId

	nodeUrl, exist := parser.Get("-u")
	this.nodeUrl = nodeUrl

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if watcher.Add(this.soPath); err != nil {
		return err
	}
	this.watcher = watcher

	return nil
}

func (this *Grid) configParser() *args.ArgParser {
	parser := args.Parser().
		AddArg("-d", "", "set the target dir of the runtime").
		// AddArg("-f", "", "select the conf file"). // default .grconf
		AddArg("-l", "", "set the latest runtime.so"). // default save the name to the file .gnruntime, and neednot next time
		AddArg("-i", "", "dynamic set the id of the node").
		AddArg("-u", "", "dynamic set the urls")
	return parser
}

func (this *Grid) createNodeInfo() *gen_server.NodeInfo {
	nodeInfo := gen_server.NewNodeInfo()

	return nodeInfo
}

func (this *Grid) autoHotUpdate() error {
	so, err := plugin.Open(this.soPath + this.soFile)
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
	so, err := plugin.Open(this.soPath + this.soFile)
	if err != nil {
		return err
	}
	funcLaunch, err := so.Lookup("Launch")
	if err != nil {
		return err
	}
	funcLaunch.(func(*gen_server.NodeInfo))(this.createNodeInfo())
	return nil
}
