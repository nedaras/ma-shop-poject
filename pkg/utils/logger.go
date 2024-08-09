package utils

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
)

var (
	logger *log.Logger
	mu     sync.Mutex
)

func Logger() *log.Logger {
	mu.Lock()
	defer mu.Unlock()

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
