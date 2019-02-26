package Structs

import (
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"os"
	"sync"
	"time"
)

const imageTimeout = time.Minute * 60 // 1 hour timeout

type MultiSegmentImage struct {
	SubProductId int
	ImageID      int
	Name         string
	Files        []string
	MaxSegments  int

	CreatedAt time.Time

	FirstSegmentHeader   *XRIT.Header
	FirstSegmentFilename string

	lock sync.Mutex
}

func StringIndexOf(v string, a []string) int {
	for i, vo := range a {
		if vo == v {
			return i
		}
	}

	return -1
}

func MakeMultiSegmentImage(Name string, SubProductId, ImageID int) *MultiSegmentImage {
	return &MultiSegmentImage{
		Name:               Name,
		SubProductId:       SubProductId,
		ImageID:            ImageID,
		CreatedAt:          time.Now(),
		MaxSegments:        -1,
		Files:              make([]string, 0),
		lock:               sync.Mutex{},
		FirstSegmentHeader: nil,
	}
}

func (msi *MultiSegmentImage) Expired() bool {
	msi.lock.Lock()
	defer msi.lock.Unlock()

	return time.Since(msi.CreatedAt) > imageTimeout
}

func (msi *MultiSegmentImage) PutSegment(filename string, xh *XRIT.Header) {
	msi.lock.Lock()
	defer msi.lock.Unlock()

	if xh.SegmentIdentificationHeader == nil {
		SLog.Error("No Segment Identification Header for segment!")
		return
	}

	if msi.MaxSegments == -1 {
		msi.MaxSegments = int(xh.SegmentIdentificationHeader.MaxSegments)
	}

	if StringIndexOf(filename, msi.Files) == -1 {
		msi.Files = append(msi.Files, filename)
	}

	if (xh.SegmentIdentificationHeader.Sequence == 0 || xh.SegmentIdentificationHeader.Sequence == 1) && msi.FirstSegmentHeader == nil {
		msi.FirstSegmentHeader = xh
		msi.FirstSegmentFilename = filename
	}
}

func (msi *MultiSegmentImage) Done() bool {
	msi.lock.Lock()
	defer msi.lock.Unlock()

	return len(msi.Files) == msi.MaxSegments
}

func (msi *MultiSegmentImage) Purge() {
	for _, v := range msi.Files {
		SLog.Debug("Removing %s", v)
		err := os.Remove(v)
		if err != nil {
			SLog.Error("Error erasing %s: %s", v, err)
		}
	}
}
