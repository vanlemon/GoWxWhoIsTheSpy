package model

import (
	"context"
	"sync"
	"time"

	"github.com/jinzhu/gorm"

	"lmf.mortal.com/GoLogs"
)

// 用户信息模型
type User struct {
	ID         int64     `json:"id" gorm:"column:id"`                   // 自增ID
	OpenId     string    `json:"openid" gorm:"column:openid"`           // openid，和微信接口的字段名一致
	NickName   string    `json:"nickName" gorm:"column:nick_name"`      // 用户昵称
	Gender     int       `json:"gender" gorm:"column:gender"`           // 用户性别
	AvatarUrl  string    `json:"avatarUrl" gorm:"column:avatar_url"`    // 用户头像
	City       string    `json:"city" gorm:"column:city"`               // 用户城市
	Province   string    `json:"province" gorm:"column:province"`       // 用户省份
	Country    string    `json:"country" gorm:"column:country"`         // 用户国家
	Language   string    `json:"language" gorm:"column:language"`       // 用户语言
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"` // 创建时间
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"` // 更新时间
}

// 用户信息表名
func (User) TableName() string {
	return "user"
}

// 用户模型访问对象
type UserDao struct {
}

// 用户模型访问对象 - 单例模式
var userDao *UserDao
var userDaoOnce sync.Once

// 用户模型访问示例
func UserDaoInstance() *UserDao {
	userDaoOnce.Do(func() {
		userDao = &UserDao{}
	})
	return userDao
}

// 创建用户
func (d *UserDao) CreateUser(ctx context.Context, openId string, user *User) (err error) {
	logs.CtxInfo(ctx, "[Model MySQL CreateUser Req] req: %s, %#v", openId, user) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model MySQL CreateUser Resp] resp: %#v", err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	user.OpenId = openId
	user.CreateTime = time.Now()
	user.UpdateTime = time.Now()
	err = UserDB.Model(&User{}).Create(user).Error

	if err != nil {
		logs.CtxError(ctx, "[Model MySQL CreateUser] Create error: %#v", err)
		return err
	}
	/************************************ 核心逻辑 ****************************************/

	return nil
}

// 更新用户
func (d *UserDao) UpdateUser(ctx context.Context, openId string, user *User) (err error) {
	logs.CtxInfo(ctx, "[Model MySQL UpdateUser Req] req: %s, %#v", openId, user) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model MySQL UpdateUser Resp] resp: %#v", err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	user.UpdateTime = time.Now()
	db := UserDB.Model(&User{}).Where("openid = ?", openId).Updates(user)

	//更新 db.RowsAffected == 0 时 err = nil，此时可能数据不存在或者数据未更新
	if db.RowsAffected == 0 {
		logs.CtxWarn(ctx, "[Model MySQL UpdateUser] Updates rows affect is zero")
	}

	err = db.Error
	if err != nil {
		logs.CtxError(ctx, "[Model MySQL UpdateUser] Updates error: %#v", err)
		return err
	}
	/************************************ 核心逻辑 ****************************************/

	return nil
}

// 查询用户，查不到返回 nil，需要判断是否查到
func (d *UserDao) QueryUser(ctx context.Context, openid string) (user *User, err error) {
	logs.CtxInfo(ctx, "[Model MySQL QueryUser Req] req: %s", openid) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model MySQL QueryUser Resp] resp: %#v, %#v", user, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	user = &User{} // 定义接收结构体
	err = UserDB.Model(&User{}).Where("openid = ?", openid).First(user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		logs.CtxError(ctx, "[Model MySQL QueryUser] Query error: %#v", err)
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		logs.CtxWarn(ctx, "[Model MySQL QueryUser] Query record not found")
		return nil, nil
	}
	/************************************ 核心逻辑 ****************************************/

	return user, nil
}
