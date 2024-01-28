package main

import (
	"encoding/json"
	"proxypool/pkg/getter"
	"proxypool/pkg/models"
	"sync"
	log "unknwon.dev/clog/v2"
)

func init() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}

	err = log.NewFile(
		log.FileConfig{
			Filename: "./logs/clog.log",
			Level:    log.LevelTrace,
			FileRotationConfig: log.FileRotationConfig{
				Rotate: true,
				Daily:  true,
				//MaxLines: 50,
			},
		},
	)

	if err != nil {
		panic("unable to create new logger with file: " + err.Error())
	}
}

type MyFunction func() []*models.IP

func run() {
	var ipList []*models.IP

	var wg sync.WaitGroup

	funs := []MyFunction{
		getter.IP3306,
		getter.IP89,
	}

	for _, fun := range funs {
		wg.Add(1)

		go func(f MyFunction) {
			temp := f()

			for _, ip := range temp {
				ipList = append(ipList, ip)
			}
			wg.Done()
		}(fun)
	}

	wg.Wait()

	jsonData, _ := json.Marshal(ipList)

	log.Info(string(jsonData))
}

func main() {
	run()
	defer log.Stop()
}
