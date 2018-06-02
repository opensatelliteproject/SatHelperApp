package main

import (
	"github.com/foize/go.fifo"
	"log"
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

