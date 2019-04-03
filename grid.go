package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"plugin"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/itfantasy/grid/utils/args"
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
	oldtime string

	version int
	vername string
	verinfo string

	nodeId  string
	nodeUrl string
	pubUrl  string
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
	if !strings.HasSuffix(proj, "/") {
		this.proj = proj + "/"
	}

	runtime, err := this.selectTheRuntime()
	if err != nil {
		return err
	}
	this.runtime = runtime

	if !this.tryK8sEvns() {
		if nodeId, b := parser.Get("i"); b {
			this.nodeId = nodeId
		}
		if nodeUrl, b := parser.Get("l"); b {
			this.nodeUrl = nodeUrl
		}
		if pubUrl, b := parser.Get("p"); b {
			this.pubUrl = pubUrl
		}
	}

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
		AddArg("d", "", "set the target dir of the runtime").
		AddArg("i", "", "dynamic set the id of the node").
		AddArg("l", "", "dynamic set the local url").
		AddArg("p", "", "dynamic set the public url")
	return parser
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
	err3 := this.addRuntimeLog("Update")
	if err3 != nil {
		fmt.Println("[RuntimeLog]::" + err3.Error())
	}
	err4 := this.mvOldRuntime()
	if err4 != nil {
		fmt.Println("[MvOldRuntime]::" + err4.Error())
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

	funcLaunch, err := so.Lookup("OnLaunch")
	if err != nil {
		return err
	}

	err3 := this.addRuntimeLog("Launch")
	if err3 != nil {
		fmt.Println("[RuntimeLog]::" + err3.Error())
	}

	funcLaunch.(func(string, string, string, string))(this.proj, this.nodeId, this.nodeUrl, this.pubUrl)

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

func (this *Grid) printVersionInfo() {
	fmt.Println("--------" + this.runtime + "--------")
	fmt.Println(" ver:	" + this.vername + "|" + strconv.Itoa(this.version))
	fmt.Println(" info:	" + this.verinfo)
	fmt.Println("----------------------------")
}

func (this *Grid) selectTheRuntime() (string, error) {
	dir, err := ioutil.ReadDir(this.proj)
	if err != nil {
		return "", err
	}
	suffix := "so"
	newVer := -1
	newFile := ""
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		fileName := fi.Name()
		if strings.HasSuffix(fileName, suffix) {
			fmt.Println("found and checking ... " + fileName)
			ver, err := this.getRunTimeInfo(fileName)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				if ver > newVer {
					newFile = fileName
					newVer = ver
					fmt.Println("a newer version: " + strconv.Itoa(newVer) + " ... " + fileName)
				} else {
					fmt.Println("old version: " + strconv.Itoa(ver) + " ... " + fileName)
				}
			}
		}
	}
	if newFile != "" {
		return newFile, nil
	}
	return "", errors.New("the appropriate runtime was not found!!")
}

func (this *Grid) tryK8sEvns() bool {
	GRID_NODE_ID := os.Getenv("GRID_NODE_ID")
	if GRID_NODE_ID == "" {
		return false
	}

	GRID_NODE_NAME := os.Getenv("GRID_NODE_NAME")
	GRID_NODE_PORT := os.Getenv("GRID_NODE_PORT")
	GRID_NODE_PROTO := os.Getenv("GRID_NODE_PROTO")
	GRID_NODE_ISPUB := os.Getenv("GRID_NODE_ISPUB")
	GRID_LOCAL_IP := os.Getenv("GRID_LOCAL_IP")
	GRID_PUB_IP := os.Getenv("GRID_PUB_IP")

	this.nodeId = GRID_NODE_ID
	this.nodeUrl = GRID_NODE_PROTO + "://" + GRID_LOCAL_IP + ":" + GRID_NODE_PORT
	if GRID_NODE_PROTO == "ws" {
		this.nodeUrl += "/" + GRID_NODE_NAME
	}
	if GRID_NODE_ISPUB == "TRUE" {
		this.pubUrl = GRID_NODE_PROTO + "://" + GRID_PUB_IP + ":" + GRID_NODE_PORT
		if GRID_NODE_PROTO == "ws" {
			this.pubUrl += "/" + GRID_NODE_NAME
		}
	} else {
		this.pubUrl = ""
	}
	return true
}
