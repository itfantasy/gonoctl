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
	fmt.Println(`
  ________      .__    .___
 /  _____/______|__| __| _/
/   \  __\_  __ \  |/ __ | 
\    \_\  \  | \/  / /_/ | 
 \______  /__|  |__\____ | 
        \/              \/  
:: An Addtional Light Engine for gonode to support Docker, Hotupdating and so on. ::

`)
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
	pub     bool
}

func NewGrid() *Grid {
	g := new(Grid)
	return g
}

func (g *Grid) initialize(parser *args.ArgParser) error {
	proj, exist := parser.Get("d")
	if !exist {
		return errors.New("(d) please set the target dir of the runtime!")
	}
	if !strings.HasSuffix(proj, "/") {
		g.proj = proj + "/"
	}

	runtime, err := g.selectTheRuntime()
	if err != nil {
		return err
	}
	g.runtime = runtime

	if !g.tryK8sEvns() {
		if nodeId, b := parser.Get("i"); b {
			g.nodeId = nodeId
		}
		if nodeUrl, b := parser.Get("l"); b {
			g.nodeUrl = nodeUrl
		}
		if pub, b := parser.Get("p"); b {
			g.pub = (pub == "y" || pub == "Y")
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if watcher.Add(g.proj); err != nil {
		return err
	}
	g.watcher = watcher

	return nil
}

func (g *Grid) configParser() *args.ArgParser {
	parser := args.Parser().
		AddArg("d", "", "set the target dir of the runtime").
		AddArg("i", "", "dynamic set the id of the node").
		AddArg("l", "", "dynamic set the local url").
		AddArg("p", "n", "dynamic set if the node is public(n or y)")
	return parser
}

func (g *Grid) autoHotUpdate() error {
	so, err := plugin.Open(g.proj + g.runtime)
	if err != nil {
		return err
	}
	err2 := g.autoVersion(so)
	if err2 != nil {
		return err2
	}
	funcUpdate, err := so.Lookup("OnHotUpdate")
	if err != nil {
		return err
	}
	err3 := g.addRuntimeLog("Update")
	if err3 != nil {
		fmt.Println("[RuntimeLog]::" + err3.Error())
	}
	err4 := g.mvOldRuntime()
	if err4 != nil {
		fmt.Println("[MvOldRuntime]::" + err4.Error())
	}
	funcUpdate.(func())()
	return nil
}

func (g *Grid) autoRun() error {
	so, err := plugin.Open(g.proj + g.runtime)
	if err != nil {
		return err
	}
	err2 := g.autoVersion(so)
	if err2 != nil {
		return err2
	}
	g.printVersionInfo()

	funcLaunch, err := so.Lookup("OnLaunch")
	if err != nil {
		return err
	}

	err3 := g.addRuntimeLog("Launch")
	if err3 != nil {
		fmt.Println("[RuntimeLog]::" + err3.Error())
	}

	funcLaunch.(func(string, string, string, bool))(g.proj, g.nodeId, g.nodeUrl, g.pub)

	return nil
}

func (g *Grid) autoVersion(so *plugin.Plugin) error {
	funcVersionName, err := so.Lookup("VersionName")
	if err != nil {
		return err
	}
	g.vername = funcVersionName.(func() string)()
	funcVersionInfo, err := so.Lookup("VersionInfo")
	if err != nil {
		return err
	}
	g.verinfo = funcVersionInfo.(func() string)()
	return nil
}

func (g *Grid) printVersionInfo() {
	fmt.Println("--------" + g.runtime + "--------")
	fmt.Println(" ver:	" + g.vername + "|" + strconv.Itoa(g.version))
	fmt.Println(" info:	" + g.verinfo)
	fmt.Println("----------------------------")
}

func (g *Grid) selectTheRuntime() (string, error) {
	dir, err := ioutil.ReadDir(g.proj)
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
			ver, err := g.getRunTimeInfo(fileName)
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

func (g *Grid) tryK8sEvns() bool {
	GRID_NODE_ID := os.Getenv("GRID_NODE_ID")
	if GRID_NODE_ID == "" {
		return false
	}

	GRID_NODE_NAME := os.Getenv("GRID_NODE_NAME")
	GRID_NODE_PORT := os.Getenv("GRID_NODE_PORT")
	GRID_NODE_PROTO := os.Getenv("GRID_NODE_PROTO")
	GRID_NODE_ISPUB := os.Getenv("GRID_NODE_ISPUB")
	GRID_LOCAL_IP := os.Getenv("GRID_LOCAL_IP")

	g.nodeId = GRID_NODE_ID
	g.nodeUrl = GRID_NODE_PROTO + "://" + GRID_LOCAL_IP + ":" + GRID_NODE_PORT
	if GRID_NODE_PROTO == "ws" {
		g.nodeUrl += "/" + GRID_NODE_NAME
	}
	if GRID_NODE_ISPUB == "TRUE" {
		g.pub = true
	} else {
		g.pub = false
	}
	return true
}
