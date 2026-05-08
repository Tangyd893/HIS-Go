// Package snowflake 雪花算法分布式ID生成器
package snowflake

import (
	"sync"
	"time"
)

const (
	epoch            = int64(1704038400000) // 起始时间戳 2024-01-01 00:00:00 UTC（毫秒）
	workerIDBits     = uint(5)
	datacenterIDBits = uint(5)
	sequenceBits     = uint(12)

	maxWorkerID     = int64(-1 ^ (-1 << workerIDBits))
	maxDatacenterID = int64(-1 ^ (-1 << datacenterIDBits))
	maxSequence     = int64(-1 ^ (-1 << sequenceBits))

	workerIDShift     = sequenceBits
	datacenterIDShift = sequenceBits + workerIDBits
	timestampShift    = sequenceBits + workerIDBits + datacenterIDBits
)

type Snowflake struct {
	mu            sync.Mutex
	epoch         int64
	workerID      int64
	datacenterID  int64
	sequence      int64
	lastTimestamp int64
}

func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, ErrInvalidWorkerID
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, ErrInvalidDatacenterID
	}
	return &Snowflake{
		epoch:        epoch,
		workerID:     workerID,
		datacenterID: datacenterID,
	}, nil
}

func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()
	if now < s.lastTimestamp {
		return 0, ErrClockBackwards
	}

	if now == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for now <= s.lastTimestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = now

	id := ((now - s.epoch) << timestampShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id, nil
}

func (s *Snowflake) NextString() (string, error) {
	id, err := s.NextID()
	if err != nil {
		return "", err
	}
	return formatStringID(id), nil
}

func formatStringID(id int64) string {
	const digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if id == 0 {
		return "0"
	}
	var result [11]byte
	var i int
	for i = len(result) - 1; id > 0 && i >= 0; i-- {
		result[i] = digits[id%62]
		id /= 62
	}
	return string(result[i+1:])
}
