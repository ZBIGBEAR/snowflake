/*
雪花算法，生成64位的数字
第一位：符号位，不用
后面41位：毫秒级时间戳
后面10位：机器id，可以部署1024个机器
后12位：序号，最大为4096.即一毫秒内最多可以生成4096个不重复的id
*/

package snowflake

import (
	"sync"
	"time"
)

var (
	snowFlake iSnowFlake
	once      sync.Once
)

type iSnowFlake interface {
	GetID() (int64, error)
	GetIDs(count int) ([]int64, error)
}

type SnowFlake struct {
	serverID      int64
	lastTimeStamp int64
	sequence      int
	mu            sync.Mutex
}

func InitSnowFlake(serverID int64) {
	initSnowFlake(serverID)
}

func initSnowFlake(serverID int64) {
	once.Do(func() {
		snowFlake = &SnowFlake{
			serverID:      serverID,
			lastTimeStamp: int64(time.Now().UnixNano() / 1000000),
		}
	})
}

func (s *SnowFlake) GetID() (int64, error) {
	ids, err := s.GetIDs(1)
	if err != nil {
		return 0, err
	}

	return ids[0], nil
}

func (s *SnowFlake) GetIDs(count int) ([]int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	curTimeStamp := s.getCurTimeStamp()
	if curTimeStamp > s.lastTimeStamp {
		// 不同毫秒内，重置计数器
		s.sequence = 0
	}

	if count > NB4095 {
		return nil, CountErr
	}

	if count > (NB4095 - s.sequence) {
		time.Sleep(time.Millisecond)
		curTimeStamp = s.getCurTimeStamp()
		s.sequence = 0
	}

	var ids []int64
	for i := 0; i < count; i++ {
		ids = append(ids, s.genID(curTimeStamp, int64(s.sequence)))
		if s.sequence > NB4095 {
			return nil, SeqStack
		}
		s.sequence++
	}
	s.lastTimeStamp = curTimeStamp
	return ids, nil
}

func (s *SnowFlake) genID(curTimeStamp int64, seq int64) int64 {
	return (((1 << 63) - 1) & (curTimeStamp << 22)) | (((1 << 22) - 1) & (s.serverID << 12)) | (((1 << 12) - 1) & seq)
}

func (s *SnowFlake) getCurTimeStamp() int64 {
	for {
		currTimeStamp := int64(time.Now().UnixNano() / 1000000)
		if currTimeStamp < s.lastTimeStamp {
			time.Sleep(time.Millisecond)
		} else {
			return currTimeStamp
		}
	}
}
