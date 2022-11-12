package utils

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var outPath string

var currPath string

var logFile *os.File

var logFilePath string

var logger *log.Logger

// procedures that gets the current paths and create, eventually, the new folders and log files based on timestamp
func Initialize(service int) {
	var err error
	currPath, err = os.Getwd()
	if err != nil {
		log.Fatal("Unable to get the current path!")
	}
	getPath(service)
}

// procedure that generates a path, different one based on the current service
func getPath(service int) {

	var err error
	switch service {
	case 0:
		outPath = filepath.Join(currPath, "logs", "registrator")

		break
	case 1:
		outPath = filepath.Join(currPath, "logs", "coordinator")
		break
	case 2:
		outPath = filepath.Join(currPath, "logs", "tkcen")
		break
	case 3:
		outPath = filepath.Join(currPath, "logs", "tkdec")
		break
	case 4:
		outPath = filepath.Join(currPath, "logs", "ricandagr")
		break
	default:
		log.Fatal("Unable to get the service requested! (must be between 0 and 4 included) ")
	}

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		err = os.Mkdir(outPath, 0755)
		if err != nil {
			log.Fatal("Unable to create the directory!")
		}
	}
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	var logFileName = "log" + timestamp + ".log"
	logFilePath = filepath.Join(outPath, logFileName)

	logFile, err = os.Create(logFilePath)
	if err != nil {
		log.Fatal("Unable to create the file!")
	}

}

// procedure to call when someone wants to write on his log file
func WriteInLog(line string, verbose bool) {
	if verbose {
		f, err := os.OpenFile(logFilePath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		if err != nil {
			log.Fatal("Unable to open the log file!")
		}
		logger = log.New(f, "", log.LstdFlags)
		logger.Println(line)
	}
}
