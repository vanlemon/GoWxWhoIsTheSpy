package service

import (
	"context"
	"encoding/json"
	"lmf.mortal.com/GoWxWhoIsTheSpy/cconst"

	"sync"

	logs "lmf.mortal.com/GoLogs"

	"lmf.mortal.com/GoWxWhoIsTheSpy/config"
	"lmf.mortal.com/GoWxWhoIsTheSpy/model"
	"lmf.mortal.com/GoWxWhoIsTheSpy/util"
)

// 用户服务
type UserService struct {
	s3idDao *model.Session3idDao // 访问用户临时 ID
	userDao *model.UserDao       // 访问用户信息
}

var userServiceInstance *UserService  // 用户服务实例
var userServiceInstanceOnce sync.Once // 单例模式

// 获取用户服务单例
func UserServiceInstance() *UserService {
	userServiceInstanceOnce.Do(func() {
		userServiceInstance = &UserService{
			s3idDao: model.Session3idDaoInstance(),
			userDao: model.UserDaoInstance(),
		}
	})
	return userServiceInstance
}

/**
1）小程序内通过wx.user接口获得code
2）将code传入后台，后台对微信服务器发起一个https请求换取openid、session_key
3）后台生成一个自身的3rd_session（以此为key值保持openid和session_key），返回给前端。PS:微信方的openid和session_key并没有发回给前端小程序
4）小程序拿到3rd_session之后保持在本地
5）小程序请求登录区内接口，通过wx.checkSession检查登录态，如果失效重新走上述登录流程，否则待上3rd_session到后台进行登录验证
*/
/**
外部接口：用户登录

req:
- code：用户登录码

resp:
- tempId: 用户临时 ID
*/
func (s *UserService) WxLogin(ctx context.Context, code string) (tempId string, err error) {
	logs.CtxInfo(ctx, "[Service UserService WxLogin Req] req: %s", code) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service UserService WxLogin Resp] resp: %s, %#v", tempId, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	var openId string    // 定义 openid
	if config.IsProd() { // 生产环境调用微信接口获取 openid
		logs.CtxInfo(ctx, "[Service UserService WxLogin] HttpGet access url: %#v", config.ConfigInstance.WxLoginUrl)
		// 调用微信 HTTP 接口
		resp, err := util.HttpGet(config.ConfigInstance.WxLoginUrl,
			config.ConfigInstance.AppId,
			config.ConfigInstance.AppSecret,
			code)
		// 调用微信 HTTP 接口失败
		if err != nil {
			logs.CtxError(ctx, "[Service UserService WxLogin] get openid error: %#v", err)
			err = util.NewErrf("[Service UserService WxLogin] get openid error: %#v", err)
			return "", err
		}
		logs.CtxInfo(ctx, "[Service UserService WxLogin] HttpGet resp: %#v", resp)
		// 转换参数
		var respStruct model.OpenIdAndSessionKey
		err = json.Unmarshal([]byte(resp), &respStruct)
		if err != nil {
			logs.CtxError(ctx, "[Service UserService WxLogin] HttpGet resp Unmarshal error: %#v, resp: %#v", err, resp)
			err = util.NewErrf("[Service UserService WxLogin] HttpGet resp Unmarshal error: %#v, resp: %#v", err, resp)
			return "", err
		}
		if util.IsNil(respStruct.OpenId) {
			logs.CtxError(ctx, "[Service UserService WxLogin] HttpGet resp error, resp: %#v", resp)
			err = util.NewErrf("[Service UserService WxLogin] HttpGet resp error, resp: %#v", resp)
			return "", err
		}
		openId = respStruct.OpenId // 获取 openId
	} else { // 开发环境 mock openid
		logs.CtxInfo(ctx, "[Service UserService WxLogin] access url: mock")
		openId = util.GenUUID(4)
	}
	// 获取 openid 成功
	logs.CtxInfo(ctx, "[Service UserService WxLogin] get openid from url success: %#v", openId)
	// 写入 redis，获取 tempId
	tempId, err = s.s3idDao.SetOpenIdGenID(ctx, openId)
	if err != nil {
		logs.CtxError(ctx, "[Service UserService WxLogin] SetOpenIdGenID error: %#v", err)
		err = util.NewErrf("[Service UserService WxLogin] SetOpenIdGenID error: %#v", err)
		return "", err
	}
	/************************************ 核心逻辑 ****************************************/

	return tempId, nil
}

/**
外部接口：更新用户信息

req:
- openId: 用户 openId
- userInfo: 用户信息

resp:
*/
func (s *UserService) UpdateUserInfo(ctx context.Context, openId string, userInfo *model.User) (err error) {
	logs.CtxInfo(ctx, "[Service UserService UpdateUserInfo Req] req: %s, %#v", openId, userInfo) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service UserService UpdateUserInfo Resp] resp: %#v", err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	// 先判断用户是否已存在，如果用户已存在则更新用户信息，否则创建新用户
	user, err := s.userDao.QueryUser(ctx, openId)
	if err != nil {
		logs.CtxError(ctx, "[Service UserService UpdateUserInfo] QueryUser error: %#v", err)
		err = util.NewErrf("[Service UserService UpdateUserInfo] QueryUser error: %#v", err)
		return err
	}
	if user == nil { // 用户不存在，创建用户
		err := s.userDao.CreateUser(ctx, openId, userInfo)
		if err != nil {
			logs.CtxError(ctx, "[Service UserService UpdateUserInfo] CreateUser error: %#v", err)
			err = util.NewErrf("[Service UserService UpdateUserInfo] CreateUser error: %#v", err)
			return err
		}
	} else { // 用户存在，更新用户
		err := s.userDao.UpdateUser(ctx, openId, userInfo)
		if err != nil {
			logs.CtxError(ctx, "[Service UserService UpdateUserInfo] UpdateUser error: %#v", err)
			err = util.NewErrf("[Service UserService UpdateUserInfo] UpdateUser error: %#v", err)
			return err
		}
	}
	/************************************ 核心逻辑 ****************************************/

	return nil
}

/**
外部接口：获取用户 openid

req:
- tempId: 用户临时 ID

resp:
- openid: 用户 openid
*/
func (s *UserService) GetUserOpenId(ctx context.Context, tempId string) (openid string, err error) {
	logs.CtxInfo(ctx, "[Service UserService getUserOpenid Req] req: %s", tempId) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Service UserService getUserOpenid Resp] resp: %s", openid) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	openid, err = s.s3idDao.GetOpenId(ctx, tempId)
	if err == cconst.RedisNil {
		logs.CtxError(ctx, "[Service UserService getUserOpenid] GetOpenId: %#v", err) // TODO：降级为 Warn
		return "", cconst.RedisNil
	}
	if err != nil {
		logs.CtxError(ctx, "[Service UserService getUserOpenid] GetOpenId error: %#v", err)
		return "", err
	}
	/************************************ 核心逻辑 ****************************************/

	return openid, nil
}
