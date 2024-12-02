package auth

import (
	"context"
	"time"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/model/table"
)

var CurrentWorkspaceID int

// InitWorkspaceID 初始化 workspace id
func InitWorkspaceID(ctx context.Context) error {
	workspaceInfo, err := table.GetCurrentWorkspace(ctx)
	if err != nil {
		return err
	}

	if workspaceInfo == nil || workspaceInfo.Id == 0 {
		return errs.New(errs.RetWebInitWorkspace, "init workspace failed")
	}

	CurrentWorkspaceID = workspaceInfo.Id
	return nil
}

// IsWorkspaceMember 是否空间成员
func IsWorkspaceMember(ctx context.Context, userid uint64, workspaceID int) error {
	if workspaceID != CurrentWorkspaceID {
		return errs.New(errs.RetWebNotIndicateSpace, "current workspace is not indicate workspace")
	}

	isNil, member, err := table.GetWorkspaceMemberByUser(ctx, workspaceID, userid)
	if err != nil {
		return err
	}

	if isNil || member == nil {
		return errs.New(errs.RetWebNotWorkspaceMember, "not workspace member, please apply")
	}

	if member.Status != consts.WorkspaceMemberStatusJoined &&
		member.Status != consts.WorkspaceMemberStatusRenewal { // 续期审批，可以在账号正常使用情况下申请
		return errs.New(errs.RetWebNotWorkspaceMember, "not workspace member, please apply")
	}

	if member.ExpireTime != 0 && int64(member.ExpireTime) < time.Now().Unix() {
		return errs.New(errs.RetWebWorkspaceMemberExpired, "workspace member permission has expired, please renewal")
	}

	return nil
}