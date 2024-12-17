package table

import (
	"context"
	"time"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
)

func GetCollectTable(ctx context.Context, userid uint64, page, size int) (*proto.Detail, []*TblCollectTable, error) {
	pageRet := proto.Detail{}

	collectTables := []*TblCollectTable{}

	where := horm.Where{}
	where["userid"] = userid

	_, err := GetTableORM("tbl_collect_table").
		FindAll(where).
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageRet, &collectTables)

	return &pageRet, collectTables, err
}

func DelCollectTable(ctx context.Context, userid uint64, tableID int) error {
	where := horm.Where{}
	where["userid"] = userid
	where["table_id"] = tableID

	_, err := GetTableORM("tbl_collect_table").Delete(where).Exec(ctx)

	return err
}

func AddCollectTable(ctx context.Context, userid uint64, tableID int) error {
	data := horm.Map{}
	data["userid"] = userid
	data["table_id"] = tableID
	data["updated_at"] = time.Now()

	_, err := GetTableORM("tbl_collect_table").Replace(data).Exec(ctx)

	return err
}

///////////////////////////////// function /////////////////////////////////////////

func GetCollectTablesID(collectTables []*TblCollectTable) []int {
	var tableIDs = []int{}
	for _, tmp := range collectTables {
		tableIDs = append(tableIDs, tmp.TableID)
	}

	return tableIDs
}
