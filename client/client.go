package client

import (
	kc "github.com/infraboard/keyauth/client"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"

	"github.com/HAOlowkey/cmdb/apps/book"
	"github.com/HAOlowkey/cmdb/apps/host"
	"github.com/HAOlowkey/cmdb/apps/resource"
	"github.com/HAOlowkey/cmdb/apps/secret"
	"github.com/HAOlowkey/cmdb/apps/task"
)

var (
	client *ClientSet
)

// SetGlobal todo
func SetGlobal(cli *ClientSet) {
	client = cli
}

// C Global
func C() *ClientSet {
	return client
}

// NewClient todo
func NewClient(conf *kc.Config) (*ClientSet, error) {
	zap.DevelopmentSetup()
	log := zap.L()

	conn, err := grpc.Dial(conf.Address(), grpc.WithInsecure(), grpc.WithPerRPCCredentials(conf.Authentication))
	if err != nil {
		return nil, err
	}

	return &ClientSet{
		conn: conn,
		log:  log,
	}, nil
}

// Client 客户端
type ClientSet struct {
	conn *grpc.ClientConn
	log  logger.Logger
}

// Book服务的SDK
func (c *ClientSet) Book() book.ServiceClient {
	return book.NewServiceClient(c.conn)
}

// Resource todo
func (c *ClientSet) Resource() resource.ServiceClient {
	return resource.NewServiceClient(c.conn)
}

// Host todos
func (c *ClientSet) Host() host.ServiceClient {
	return host.NewServiceClient(c.conn)
}

// Host todos
func (c *ClientSet) Secret() secret.ServiceClient {
	return secret.NewServiceClient(c.conn)
}

func (c *ClientSet) Task() task.ServiceClient {
	return task.NewServiceClient(c.conn)
}
