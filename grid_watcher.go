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
func (this *Grid) watchingDirectory() {
	defer this.watcher.Close()
	for {
		select {
		case ev := <-this.watcher.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if this.isSoLibFile(ev.Name) {
						fmt.Println("find a new file : ", ev.Name)
						ver, err := this.getRunTimeInfo(ev.Name)
						if err != nil {
							fmt.Println(err.Error())
							continue
						}
						if ver > this.version {
							this.version = ver
							this.runtime = ev.Name
							if this.autoHotUpdate(); err != nil {
								fmt.Println(err.Error())
								continue
							}
							fmt.Println("an new version : " + ev.Name + " has been loaded !")
						} else {
							fmt.Println("no need update !")
						}
					}
				}
			}
		case err := <-this.watcher.Errors:
			{
				fmt.Println("error : ", err)
				return
			}
		}
	}
}

func (this *Grid) isSoLibFile(fileName string) bool {
	infos := strings.Split(fileName, ".")
	extendName := infos[len(infos)-1]
	if extendName == "so" {
		return true
	}
	return false
}

func (this *Grid) getRunTimeInfo(fileName string) (int, error) {

	return 0, nil
}
