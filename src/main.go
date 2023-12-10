package main

import (
	log "unknwon.dev/clog/v2"
)

func init() {
	err := log.NewConsole(
		log.DefaultFileName,
		0,
		log.ConsoleConfig{Level: log.LevelTrace},
	)

	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}

func main() {
	log.Trace("hello %s", "world")
	log.Info("Hello %s!", "World") // YYYY/MM/DD 12:34:56 [ INFO] Hello World!
	log.Warn("Hello %s!", "World") // YYYY/MM/DD 12:34:56 [ WARN] Hello World!

	defer log.Stop()
}
