package service

import "context"

type HelloService struct{}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (oSelf *HelloService) TestHello(ctx context.Context, sReq string, sRsp *string) (int32, error) {
	*sRsp = "Hello " + sReq
	return 0, nil
}
