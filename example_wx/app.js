//app.js

const DefaultPostHeader = {
    "Content-Type": "application/x-www-form-urlencoded"
};
const FilePostHeader = {
    "Content-Type": "multipart/form-data"
};
const JsonPostHeader = {
    "Content-Type": "application/json"
};
const TextPostHeader = {
    "Content-Type": "text/xml"
};

const server = "http://127.0.0.1:9205"
// const server = "http://something:9205"

const WxLoginUrl = server + "/user/wx_login";
const UpdateUserInfoUrl = server + "/user/update_userInfo";

var md5 = require('./md5.js');
var salt = "something"

App({
    globalData: {
        userInfo: null,
        id: null,
    },
    onLaunch: function () {
        // TODO test sign
        var sign = md5.hexMD5("哈哈哈哈哈");
        console.log(sign);// 加密后的"哈哈哈哈哈"： d98293b3ec6c0050f759e9492b2ba24d

        let that = this;
        wx.showLoading({
            title: '登录中',
            mask: true
        });
        // 获取id
        wx.checkSession({
            success() {
                // session_key 未过期，并且在本生命周期一直有效
                // 第三方id丢失，重新认证
                if (!that.globalData.id) {
                    console.log("no session")
                    that.wxLogin()
                } else {
                    wx.hideLoading();
                    console.log("has session, login success, id: " + that.globalData.id)
                }
            },
            fail() {
                // session_key 已经失效，需要重新执行登录流程
                that.wxLogin();
            }
        });
        // 如果已经授权，更新个人信息
        wx.getSetting({
            success: res => {
                if (res.authSetting['scope.userInfo']) {
                    // 已经授权，可以直接调用 getUserInfo 获取头像昵称，不会弹框
                    wx.getUserInfo({
                        success: res => {
                            console.log("getUserInfo")
                            // 可以将 res 发送给后台解码出 unionId
                            that.globalData.userInfo = res.userInfo;
                            // 由于 getUserInfo 是网络请求，可能会在 Page.onLoad 之后才返回
                            // 所以此处加入 callback 以防止这种情况
                            if (that.userInfoReadyCallback) {
                                that.userInfoReadyCallback(res.userInfo)
                            }
                            // 如果已经获取到 id，直接更新 userInfo
                            // 如果没有获取到 id，定义 OpenIdReady 的回调
                            if (that.globalData.id) {
                                console.log("direct: that.updateUserInfo")
                                that.updateUserInfo();
                            } else {
                                this.idReadyCallback = function () {
                                    console.log("callback: that.updateUserInfo")
                                    that.updateUserInfo();
                                };
                            }
                        }
                    })
                }
            }
        })
    },
    wxLogin: function () {
        console.log("wx login");
        let that = this;
        wx.login({
            success: res => {
                console.log(res);
                if (res.code) {
                    let data = {code: res.code};
                    let successFunc = function (resp) {
                        wx.hideLoading();
                        that.globalData.id = resp.Data;
                        console.log("login success, id: " + that.globalData.id);
                        if (that.idReadyCallback) {
                            that.idReadyCallback()
                        }
                    };
                    let requestFailFunc = function () {
                        wx.hideLoading();
                        wx.showToast({
                            title: '服务器维护中',
                            icon: 'none'
                        })
                    };
                    let responseFailFunc = function () {
                        wx.hideLoading();
                        wx.showToast({
                            title: '登录失败',
                            icon: 'none'
                        })
                    };
                    that.WxPostRequest(WxLoginUrl, DefaultPostHeader, data, successFunc, requestFailFunc, responseFailFunc);
                }
            }
        });
    },
    updateUserInfo: function () {
        console.log("update userInfo");
        let data = {
            tempId: this.globalData.id,
            userInfo: JSON.stringify(this.globalData.userInfo)
        };
        console.log(data);
        this.WxPostRequest(UpdateUserInfoUrl, DefaultPostHeader, data)
    },
    // url, header, data 必传参数
    // successFunc, requestFailFunc, responseFailFunc 选传参数，成功函数，请求失败函数（未成功请求），回应失败函数（请求的返回值为失败）
    // 返回值 data 格式：
    // type HTTPResponse struct {
    //     Success bool
    //     Message string
    //     Data    interface{}
    // }
    WxPostRequest: function (url, header, data, successFunc, requestFailFunc, responseFailFunc) {
        let reqDataJsonString = JSON.stringify(data)

        wx.request({
            url: url,
            header: header,
            // data: data,
            data: {
                "data": reqDataJsonString,
                "sign": md5.hexMD5(salt + reqDataJsonString + salt),
            },
            method: 'POST',
            success: function (res) {
                console.log(url + " response: ");
                console.log(res);
                if (res.data.Success) {
                    console.log(url + ": success");
                    if (successFunc) {
                        successFunc(res.data)
                    }
                } else {
                    console.log(url + ": response fail, msg: " + res.data.Message);
                    if (responseFailFunc) {
                        responseFailFunc()
                    }
                }
            },
            fail: function () {
                console.log(url + ": request fail");
                if (requestFailFunc) {
                    requestFailFunc()
                }
            }
        })
    }
})
;