package consensus

import "time"

type Customize struct {
	// 前回の Commit から次の Commit までの間隔の最大値
	MaxWaitngCommitInterval time.Duration
}

