package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/itfantasy/gonode/utils/args"
	gio "github.com/itfantasy/gonode/utils/io"
	"github.com/itfantasy/gonode/utils/ts"
)

func run() error {
	parser := args.Parser().
		AddArg("proj", "", "set the project name of the runtime")
	projName, _ := parser.Get("proj")
	if projName != "" {
		return buildTheRunTime(projName)
	}
	directory, err := ioutil.ReadDir(gio.CurrentDir())
	if err != nil {
		return err
	}
	for _, fi := range directory {
		if fi.IsDir() {
			projName := fi.Name()
			return buildTheRunTime(projName)
		}
	}
	return nil
}

func buildTheRunTime(projName string) error {

	time := ts.Time()
	fileName := projName
	fileName += "_" + ts.TimeToStr(time, ts.FORMAT_NOW_C)
	fileName += ".so"

	projPath := projName + "/"

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", projPath+fileName, projPath+"runtime.go", projPath+"node.go")
	//fmt.Println(cmd.Args)
	fmt.Print(projName + " buiding...  ")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			fmt.Println(err2)
			break
		}
		fmt.Println(line)
	}
	cmd.Wait()

	// create the .meta file
	createMetaFile(projPath, fileName, projName, ts.TimeToStr(time, ts.FORMAT_NOW_A))

	fmt.Println(projPath + fileName + " build succeed!!")
	return nil
}

func createMetaFile(projPath string, fileName string, projName string, date string) {
	meta := ""
	meta = "project: " + projName + "\r\n"
	meta += "date: " + date + "\r\n"
	file, err := os.Open(projPath + "runtime.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		bytes, _, c := reader.ReadLine()
		if c == io.EOF {
			break
		}
		txt := string(bytes)
		if strings.HasPrefix(txt, "//") && strings.Contains(txt, "@Meta.VersionName") {
			reader.ReadLine()
			b, _, _ := reader.ReadLine()
			str := string(b)
			str = strings.Replace(str, "return", "", -1)
			str = strings.Replace(str, "\"", "", -1)
			str = strings.Replace(str, " ", "", -1)
			meta += "versionName: " + str
		}

		if strings.HasPrefix(txt, "//") && strings.Contains(txt, "@Meta.VersionCode") {
			reader.ReadLine()
			b, _, _ := reader.ReadLine()
			str := string(b)
			str = strings.Replace(str, "return", "", -1)
			str = strings.Replace(str, " ", "", -1)
			meta += "\r\n"
			meta += "versionCode: " + str
		}

		if strings.HasPrefix(txt, "//") && strings.Contains(txt, "@Meta.VersionInfo") {
			reader.ReadLine()
			b, _, _ := reader.ReadLine()
			str := string(b)
			str = strings.Replace(str, "return", "", -1)
			str = strings.Replace(str, "\"", "", -1)
			str = strings.Replace(str, " ", "", -1)
			meta += "\r\n"
			meta += "versionInfo: " + str
		}
	}

	if err := gio.SaveFile(projPath+strings.Replace(fileName, ".so", ".meta", -1), meta); err != nil {
		fmt.Println(err)
	}
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
	}
}
