package controllers

import (
	"context"
	"encoding/json"
	"log"
	"oyster-iot/init/cache"
	"oyster-iot/services"
	"time"

	bcrypt "oyster-iot/utils"
	jwt "oyster-iot/utils"

	djwt "github.com/dgrijalva/jwt-go"

	"github.com/beego/beego/v2/core/validation"
)

type AuthController struct {
	BaseController
}

type TokenData struct {
	Msg       string `json:"msg"`
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
func (this *AuthController) Login() {
	log.Println(this.Ctx.Input.RequestBody)
	loginInfo := LoginInfo{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &loginInfo)
	if err != nil {
		log.Println("Json Unmarshal Failed!", err.Error())
		this.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&loginInfo)
	if err != nil {
		// handler error
		this.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			log.Println(err.Key, err.Message)
		}
		this.Response(400, "输入参数错误")
		return
	}
	//检查用户名或者密码

	var UserService services.UserService
	user, err := UserService.GetUserByUsername(loginInfo.Username)
	if err != nil || !bcrypt.ComparePasswords(user.Password, []byte(loginInfo.Password)) {
		this.Response(400, "用户名或密码错误")
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
		this.Response(400, "JWT失败")
		return
	}
	d := TokenData{
		Msg:       "登录成功",
		Token:     token,
		TokenType: "oyster",
		Expires:   3600,
	}

	//将Token写入缓存
	cache.Bm.Put(context.TODO(), token, 1, 3000*time.Second)
	// 登录成功
	this.Response(200, d)
}

// 退出登录
func (this *AuthController) Logout() {
	authorization := this.Ctx.Request.Header["Authorization"][0]
	userToken := authorization[len(jwt.JWTType)+1:]
	_, err := jwt.ParseCliamsToken(userToken)
	if err != nil {
		this.Response(400, "token异常")
		return
	}
	s, err := cache.Bm.IsExist(context.TODO(), userToken)
	if s {
		cache.Bm.Delete(context.TODO(), userToken)
	}
	this.Response(200, "退出成功")
	return
}

// 刷新token
func (this *AuthController) Refresh() {
	authorization := this.Ctx.Request.Header["Authorization"][0]
	userToken := authorization[len(jwt.JWTType)+1:]
	user, err := jwt.ParseCliamsToken(userToken)
	if err != nil {
		this.Response(400, "token异常")
		return
	}
	s, _ := cache.Bm.IsExist(context.TODO(), userToken)
	if s {
		cache.Bm.Delete(context.TODO(), userToken)
	}
	var UserService services.UserService
	u, err := UserService.GetUserById(user.Id)
	if err != nil {
		this.Response(400, "该账户不存在")
		return
	}
	// 生成jwt
	tokenCliams := jwt.UserClaims{
		Id:         u.Id,
		Username:   u.Username,
		CreateTime: time.Now(),
		StandardClaims: djwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 3600,
		},
	}
	token, err := jwt.MakeCliamsToken(&tokenCliams)
	if err != nil {
		// JWT失败
		this.Response(400, "JWT失败")
		return
	}

	d := TokenData{
		Msg:       "Token刷新成功",
		Token:     token,
		TokenType: jwt.JWTType,
		Expires:   3600,
	}

	//将Token写入缓存
	cache.Bm.Put(context.TODO(), token, 1, 3000*time.Second)
	// 登录成功
	this.Response(200, d)
}
