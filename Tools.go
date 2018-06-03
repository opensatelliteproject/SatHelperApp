package main

import (
	"github.com/foize/go.fifo"
	"log"
	"runtime"
	"os/exec"
	"os"
)

func AddToFifoC64(fifo *fifo.Queue, arr []complex64, length int) {
	for i := 0; i < length; i++ {
		if fifo.Len() >= FifoSize {
			log.Printf("FIFO Overflowing!!")
			break
		}
		fifo.Add(arr[i])
	}
}

func AddToFifoS16(fifo *fifo.Queue, arr []int16, length int) {
	for i := 0; i < length; i++ {
		if fifo.Len() >= FifoSize {
			log.Printf("FIFO Overflowing!!")
			break
		}
		fifo.Add(arr[i])
	}
}

func AddToFifoS8(fifo *fifo.Queue, arr []int8, length int) {
	for i := 0; i < length; i++ {
		if fifo.Len() >= FifoSize {
			log.Printf("FIFO Overflowing!!")
			break
		}
		fifo.Add(arr[i])
	}
}

func swapBuffers(a **complex64, b **complex64) {
	c := *b
	*b = *a
	*a = c
}

func checkAndResizeBuffers(length int) {
	if len(buffer0) < length {
		buffer0 = make([]complex64, length)
	}
	if len(buffer1) < length {
		buffer1 = make([]complex64, length)
	}
}

func shiftWithConstantSize(arr *[]byte, pos int, length int) {
	for i := 0; i < length-pos; i++ {
		(*arr)[i] = (*arr)[pos+i]
	}
}

var clear = map[string]func() {
	"linux": func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	},
	"windows" : func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	},
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok { //if we defined a clear func for that platform:
		value()  //we execute it
	} else { //unsupported platform
		clear["linux"]()
	}
}