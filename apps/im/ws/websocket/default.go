package websocket

import (
	"math"
	"time"
)

const (
	defaultMaxConnectionIdle = time.Duration(math.MaxInt64)
	// defaultMaxConnectionIdle = 5 * time.Second
	defaultAckTimeout        = 30 * time.Second
	defaultConcurrency =10
)

