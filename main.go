package main

import (
	"log"
	"time"
	"yiban/app"

	"github.com/fatih/color"
)

func main() {
	log.SetFlags(log.Lmsgprefix)
	log.Default().SetPrefix(color.HiCyanString("[%s]", time.Now().Format("2006-01-02 15:04:05")))
	app.App().Execute()
}
