package main

import (
	"github.com/pagedegeek/golog"
	"io"
	"os"
)

func main() {
	f, err := os.Create("/tmp/log")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := io.MultiWriter(os.Stdout, f)

	l := golog.New(golog.INFO, w, nil)
	l.Debug("Debug msg")
	l.Info("Info msg")
	l.Error("Error msg")
}
