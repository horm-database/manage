package auth

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/horm-database/common/crypto"
	"github.com/horm-database/manage/srv/transport/web/head"
)

var SameRequestLock = new(sync.RWMutex)
var SameRequest = map[string]bool{}

func init() {
	go func() {
		for {
			SameRequestLock.Lock()
			SameRequest = map[string]bool{}
			SameRequestLock.Unlock()
			time.Sleep(time.Duration(30+rand.Intn(30)) * time.Second) // 0.5~1 分钟清空一次
		}
	}()
}

// SignSuccess 签名是否正确
func SignSuccess(header *head.WebReqHeader, token string) bool {
	if header.Userid == 0 {
		return false
	}

	if token == "" {
		return false
	}

	sign := crypto.MD5Str(fmt.Sprintf("%d%d%s%s%d%d%d%s%d", header.WorkspaceId,
		header.Userid, token, header.Version, header.RequestId, header.Timestamp,
		header.Timeout, header.Caller, header.AuthRand))

	if sign != header.Sign {
		return false
	}

	requestUniq := fmt.Sprintf("%d%s%d", header.Timestamp, header.Ip, header.AuthRand)

	SameRequestLock.Lock()
	isSame := SameRequest[requestUniq]
	isSame2 := SameRequest[sign]
	SameRequest[requestUniq] = true
	SameRequest[sign] = true
	SameRequestLock.Unlock()

	if isSame || isSame2 {
		return false
	}

	return true
}
