package logic

import (
	"golang.org/x/net/context"
	"zg5/work/work07/server/proto/server"
)

type InitService struct {
	server.UnimplementedServerServer
}

func (c *InitService) UserReg(ctx context.Context, in *server.UserRegRequest) (*server.Response, error) {
	user := in.Username
	pwd := in.Password
	return &server.Response{
		Msg: user + pwd,
	}, nil
}
