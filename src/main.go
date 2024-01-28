package main

import (
	"proxypool/pkg/getter"
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

func main() {
	//getter.IP89()
	getter.IP3306()

	defer log.Stop()
}
