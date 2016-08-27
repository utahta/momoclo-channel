package client

import (
	pb "github.com/utahta/momoclo-channel/line/protos"
	"google.golang.org/grpc"
	"github.com/pkg/errors"
)

type LineClient struct {
	conn *grpc.ClientConn
	pb.LineClient
}

func Dial(address string) (*LineClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "did not connect. address:%s", address)
	}
	return &LineClient{ conn: conn, LineClient: pb.NewLineClient(conn) }, nil
}

func (c *LineClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
