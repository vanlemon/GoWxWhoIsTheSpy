package model

import (
	"github.com/bmizerany/assert"
	"lmf.mortal.com/GoLogs"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"
	"testing"
)

func TestUserDao_Create_Update_QueryUser(t *testing.T) {
	ctx := logs.TestCtx()

	// 生成随机 openid
	openId := util.GenUUID(4)

	user := &User{OpenId: openId}

	// 测试创建用户
	err := UserDaoInstance().CreateUser(ctx, openId, user)
	assert.Equal(t, err, nil)

	// 测试查询用户
	user, err = UserDaoInstance().QueryUser(ctx, openId)
	assert.Equal(t, err, nil)
	assert.Equal(t, user.OpenId, openId)
	assert.NotEqual(t, user.CreateTime, nil)
	assert.NotEqual(t, user.UpdateTime, nil)

	// 测试更新用户
	user.Gender = 1
	err = UserDaoInstance().UpdateUser(ctx, openId, user)
	assert.Equal(t, err, nil)

	// 测试查询用户
	user, err = UserDaoInstance().QueryUser(ctx, openId)
	assert.Equal(t, err, nil)
	assert.Equal(t, user.Gender, 1)
}
