package snowflake

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	MaxCountID = 1000000 // 10万，43毫秒左右，100万,430毫秒左右；1000万，4300毫秒左右
)

func TestSnowFlake_GetID(t *testing.T) {
	begin := time.Now()
	for i:=0;i<MaxCountID;i++{
		GetID()
	}

	fmt.Println(fmt.Sprintf("time cost:%d", time.Now().Sub(begin).Milliseconds()))
}

// 批量取跟单个取耗时差不多
func TestSnowFlake_GetIDs(t *testing.T) {
	begin := time.Now()
	GetIDs(MaxCountID)

	fmt.Println(fmt.Sprintf("time cost:%d", time.Now().Sub(begin).Milliseconds()))
}

func BenchmarkSnowFlake_GetID(b *testing.B) {
	for i:=0;i< b.N;i++ {
		idSet := make(map[int64]bool)
		for j:=0;j<MaxCountID;j++{
			id, err := GetID()
			assert.Equal(b, nil, err)
			idSet[id] = true
		}
		assert.Equal(b, MaxCountID, len(idSet))
	}
}

// 压测的时候批量还是快一些，大约快一倍
func BenchmarkSnowFlake_GetIDs(b *testing.B) {
	for i:=0;i< b.N;i++ {
		ids, err := GetIDs(MaxCountID)
		assert.Equal(b, nil, err)

		idSet := make(map[int64]bool)
		for i:=range ids{
			idSet[ids[i]] = true
		}
		assert.Equal(b, MaxCountID, len(idSet))
	}
}