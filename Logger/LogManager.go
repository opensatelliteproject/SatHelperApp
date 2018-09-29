package SLog

import (
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
)

var displayOnTermUI = false

var logStarted = false
var logPath string
var logPathFlag = flag.String("logdir", "", "Log folder")
var logFileHandle *os.File

func StartLog() {
	if !logStarted {
		flag.Parse()
		home, _ := homedir.Dir()
		logPath = fmt.Sprintf("%s/SatHelperApp/logs", home)
		if *logPathFlag != "" {
			logPath = *logPathFlag
		}
		logStarted = true
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
		file, err := os.Create(fmt.Sprintf("%s/log.txt", logPath))
		if err != nil {
			Error("Error opening log file: %s. Won't try again...", err)
			return
		}
		Info(aurora.Bold("Log Started at %s/log.txt . If no message apears please check the log.").String(), logPath)
		logFileHandle = file
	}
}

func EndLog() {
	if logFileHandle != nil {
		logFileHandle.Close()
	}
}

func SetTermUiDisplay(b bool) {
	displayOnTermUI = b
}

func Info(str string, v ...interface{}) {
	Log(str, v...)
}

func Log(str string, v ...interface{}) {
	if displayOnTermUI {
		log.Printf("[I](fg-bold) [%s](fg-cyan)\n", fmt.Sprintf(str, v...))
	} else {
		log.Printf(aurora.Cyan("%s").String(), fmt.Sprintf(str, v...))
	}

	if logFileHandle != nil {
		_, err := logFileHandle.WriteString(aurora.Cyan(fmt.Sprintf("[I] %s\n", fmt.Sprintf(str, v...))).String())
		if err != nil {
			panic(err)
		}
	}
}

func Debug(str string, v ...interface{}) {
	if displayOnTermUI {
		log.Printf("[D](fg-bold) [%s](fg-magenta)\n", fmt.Sprintf(str, v...))
	} else {
		log.Printf(aurora.Magenta("%s").String(), fmt.Sprintf(str, v...))
	}
	if logFileHandle != nil {
		_, err := logFileHandle.WriteString(aurora.Magenta(fmt.Sprintf("[D] %s\n", fmt.Sprintf(str, v...))).String())
		if err != nil {
			panic(err)
		}
	}
}

func Warn(str string, v ...interface{}) {
	if displayOnTermUI {
		log.Printf("[W](fg-bold) [%s](fg-yellow)\n", fmt.Sprintf(str, v...))
	} else {
		log.Printf(aurora.Brown("%s").String(), fmt.Sprintf(str, v...))
	}
	if logFileHandle != nil {
		_, err := logFileHandle.WriteString(aurora.Brown(fmt.Sprintf("[W] %s\n", fmt.Sprintf(str, v...))).String())
		if err != nil {
			panic(err)
		}
	}
}

func Error(str string, v ...interface{}) {
	if displayOnTermUI {
		log.Printf("[E](fg-bold) [%s](fg-red)\n", fmt.Sprintf(str, v...))
	} else {
		log.Printf(aurora.Red("%s").String(), fmt.Sprintf(str, v...))
	}
	if logFileHandle != nil {
		_, err := logFileHandle.WriteString(aurora.Red(fmt.Sprintf("[E] %s\n", fmt.Sprintf(str, v...))).String())
		if err != nil {
			panic(err)
		}
	}
}
