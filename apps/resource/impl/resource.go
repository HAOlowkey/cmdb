package impl

import (
	"context"

	"github.com/HAOlowkey/cmdb/apps/resource"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/sqlbuilder"
)

func (s *service) Search(ctx context.Context, req *resource.SearchRequest) (
	*resource.ResourceSet, error) {
	query := sqlbuilder.NewQuery(SQLQueryResource)

	if req.Keywords != "" {
		query.Where("name LIKE ? OR id = ? OR instance_id = ? OR private_ip LIKE ? OR public_ip LIKE ?",
			"%"+req.Keywords+"%",
			req.Keywords,
			req.Keywords,
			req.Keywords+"%",
			req.Keywords+"%",
		)
	}

	if req.Vendor != nil {
		query.Where("vendor = ?", req.Vendor)
	}

	querySQL, args := query.Order("sync_at").Desc().Limit(req.Page.ComputeOffset(), uint(req.Page.PageSize)).BuildQuery()
	s.log.Debugf("sql: %s", querySQL)

	queryStmt, err := s.db.Prepare(querySQL)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare query host error, %s", err.Error())
	}
	defer queryStmt.Close()

	rows, err := queryStmt.Query(args...)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	var (
		// ip的列表,  127.0.0.1,127.0.0.2  ==> ["127.0.0.1", "127.0.0.2"]
		publicIPList, privateIPList string
	)
	set := resource.NewResourceSet()
	for rows.Next() {
		ins := resource.NewDefaultResource()
		base := ins.Base
		info := ins.Information
		err := rows.Scan(
			&base.Id, &base.Vendor, &base.Region, &base.Zone, &base.CreateAt, &info.ExpireAt,
			&info.Category, &info.Type, &info.Name, &info.Description,
			&info.Status, &info.UpdateAt, &base.SyncAt, &info.SyncAccount,
			&publicIPList, &privateIPList, &info.PayType, &base.DescribeHash, &base.ResourceHash,
		)
		if err != nil {
			return nil, exception.NewInternalServerError("query host error, %s", err.Error())
		}
		info.LoadPrivateIPString(privateIPList)
		info.LoadPublicIPString(publicIPList)
		set.Add(ins)
	}

	// 获取total SELECT COUNT(*) FROMT t Where ....
	countSQL, args := query.BuildCount()
	countStmt, err := s.db.Prepare(countSQL)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}

	defer countStmt.Close()
	err = countStmt.QueryRow(args...).Scan(&set.Total)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}

	return set, nil
}
