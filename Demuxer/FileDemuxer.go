package Demuxer

import (
	"container/list"
	"os"
	"sync"
)

// region Struct Definition
type FileDemuxer struct {
	filename string
	clients  *list.List
	syncMtx  *sync.Mutex
	running  bool
	handle   *os.File
}

// endregion
// region Constructor
func NewFileDemuxer(filename string) *FileDemuxer {
	return &FileDemuxer{
		filename: filename,
		syncMtx:  &sync.Mutex{},
	}
}

// endregion
// region BaseDemuxer Methods
func (f *FileDemuxer) Init() {
	f.clients = list.New()
}
func (f *FileDemuxer) Start() {
	f.running = true
	var err error
	f.handle, err = os.Create(f.filename)

	if err != nil {
		panic(err)
	}
}
func (f *FileDemuxer) Stop() {
	f.running = false
	if f.handle != nil {
		f.handle.Close()
		f.handle = nil
	}
}
func (f *FileDemuxer) SendFrame(frame []byte) {
	_, err := f.handle.Write(frame)
	if err != nil {
		panic(err)
	}
}
func (f *FileDemuxer) GetName() string {
	return "File"
}

// endregion
