package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

func Init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	})

	err := os.MkdirAll("logs", 0644)
	if err != nil {
		panic(err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)
}
