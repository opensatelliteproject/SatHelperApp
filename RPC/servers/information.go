package servers

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opensatelliteproject/SatHelperApp/RPC/sathelperapp"
)

type InformationFetcher interface {
	GetStatistics() (*sathelperapp.StatData, error)
	GetConsoleLines() (*sathelperapp.ConsoleData, error)
}

type informationServer struct {
	infoSource InformationFetcher
}

func MakeInformationServer(source InformationFetcher) sathelperapp.InformationServer {
	return &informationServer{
		infoSource: source,
	}
}

func (s *informationServer) GetStatistics(context.Context, *empty.Empty) (*sathelperapp.StatData, error) {
	return s.infoSource.GetStatistics()
}

func (s *informationServer) GetConsoleLines(context.Context, *empty.Empty) (*sathelperapp.ConsoleData, error) {
	return s.infoSource.GetConsoleLines()
}
