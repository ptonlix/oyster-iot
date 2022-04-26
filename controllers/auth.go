package controllers

import (
	"context"
	"encoding/json"
	"oyster-iot/init/cache"
	"oyster-iot/services"
	"time"

	bcrypt "oyster-iot/utils"
	jwt "oyster-iot/utils"

	djwt "github.com/dgrijalva/jwt-go"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type AuthController struct {
	BaseController
}

type TokenData struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	Expires   int    `json:"expires"`
}

// TODO:正则校验，防止SQL注入
type LoginInfo struct {
	Username string `json:"username" valid:"Required;MaxSize(100);Match(/[a-zA-Z0-9_]{3,16}/)"`
	Password string `json:"passwd" valid:"Required;MinSize(6);MaxSize(100)"`
}

// 如果你的 struct 实现了接口 validation.ValidFormer
// 当 StructTag 中的测试都成功时，将会执行 Valid 函数进行自定义验证
func (u *LoginInfo) Valid(v *validation.Validation) {
	//if strings.Index(u.Name, "admin") != -1 {
	// 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
	//	v.SetError("Name", "名称里不能含有 admin")
	//}
}

// 登录
func (u *AuthController) Login() {
	loginInfo := LoginInfo{}
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &loginInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		u.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&loginInfo)
	if err != nil {
		// handler error
		u.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		u.Response(400, "输入参数错误")
		return
	}
	//检查用户名或者密码

	var UserService services.UserService
	user, err := UserService.GetUserByUsername(loginInfo.Username)
	if err != nil || !bcrypt.ComparePasswords(user.Password, []byte(loginInfo.Password)) {
		u.Response(400, "用户名或密码错误")
		return
	}

	// 生成JWT
	tokenCliams := jwt.UserClaims{
		Id:         user.Id,
		Username:   user.Username,
		CreateTime: time.Now(),
		StandardClaims: djwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 3600,
		},
	}
	token, err := jwt.MakeCliamsToken(&tokenCliams)
	if err != nil {
		// jwt失败
		u.Response(500, "JWT失败")
		return
	}
	d := TokenData{
		Token:     token,
		TokenType: "oyster",
		Expires:   3600,
	}

	//将Token写入缓存
	cache.Bm.Put(context.TODO(), token, 1, 3000*time.Second)
	// 登录成功
	u.Response(200, "登录成功", d)
}

// 退出登录
func (u *AuthController) Logout() {
	authorization := u.Ctx.Request.Header["Authorization"][0]
	userToken := authorization[len(jwt.JWTType)+1:]
	_, err := jwt.ParseCliamsToken(userToken)
	if err != nil {
		u.Response(400, "token异常", nil)
		return
	}
	s, err := cache.Bm.IsExist(context.TODO(), userToken)
	if s {
		cache.Bm.Delete(context.TODO(), userToken)
	}
	u.Response(200, "退出成功")
	return
}

// 刷新token
func (u *AuthController) Refresh() {
	authorization := u.Ctx.Request.Header["Authorization"][0]
	userToken := authorization[len(jwt.JWTType)+1:]
	userInfo, err := jwt.ParseCliamsToken(userToken)
	if err != nil {
		u.Response(400, "token异常")
		return
	}
	s, _ := cache.Bm.IsExist(context.TODO(), userToken)
	if s {
		cache.Bm.Delete(context.TODO(), userToken)
	}
	var UserService services.UserService
	user, err := UserService.GetUserById(userInfo.Id)
	if err != nil {
		u.Response(500, "该账户不存在")
		return
	}
	// 生成jwt
	tokenCliams := jwt.UserClaims{
		Id:         user.Id,
		Username:   user.Username,
		CreateTime: time.Now(),
		StandardClaims: djwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 3600,
		},
	}
	token, err := jwt.MakeCliamsToken(&tokenCliams)
	if err != nil {
		// JWT失败
		u.Response(500, "JWT失败")
		return
	}

	d := TokenData{
		Token:     token,
		TokenType: jwt.JWTType,
		Expires:   3600,
	}

	//将Token写入缓存
	cache.Bm.Put(context.TODO(), token, 1, 3000*time.Second)
	// 登录成功
	u.Response(200, "Token刷新成功", d)
}
