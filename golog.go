package golog

import (
	"fmt"
	"io"
	"time"
)

type (
	Golog struct {
		Level    int
		Writer   io.Writer
		Formater Formater
	}

	Formater interface {
		Format(level int, format string, a ...interface{}) string
	}

	StdFormater struct{}
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (f *StdFormater) Format(level int, format string, a ...interface{}) string {
	var levelName string
	switch level {
	case 1:
		levelName = "INFO"
	case 2:
		levelName = "WARN"
	case 3:
		levelName = "ERROR"
	case 4:
		levelName = "FATAL"
	default:
		levelName = "DEBUG"
	}

	msg := fmt.Sprintf(format, a...)
	return fmt.Sprintf("%s [%s] %s\r\n", time.Now().Format(time.RFC3339), levelName, msg)
}

func New(level int, w io.Writer, formater Formater) *Golog {
	if formater == nil {
		formater = &StdFormater{}
	}
	return &Golog{
		Level:    level,
		Writer:   w,
		Formater: formater,
	}
}

func (g *Golog) Printf(level int, format string, a ...interface{}) error {
	if level >= g.Level {
		msg := g.Formater.Format(level, format, a...)
		if _, err := g.Writer.Write([]byte(msg)); err != nil {
			return err
		}

	}
	return nil
}

func (g *Golog) Debug(format string, a ...interface{}) error {
	return g.Printf(DEBUG, format, a...)
}

func (g *Golog) Info(format string, a ...interface{}) error {
	return g.Printf(INFO, format, a...)
}

func (g *Golog) Warn(format string, a ...interface{}) error {
	return g.Printf(WARN, format, a...)
}

func (g *Golog) Error(format string, a ...interface{}) error {
	return g.Printf(ERROR, format, a...)
}

func (g *Golog) Fatal(format string, a ...interface{}) error {
	return g.Printf(FATAL, format, a...)
}
