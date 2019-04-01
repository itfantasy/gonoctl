package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/itfantasy/grid/utils/crypt"
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
						fmt.Println("[watcher]::find a new runtime : ", ev.Name)
						ver, err := this.getRunTimeInfo(ev.Name)
						if err != nil {
							fmt.Println(err.Error())
							continue
						}
						if ver > this.version {
							this.version = ver
							this.oldtime = this.runtime
							this.runtime = ev.Name
							if this.autoHotUpdate(); err != nil {
								fmt.Println(err.Error())
								continue
							}
							fmt.Println("[watcher]::an new version : " + ev.Name + " has been loaded !")
							this.printVersionInfo()
						} else {
							fmt.Println("[watcher]::no need update !")
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
	infos := strings.Split(fileName, ".")
	if len(infos) != 2 {
		return -1, errors.New("illegal runtime file name! -1")
	}
	runtimeName := infos[0]
	strs := strings.Split(runtimeName, "_")
	if len(strs) != 4 {
		return -1, errors.New("illegal runtime file name! -2")
	}
	ver, err := strconv.Atoi(strs[1])
	if err != nil {
		return -1, errors.New("illegal runtime file name! -3")
	}
	md5 := strs[3]
	fileString := strs[0] + "_" + strs[1] + "_" + strs[2]
	sign := crypt.Md5("ITFANTASY-GRID-" + fileString)
	if md5 != sign {
		return -1, errors.New("illegal sign code!")
	}

	return ver, nil
}
