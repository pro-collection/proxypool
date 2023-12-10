package main

import (
	"errors"
	log "unknwon.dev/clog/v2"
)

func init() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}

	err = log.NewFile(log.FileConfig{
		Filename: "./logs/clog.log",
		Level:    log.LevelTrace,
		FileRotationConfig: log.FileRotationConfig{
			Rotate: true,
			Daily:  true,
			//MaxLines: 50,
		},
	})

	if err != nil {
		panic("unable to create new logger with file: " + err.Error())
	}
}

func main() {
	err := errors.New("错误信息")

	log.Trace("Hello %s!", "World")
	log.Info("Hello %s!", "World")
	log.Warn("Hello %s!", "World")
	log.Error("So bad... %v", err)
	log.Fatal("Boom! %v", err)

	defer log.Stop()
}
