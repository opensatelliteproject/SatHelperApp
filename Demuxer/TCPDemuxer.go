package Demuxer

import (
	"net"
	"fmt"
	"github.com/prometheus/common/log"
	"sync"
	"container/list"
)
// region Struct Definition
type TCPServerDemuxer struct {
	port int
	host string
	connection net.Listener
	clients *list.List
	syncMtx *sync.Mutex
	running bool
}
// endregion
// region Constructor
func NewTCPDemuxer(host string, port int) *TCPServerDemuxer {
	return &TCPServerDemuxer{
		port: port,
		host: host,
		syncMtx: &sync.Mutex{},
	}
}
// endregion
// region BaseDemuxer Methods
func (f *TCPServerDemuxer) Init() {
	f.clients = list.New()
}
func (f *TCPServerDemuxer) Start() {
	f.running = true
	go f.loop()
}
func (f *TCPServerDemuxer) Stop() {
	f.running = false
}
func (f *TCPServerDemuxer) SendFrame(frame []byte) {
	go func() {
		f.syncMtx.Lock()
		var next *list.Element
		for e := f.clients.Front(); e != nil; e = next {
			client := e.Value.(net.Conn)
			n, err := client.Write(frame)
			if n != len(frame) || err != nil {
				log.Error(err)
				next = e.Next()
				f.clients.Remove(e)
				log.Infof("Client disconnected %s", client.RemoteAddr())
			}
		}
		f.syncMtx.Unlock()
	}()
}
// endregion
// region Loop Function
func (f *TCPServerDemuxer) loop() {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", f.host, f.port))
	if err != nil {
		log.Fatalf("Error opening TCP Server Socket: %s\n", err)
	}
	for f.running {
		conn, err := ln.Accept()
		if err != nil {

		} else {
			f.syncMtx.Lock()
			f.clients.PushBack(conn)
			f.syncMtx.Unlock()
			log.Infof("Client connected from %s", conn.RemoteAddr())
			// go f.handleConnection(conn)
		}
	}

}
func (f *TCPServerDemuxer) handleConnection(conn net.Conn) {
	// TODO: Needed?
}
// endregion
