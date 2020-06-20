package main

import (
	"fmt"
	"strings"

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
					if g.isSoLibFile(ev.Name) {
						fmt.Println("[watcher]::find a new runtime : ", ev.Name)
						g.runtime = ev.Name
						if err := g.autoHotUpdate(); err != nil {
							fmt.Println(err.Error())
							continue
						}
						fmt.Println("[watcher]::an new version : " + ev.Name + " has been loaded !")
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
