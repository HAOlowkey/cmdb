package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"

	"github.com/HAOlowkey/cmdb/apps/resource"
	"github.com/HAOlowkey/cmdb/conf"
)

var (
	// Service 服务实例
	svr = &service{}
)

type service struct {
	db       *sql.DB
	log      logger.Logger
	resource resource.ServiceServer
	resource.UnimplementedServiceServer
}

func (s *service) Config() error {
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}

	s.log = zap.L().Named(s.Name())
	s.db = db
	s.resource = app.GetGrpcApp(resource.AppName).(resource.ServiceServer)
	return nil
}

func (s *service) Name() string {
	return resource.AppName
}

func (s *service) Registry(server *grpc.Server) {
	resource.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}
