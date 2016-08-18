package server

import (
	"net"
	"log"

	pb "github.com/utahta/momoclo-channel/grpc/line/protos"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/pkg/errors"
)

type lineServer struct{}

func (s *lineServer) NotifyChannel(c context.Context, r *pb.NotifyChannelRequest) (*pb.NotifyChannelResponse, error) {
	return &pb.NotifyChannelResponse{}, nil
}

func (s *lineServer) AppendUser(c context.Context, r *pb.AppendUserRequest) (*pb.AppendUserResponse, error) {
	log.Printf("append user. request:%#v", r)
	return &pb.AppendUserResponse{}, nil
}

func (s *lineServer) DeleteUser(c context.Context, r *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	return &pb.DeleteUserResponse{}, nil
}

func Run(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrapf(err, "failed to listen. port:%s", port)
	}

	s := grpc.NewServer()
	pb.RegisterLineServer(s, &lineServer{})
	s.Serve(lis)
	return nil
}
