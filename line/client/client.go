package client

import (
	"github.com/pkg/errors"
	pb "github.com/utahta/momoclo-channel/line/protos"
	"github.com/utahta/momoclo-channel/log"
	"google.golang.org/grpc"
)

type LineClient struct {
	conn *grpc.ClientConn
	pb.LineClient

	Log log.Logger
}

func Dial(address string, opts ...grpc.DialOption) (*LineClient, error) {
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "did not connect. address:%s", address)
	}
	return &LineClient{conn: conn, LineClient: pb.NewLineClient(conn), Log: log.NewSilentLogger()}, nil
}

func (c *LineClient) Close() {
	if c.conn == nil {
		return
	}
	if err := c.conn.Close(); err != nil {
		c.Log.Errorf("Failed to close connection. error:%v", err)
	}
}
