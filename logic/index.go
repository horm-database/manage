package logic

import (
	"context"

	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/model/table"
)

func IndexTableList(ctx context.Context, req *pb.IndexTableListRequest) (*pb.IndexTableListResponse, error) {
	pageInfo, tables, err := table.GetIndexTable(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	ret := pb.IndexTableListResponse{
		Total:     pageInfo.Total,
		TotalPage: pageInfo.TotalPage,
		Page:      pageInfo.Page,
		Size:      pageInfo.Size,
		Tables:    []*pb.TableBase{},
	}

	for _, v := range tables {
		ret.Tables = append(ret.Tables, GetTableBase(v))
	}

	return &ret, nil
}

func CollectTableList(ctx context.Context, userid uint64,
	req *pb.CollectTableListRequest) (*pb.IndexTableListResponse, error) {
	pageInfo, collectTables, err := table.GetCollectTable(ctx, userid, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	ret := pb.IndexTableListResponse{
		Total:     pageInfo.Total,
		TotalPage: pageInfo.TotalPage,
		Page:      pageInfo.Page,
		Size:      pageInfo.Size,
		Tables:    []*pb.TableBase{},
	}

	tableIds := table.GetCollectTablesID(collectTables)

	tables, err := table.GetTableByIds(ctx, tableIds)
	if err != nil {
		return nil, err
	}

	tableMap := table.TablesToMap(tables)

	for _, ct := range collectTables {
		tbl := tableMap[ct.TableID]
		if tbl == nil {
			continue
		}

		ret.Tables = append(ret.Tables, GetTableBase(tbl))
	}

	return &ret, nil
}

func CollectTable(ctx context.Context, userid uint64,
	req *pb.CollectTableRequest) (err error) {
	if req.Status == consts.CollectTableAdd {
		err = table.AddCollectTable(ctx, userid, req.TableID)
	} else {
		err = table.DelCollectTable(ctx, userid, req.TableID)
	}

	return err
}
