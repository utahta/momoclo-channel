package server

import (
	"net"
	"fmt"

	pb "github.com/utahta/momoclo-channel/line/protos"
	"github.com/utahta/momoclo-channel/log"
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/pkg/errors"
)

type notificationServer struct {
	Client *linebot.Client
	Log log.Logger
}

func New(channelID int64, channelSecret, channelMID string) (*notificationServer, error) {
	client, err := linebot.NewClient(channelID, channelSecret, channelMID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init linebot client")
	}
	return &notificationServer{Client: client, Log: log.NewBasicLogger()}, nil
}

func (s *notificationServer) NotifyChannel(c context.Context, r *pb.NotifyChannelRequest) (*pb.NotifyChannelResponse, error) {
	mm := s.Client.NewMultipleMessage()
	mm.AddText(fmt.Sprintf("%s\n%s\n%s", r.Title, r.Item.Title, r.Item.Url))
	for _, img := range r.Item.Images {
		mm.AddImage(img.Url, img.Url)
	}
	_, err := mm.Send(r.To)
	if err != nil {
		s.Log.Errorf("Failed to send channel. error:%v", err)
	}
	return &pb.NotifyChannelResponse{}, nil
}

func (s *notificationServer) AppendUser(c context.Context, r *pb.AppendUserRequest) (*pb.AppendUserResponse, error) {
	_, err := s.Client.SendText([]string{r.To}, "通知ノフ設定オンにしました（・Θ・）")
	if err != nil {
		s.Log.Errorf("failed to send text. error:%v", err)
	}
	return &pb.AppendUserResponse{}, nil
}

func (s *notificationServer) DeleteUser(c context.Context, r *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_, err := s.Client.SendText([]string{r.To}, "通知ノフ設定オフにしました（・Θ・）")
	if err != nil {
		s.Log.Errorf("failed to send text. error:%v", err)
	}
	return &pb.DeleteUserResponse{}, nil
}

func (s *notificationServer) Run(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrapf(err, "failed to listen. port:%s", port)
	}

	gs := grpc.NewServer()
	pb.RegisterLineServer(gs, s)
	gs.Serve(lis)
	return nil
}
