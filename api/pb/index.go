package pb

type IndexTableListRequest struct {
	Page int `json:"page"` // 分页
	Size int `json:"size"` // 每页大小
}

type IndexTableListResponse struct {
	Total     uint64       `json:"total"`      // 总数
	TotalPage uint32       `json:"total_page"` // 总页数
	Page      int          `json:"page"`       // 分页
	Size      int          `json:"size"`       // 每页大小
	Tables    []*TableBase `json:"tables"`     // 用户列表
}

type CollectTableListRequest struct {
	Page int `json:"page"` // 分页
	Size int `json:"size"` // 每页大小
}

type CollectTableRequest struct {
	TableID int `json:"table_id"` // 表id
	Status  int `json:"status"`   // 1-收藏 2-取关
}
