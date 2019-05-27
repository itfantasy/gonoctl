package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/itfantasy/grid/utils/ts"
)

type RuntimeLog struct {
	Date    string `json:"date"`
	Action  string `json:"action"`
	Runtime string `json:"runtime"`
	VerName string `json:"vername"`
	Version int    `json:"version"`
	VerInfo string `json:"verinfo"`
	Status  string `json:"status"`
}

type LogRecord struct {
	Logs []*RuntimeLog
}

func NewLogRecord() *LogRecord {
	g := new(LogRecord)
	g.Logs = make([]*RuntimeLog, 0, 0)
	return g
}

func NewRuntimeLog() *RuntimeLog {
	g := new(RuntimeLog)
	return g
}

func (g *Grid) addRuntimeLog(action string) error {
	logFile := g.proj + "runtime_log.json"
	rec := NewLogRecord()
	exist := g.fileExists(logFile)
	if exist {
		bytes, err := ioutil.ReadFile(logFile)
		if err != nil {
			return err
		}
		err2 := json.Unmarshal(bytes, rec)
		if err2 != nil {
			return err2
		}
	}

	newLog := NewRuntimeLog()
	newLog.Date = ts.NowToStr(ts.Now(), ts.FORMAT_NOW_A)
	newLog.Action = action
	newLog.Runtime = g.runtime
	newLog.VerName = g.vername
	newLog.Version = g.version
	newLog.VerInfo = g.verinfo

	rec.Logs = append(rec.Logs, newLog)

	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	err2 := ioutil.WriteFile(logFile, data, 0644)
	if err2 != nil {
		return err2
	}

	return nil
}

func (g *Grid) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (g *Grid) mvOldRuntime() error {
	err := os.MkdirAll(g.proj+"backoff/", os.ModePerm)
	if err != nil {
		return err
	}
	err2 := os.Rename(g.proj+g.oldtime, g.proj+"backoff/"+g.oldtime)
	if err2 != nil {
		return err2
	}
	return nil
}
