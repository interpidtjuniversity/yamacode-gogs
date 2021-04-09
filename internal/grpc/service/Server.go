package service

import (
	"gogs.io/gogs/internal/grpc/service/serviceImpl"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Start() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//建立 gPRC 服务器，并注册服务
	s := grpc.NewServer()
	// start all grpc service
	serviceImpl.RegisterYaMaHubBranchServiceServer(s, &serviceImpl.BranchService{})
	serviceImpl.RegisterYaMaHubApplicationServiceServer(s, &serviceImpl.ApplicationService{})

	log.Println("Server run ...")
	//启动服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("fail to serve: %v", err)
	}
}
