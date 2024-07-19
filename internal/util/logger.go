package util

import (
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/log"
)

var (
	Logger *log.Logger
	once   sync.Once
)

func init() {
	once.Do(func() {
		// Check if the directory is created
		nvmDir, err := GetNvmDirectory()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		logFile, err := os.OpenFile(fmt.Sprintf("%s/debug.log", nvmDir), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("Failed to open log file:", err)
		}

		Logger = log.New(logFile)
		Logger.SetReportTimestamp(true)
		Logger.SetReportCaller(true)
		Logger.SetLevel(log.DebugLevel)
	})
}
