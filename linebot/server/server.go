package server

import (
	"net"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
	pb "github.com/utahta/momoclo-channel/linebot/protos"
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

func (s *notificationServer) NotifyMessage(c context.Context, r *pb.NotifyMessageRequest) (*pb.NotifyMessageResponse, error) {
	s.Log.Infof("start notify message. params:%#v", r)

	mm := s.Client.NewMultipleMessage()
	mm.AddText(r.Text)
	for _, url := range r.ImageUrls {
		mm.AddImage(url, url)
	}
	_, err := mm.Send(r.To)
	if err != nil {
		s.Log.Errorf("Failed to send message. error:%v", err)
	}

	s.Log.Info("end notify message.")
	return &pb.NotifyMessageResponse{}, nil
}

func (s *notificationServer) Run(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrapf(err, "failed to listen. port:%s", port)
	}

	gs := grpc.NewServer()
	pb.RegisterLineBotServer(gs, s)
	gs.Serve(lis)
	return nil
}
