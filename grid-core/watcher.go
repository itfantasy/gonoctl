package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// watch the target Directory
// check if there is a new file been created
// check if the file version is greater than runtime version
// then replace the runtime function interface of now
func (g *Grid) watchingDirectory() {
	defer g.watcher.Close()
	for {
		select {
		case ev := <-g.watcher.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					wholePath := ev.Name
					infos := strings.Split(wholePath, "/")
					soName := infos[len(infos)-1]
					if g.isSoLibFile(soName) {
						fmt.Println("[grid-core]::find a new runtime : ", soName)
						g.runtime = soName
						<-time.After(time.Millisecond * time.Duration(1000))
						fmt.Print("[grid-core]::loading .")
						<-time.After(time.Millisecond * time.Duration(1000))
						fmt.Print(".")
						<-time.After(time.Millisecond * time.Duration(1000))
						fmt.Println(".")
						if err := g.autoHotUpdate(); err != nil {
							fmt.Println(err.Error())
							continue
						}
						fmt.Println("[grid-core]::an new version : " + soName + " has been loaded !")
						g.printVersionInfo()
					}
				}
			}
		case err := <-g.watcher.Errors:
			{
				fmt.Println("error : ", err)
				return
			}
		}
	}
}

func (g *Grid) isSoLibFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".so")
}
