package server

import (
	"fmt"
	"net"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
	pb "github.com/utahta/momoclo-channel/line/protos"
	"github.com/utahta/momoclo-channel/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type notificationServer struct {
	Client *linebot.Client
	Log    log.Logger
}

func New(channelID int64, channelSecret, channelMID string) (*notificationServer, error) {
	client, err := linebot.NewClient(channelID, channelSecret, channelMID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init linebot client")
	}
	return &notificationServer{Client: client, Log: log.NewBasicLogger()}, nil
}

func (s *notificationServer) NotifyChannel(c context.Context, r *pb.NotifyChannelRequest) (*pb.NotifyChannelResponse, error) {
	s.Log.Infof("start notify channel. params:%#v", r)

	mm := s.Client.NewMultipleMessage()
	mm.AddText(fmt.Sprintf("%s\n%s\n%s", r.Title, r.Item.Title, r.Item.Url))
	for _, img := range r.Item.Images {
		mm.AddImage(img.Url, img.Url)
	}
	_, err := mm.Send(r.To)
	if err != nil {
		s.Log.Errorf("Failed to send channel. error:%v", err)
	}

	s.Log.Info("end notify channel.")
	return &pb.NotifyChannelResponse{}, nil
}

func (s *notificationServer) NotifyAppendUser(c context.Context, r *pb.NotifyAppendUserRequest) (*pb.NotifyAppendUserResponse, error) {
	s.Log.Infof("start append user. params:%#v", r)

	_, err := s.Client.SendText([]string{r.To}, "通知ノフ設定オンにしました（・Θ・）")
	if err != nil {
		s.Log.Errorf("failed to send text. error:%v", err)
	}

	s.Log.Info("end append user.")
	return &pb.NotifyAppendUserResponse{}, nil
}

func (s *notificationServer) NotifyDeleteUser(c context.Context, r *pb.NotifyDeleteUserRequest) (*pb.NotifyDeleteUserResponse, error) {
	s.Log.Infof("start delete user. params:%#v", r)

	_, err := s.Client.SendText([]string{r.To}, "通知ノフ設定オフにしました（・Θ・）")
	if err != nil {
		s.Log.Errorf("failed to send text. error:%v", err)
	}

	s.Log.Infof("end delete user.")
	return &pb.NotifyDeleteUserResponse{}, nil
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
