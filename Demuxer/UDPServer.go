package Demuxer

import (
	"container/list"
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"net"
	"sync"
	"time"
)

// region Struct Definition
type UDPServer struct {
	port    int
	host    string
	clients *list.List
	syncMtx *sync.Mutex
	running bool
	conn    *net.UDPConn
	target  *net.UDPAddr
}

// endregion
// region Constructor
func NewUDPServer(host string, port int) *UDPServer {
	var udp = &UDPServer{
		port:    port,
		host:    host,
		syncMtx: &sync.Mutex{},
	}

	udp.Init()

	return udp
}

// endregion
// region BaseDemuxer Methods
func (f *UDPServer) Init() {
	f.clients = list.New()
}
func (f *UDPServer) Start() {
	f.syncMtx.Lock()
	f.running = true
	go f.loop()
	f.syncMtx.Unlock()
}
func (f *UDPServer) Stop() {
	f.syncMtx.Lock()
	f.running = false
	f.syncMtx.Unlock()
}
func (f *UDPServer) SendData(data []byte) {
	go func() {
		f.syncMtx.Lock()
		if f.conn != nil {
			_, err := f.conn.WriteToUDP(data, f.target)
			if err != nil {
				SLog.Error("Error sending payload to client: ", err)
			}
		}
		f.syncMtx.Unlock()
	}()
}
func (f *UDPServer) GetName() string {
	return "UDP Server"
}

func (f *UDPServer) isRunning() bool {
	f.syncMtx.Lock()
	defer f.syncMtx.Unlock()
	return f.running
}

// endregion
// region Loop Function
func (f *UDPServer) loop() {
	f.syncMtx.Lock()
	SLog.Info("Starting UDP Server at port %d", f.port)
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", 9001))

	if err != nil {
		SLog.Error("Error opening UDP Server Socket: %s\n", err)
		f.syncMtx.Unlock()
		return
	}

	ln, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		SLog.Error("Error opening UDP Server Socket: %s\n", err)
		f.syncMtx.Unlock()
		return
	}
	defer ln.Close()

	target, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", f.host, f.port))

	if err != nil {
		SLog.Error("Error opening UDP Server Socket: %s\n", err)
		f.syncMtx.Unlock()
		return
	}

	SLog.Info("UDP Server Started")
	f.conn = ln
	f.target = target
	f.syncMtx.Unlock()

	for f.isRunning() {
		time.Sleep(time.Millisecond * 100)
	}

}

// endregion
