package cconst

import (
	"github.com/bmizerany/assert"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"
	"testing"
)

func TestRedisLockDuplicate(t *testing.T) {
	// 错误是可以比较相同的
	assert.T(t, RedisLockDuplicate == RedisLockDuplicate)
	// 新建的错误是不可以比较相同的
	assert.T(t, RedisLockDuplicate != util.NewErrf("RedisLockDuplicate"))
}
