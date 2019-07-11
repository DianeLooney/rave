package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dianelooney/directive"
	"github.com/dianelooney/rave"
	"github.com/fsnotify/fsnotify"
)

const src = "client/src.rave"

var ctx rave.Context

func init() {
	ctx.Init()
}

func loadFile() {
	bytes, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatalf("Unable to read file: %v\n", err)
	}
	e, err := directive.Prepare(bytes)
	if err != nil {
		log.Fatalf("Unable to prepare directives: %v\n", err)
	}
	d := rave.Doc{}
	err = e.Execute(&d)
	if err != nil {
		log.Printf("Unable to execute directive: %v\n", err)
	}
	ctx.Load(&d)
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	watcher.Add(src)

	go func() {
		for {
			if (<-watcher.Events).Op != fsnotify.Write {
				continue
			}
			loadFile()
		}
	}()

	loadFile()

	<-make(chan bool)
}
