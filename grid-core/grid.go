package main

import (
	"encoding/json"
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
	fmt.Println(`
  ________      .__    .___
 /  _____/______|__| __| _/
/   \  __\_  __ \  |/ __ | 
\    \_\  \  | \/  / /_/ | 
 \______  /__|  |__\____ | 
        \/              \/  
:: Grid is An Addtional Light Engine for gonode to support Docker, K8s, Hotupdating and so on. ::

`)
	grid := NewGrid()
	parser := grid.configParser()
	if err := grid.initialize(parser); err != nil {
		fmt.Println(err)
	}
	go grid.watchingDirectory()
	if err := grid.autoRun(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("==========================================================")
}

type Grid struct {
	watcher *fsnotify.Watcher

	proj    string
	runtime string
	oldtime string

	version int
	vername string
	verinfo string

	namespace string
	nodeid    string
	endpoints []string
	etc       string
}

func NewGrid() *Grid {
	g := new(Grid)
	g.endpoints = make([]string, 0, 1)
	return g
}

func (g *Grid) initialize(parser *args.ArgParser) error {
	proj, exist := parser.Get("proj")
	if !exist {
		return errors.New("project dir (-proj) is necessary!")
	}
	if !strings.HasSuffix(proj, "/") {
		g.proj = proj + "/"
	}

	runtime, err := g.selectTheRuntime()
	if err != nil {
		return err
	}
	g.runtime = runtime

	k8sEvn, err := g.tryK8sEvns()
	if err != nil {
		return err
	}

	if !k8sEvn {
		if namespace, b := parser.Get("namespace"); b {
			g.namespace = namespace
		}
		if nodeid, b := parser.Get("nodeid"); b {
			g.nodeid = nodeid
		}
		if endpoints, b := parser.Get("endpoints"); b {
			g.endpoints = strings.Split(endpoints, ",")
		}
		if etc, b := parser.Get("etc"); b {
			g.etc = etc
		}
	}

	if g.nodeid == "" {
		return errors.New("nodeid (-nodeid) is necessary!")
	}

	if len(g.endpoints) <= 0 {
		return errors.New("endpoint list (-endpoints) is necessary!")
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
		AddArg("proj", "", "the project dir which will mount for grid").
		AddArg("namespace", "", "set the namespace, such as 'itfantasy'").
		AddArg("nodeid", "", "set the nodeid, such as 'game_1024'").
		AddArg("endpoints", "", "set the endpoint list, such as 'tcp://yourserver:30005,kcp://yourserver:30006,ws://yourserver:30007/game_1024'").
		AddArg("etc", "", "extra configs")
	return parser
}

func (g *Grid) autoHotUpdate() error {
	so, err := plugin.Open(g.proj + g.runtime)
	if err != nil {
		return err
	}
	if err := g.autoVersion(so); err != nil {
		return err
	}
	funcUpdate, err := so.Lookup("OnHotUpdate")
	if err != nil {
		return err
	}
	funcUpdate.(func())()
	return nil
}

func (g *Grid) autoRun() error {
	so, err := plugin.Open(g.proj + g.runtime)
	if err != nil {
		return err
	}
	if err := g.autoVersion(so); err != nil {
		return err
	}
	g.printVersionInfo()
	funcLaunch, err := so.Lookup("OnLaunch")
	if err != nil {
		return err
	}
	funcLaunch.(func(string, string, string, []string, string))(g.proj, g.namespace, g.nodeid, g.endpoints, g.etc)
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
	infos := strings.Split(g.runtime, "_")
	fmt.Println("")
	fmt.Println("--------" + g.runtime + "--------")
	fmt.Println(" proj:	" + infos[0])
	fmt.Println(" date:	" + infos[1])
	fmt.Println(" ver:	" + g.vername + "|" + strconv.Itoa(g.version))
	fmt.Println(" info:	" + g.verinfo)
	fmt.Println("------------------------------------------------")
	fmt.Println("")
}

func (g *Grid) selectTheRuntime() (string, error) {
	lstTime := 0
	lstFile := ""
	dir, err := ioutil.ReadDir(g.proj)
	if err != nil {
		return "", err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		fileName := fi.Name()
		if g.isSoLibFile(fileName) {
			fmt.Println("found and checking ... " + fileName)
			infos := strings.Split(strings.TrimSuffix(fileName, ".so"), "_")
			if len(infos) == 2 {
				if time, err := strconv.Atoi(infos[1]); err == nil {
					if time > lstTime {
						lstTime = time
						lstFile = fileName
					}
				}
			}
		}
	}
	if lstFile != "" {
		return lstFile, nil
	}
	return "", errors.New("the appropriate runtime was not found!!")
}

func (g *Grid) tryK8sEvns() (bool, error) {
	GRID_NODE_ID := os.Getenv("GRID_NODE_ID")
	if GRID_NODE_ID == "" {
		return false, nil
	}
	GRID_NODE_NAMESPACE := os.Getenv("GRID_NODE_NAMESPACE")
	GRID_NODE_ENDPOINTS := os.Getenv("GRID_NODE_ENDPOINTS")
	GRID_LOCAL_IP := os.Getenv("GRID_LOCAL_IP")

	g.nodeid = GRID_NODE_ID
	g.namespace = GRID_NODE_NAMESPACE

	endpoints := make([]string, 0, 1)
	if err := json.Unmarshal([]byte(GRID_NODE_ENDPOINTS), endpoints); err != nil {
		return false, err
	}
	for _, endpoint := range endpoints {
		infos := strings.Split(endpoint, "://:")
		if len(infos) < 2 {
			return false, errors.New("illegal endpoints!!")
		}
		g.endpoints = append(g.endpoints, infos[0]+"://"+GRID_LOCAL_IP+":"+infos[1])
	}

	return true, nil
}
