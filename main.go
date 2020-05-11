package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

func main() {

	// Set logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	path, err := os.Getwd()

	if err != nil {
		path = "."
	}

	// Create astilectron
	a, err := astilectron.New(l, astilectron.Options{
		AppName:            "Dream Tube",
		AppIconDefaultPath: path + "/content/icons/icon.png",
	})
	if err != nil {
		l.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		l.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	// New window
	var w *astilectron.Window
	if w, err = a.NewWindow("./content/html/index.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(700),
		Width:  astikit.IntPtr(700),
	}); err != nil {
		l.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

	// Create windows
	if err = w.Create(); err != nil {
		l.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}

	videoIndex := 1
	audioIndex := 1
	// This will listen to messages sent by Javascript
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var s string
		m.Unmarshal(&s)

		// Process next video
		if s == "getNextVideo" {
			files := getFiles("./content/mp4")

			if len(files) == 0 {
				return nil
			} else if len(files) > videoIndex {

				file := files[videoIndex]
				videoIndex++
				return file
			} else {
				videoIndex = 1
				return files[videoIndex]
			}
		}

		// Process next video
		if s == "getNextAudio" {
			files := getFiles("./content/mp3")

			if len(files) == 0 {
				return nil
			} else if len(files) > audioIndex {

				file := files[audioIndex]
				audioIndex++
				return file
			} else {
				audioIndex = 1
				return files[audioIndex]
			}
		}
		return nil
	})
	w.OpenDevTools()
	// Blocking pattern
	a.Wait()

}

func getFiles(root string) []string {
	var files []string

	filesLen, _ := ioutil.ReadDir(root)
	if len(filesLen) == 0 {
		return files
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, strings.Replace(path, "content\\", "..\\", -1))
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}
