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
	this := new(LogRecord)
	this.Logs = make([]*RuntimeLog, 0, 0)
	return this
}

func NewRuntimeLog() *RuntimeLog {
	this := new(RuntimeLog)
	return this
}

func (this *Grid) addRuntimeLog(action string) error {
	logFile := this.proj + "runtime_log.json"
	rec := NewLogRecord()
	exist := this.fileExists(logFile)
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
	newLog.Runtime = this.runtime
	newLog.VerName = this.vername
	newLog.Version = this.version
	newLog.VerInfo = this.verinfo

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

func (this *Grid) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (this *Grid) mvOldRuntime() error {
	err := os.MkdirAll(this.proj+"backoff/", os.ModePerm)
	if err != nil {
		return err
	}
	err2 := os.Rename(this.proj+this.oldtime, this.proj+"backoff/"+this.oldtime)
	if err2 != nil {
		return err2
	}
	return nil
}
