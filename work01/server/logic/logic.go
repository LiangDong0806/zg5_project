package logic

import (
	"context"
	"github.com/go-errors/errors"
	"net/http"
	"zg5/work01/server/models"
	server "zg5/work01/server/proto"
)

func (r *ServerRpc) UserReg(ctx context.Context, in *server.UserRegRequest) (*server.Response, error) {
	res, err := models.GetUserByUsername(in.Username)
	if err != nil {
		return &server.Response{}, err
	}
	if res.Password != in.Password {
		return &server.Response{}, errors.New("密码错误")
	}
	return &server.Response{
		Msg: "登录成功",
	}, nil
}

func (c *ServerRpc) CreateCustomer(ctx context.Context, in *server.CreateOrderRequest) (*server.Response, error) {

	order := map[string]interface{}{
		"OrderId":     in.OrderId,
		"OrderName":   in.OrderName,
		"OrderPhone":  in.OrderPhone,
		"OrderNum":    in.OrderNum,
		"OrderStatus": in.OrderStatus,
	}

	err := models.EscIns(order)
	if err != nil {
		return &server.Response{}, errors.New("订单添加失败")
	}
	return &server.Response{
		Code: http.StatusOK,
		Msg:  "订单添加成功",
	}, nil
}
