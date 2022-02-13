package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"

	"github.com/HAOlowkey/cmdb/apps/host"
	"github.com/HAOlowkey/cmdb/apps/secret"
	"github.com/HAOlowkey/cmdb/apps/task"
	"github.com/HAOlowkey/cmdb/conf"
)

var (
	// Service 服务实例
	svr = &service{}
)

type service struct {
	db     *sql.DB
	log    logger.Logger
	secret secret.ServiceServer
	host   host.ServiceServer
	task.UnimplementedServiceServer
}

func (s *service) Config() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}

	s.log = zap.L().Named(s.Name())
	s.db = db
	s.secret = app.GetGrpcApp(secret.AppName).(secret.ServiceServer)
	s.host = app.GetGrpcApp(host.AppName).(host.ServiceServer)

	return nil
}

func (s *service) Name() string {
	return task.AppName
}

func (s *service) Registry(server *grpc.Server) {
	task.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}
