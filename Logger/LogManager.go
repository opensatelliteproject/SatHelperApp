package SLog

import (
	"log"
	"fmt"
	"github.com/logrusorgru/aurora"
)

var displayOnTermUI = false

func SetTermUiDisplay(b bool) {
	displayOnTermUI = b
}

func Info(str string, v ... interface{}) {
	Log(str, v ...)
}

func Log(str string, v ...interface{}) {
	if displayOnTermUI {
		log.Printf("[I](fg-bold) [%s](fg-cyan)\n", fmt.Sprintf(str, v...))
	} else {
		log.Printf(aurora.Cyan("%s").String(), fmt.Sprintf(str, v...))
	}
}

func Debug(str string, v ...interface{}) {
	if displayOnTermUI {
		log.Printf("[D](fg-bold) [%s](fg-magenta)\n", fmt.Sprintf(str, v...))
	} else {
		log.Printf(aurora.Magenta("%s").String(), fmt.Sprintf(str, v...))
	}
}

func Warn(str string, v ...interface{}) {
	if displayOnTermUI {
		log.Printf("[W](fg-bold) [%s](fg-yellow)\n", fmt.Sprintf(str, v...))
	} else {
		log.Printf(aurora.Brown("%s").String(), fmt.Sprintf(str, v...))
	}
}

func Error(str string, v ...interface{}) {
	if displayOnTermUI {
		log.Printf("[E](fg-bold) [%s](fg-red)\n", fmt.Sprintf(str, v...))
	} else {
		log.Printf(aurora.Red("%s").String(), fmt.Sprintf(str, v...))
	}
}