package RPC

import (
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp/Logger"
	"github.com/opensatelliteproject/SatHelperApp/RPC/sathelperapp"
	"github.com/opensatelliteproject/SatHelperApp/RPC/servers"
	"google.golang.org/grpc"
	"net"
)

type DataSource interface {
	servers.InformationFetcher
}

type Server struct {
	infoServer sathelperapp.InformationServer
	grpcServer *grpc.Server
}

func MakeRPCServer(source DataSource) *Server {
	return &Server{
		infoServer: servers.MakeInformationServer(source),
	}
}

func (s *Server) Listen(address string) error {
	if s.grpcServer != nil {
		return fmt.Errorf("server already runing")
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.grpcServer = grpc.NewServer()
	sathelperapp.RegisterInformationServer(s.grpcServer, s.infoServer)
	go s.serve(lis)
	return nil
}

func (s *Server) serve(conn net.Listener) {
	err := s.grpcServer.Serve(conn)
	if err != nil {
		SLog.Error("RPC Error: %s", err)
	}
	s.Stop()
}

func (s *Server) Stop() {
	if s.grpcServer == nil {
		SLog.Error("RPC Server Already stopped")
		return
	}
	SLog.Info("Stopping RPC Server")
	s.grpcServer.Stop()
	s.grpcServer = nil
}
