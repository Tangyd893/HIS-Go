// Package snowflake 雪花算法错误定义
package snowflake

import "errors"

var (
	ErrInvalidWorkerID     = errors.New("worker ID 超出范围")
	ErrInvalidDatacenterID = errors.New("datacenter ID 超出范围")
	ErrClockBackwards      = errors.New("时钟回拨，拒绝生成ID")
)
