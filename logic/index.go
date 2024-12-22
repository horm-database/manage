// Copyright (c) 2024 The horm-database Authors (such as CaoHao <18500482693@163.com>). All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
