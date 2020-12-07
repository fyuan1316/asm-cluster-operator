package error

import (
	"errors"
	"time"
)

var (
	ErrNeedRetry = errors.New("need retry")
)

const (
	ReconcileAfterDuration            = time.Second * 2
	HealthCheckReconcileAfterDuration = time.Second * 10 //300
)
