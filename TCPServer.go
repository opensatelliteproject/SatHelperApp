package main

import (
	"container/list"
	"fmt"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"github.com/prometheus/common/log"
	"net"
	"sync"
)

// region Struct Definition
type TCPServer struct {
	port    int
	host    string
	clients *list.List
	syncMtx *sync.Mutex
	running bool
}

// endregion
// region Constructor
func NewTCPServer(host string, port int) *TCPServer {
	var tcp = &TCPServer{
		port:    port,
		host:    host,
		syncMtx: &sync.Mutex{},
	}

	tcp.Init()

	return tcp
}

// endregion
// region BaseDemuxer Methods
func (f *TCPServer) Init() {
	f.clients = list.New()
}
func (f *TCPServer) Start() {
	f.running = true
	go f.loop()
}
func (f *TCPServer) Stop() {
	f.running = false
}
func (f *TCPServer) SendData(data []byte) {
	go func() {
		f.syncMtx.Lock()
		var next *list.Element
		for e := f.clients.Front(); e != nil; e = next {
			client := e.Value.(net.Conn)
			n, err := client.Write(data)
			if n != len(data) || err != nil {
				SLog.Error("%s", err)
				next = e.Next()
				f.clients.Remove(e)
				SLog.Info("Client disconnected %s", client.RemoteAddr())
			}
		}
		f.syncMtx.Unlock()
	}()
}
func (f *TCPServer) GetName() string {
	return "TCP Server"
}

// endregion
// region Loop Function
func (f *TCPServer) loop() {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", f.host, f.port))
	if err != nil {
		SLog.Error("Error opening TCP Server Socket: %s\n", err)
		return
	}
	for f.running {
		conn, err := ln.Accept()
		if err != nil {
			log.Error(err)
		} else {
			f.syncMtx.Lock()
			f.clients.PushBack(conn)
			f.syncMtx.Unlock()
			SLog.Info("Client connected from %s", conn.RemoteAddr())
			// go f.handleConnection(conn)
		}
	}

}

//func (f *TCPServer) handleConnection(conn net.Conn) {
//	// TODO: Needed?
//}

// endregion
