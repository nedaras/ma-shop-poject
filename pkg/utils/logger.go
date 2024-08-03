package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/gommon/log"
)

var (
	logger *log.Logger
)

func Logger() *log.Logger {
	fmt.Println("geting the logger")
	if logger == nil {
		exe, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		logs := filepath.Join(filepath.Dir(exe), "logs")
		file := filepath.Join(logs, time.Now().Format(time.DateTime)+".txt")

		_, err = os.Stat(logs)
		if os.IsNotExist(err) {
			err = os.Mkdir(logs, os.ModePerm)
		}

		if err != nil {
			log.Fatal(err)
		}

		out, err := os.Create(file)
		if err != nil {
			log.Fatal(err)
		}

		logger = log.New("logs")
		logger.SetOutput(out)
	}

	return logger
}
