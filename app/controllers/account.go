package controllers

import (
	"github.com/robfig/revel"
	"regexp"
	"revel_userm/app/models"
)

type Account struct {
	*revel.Controller
}

var USERNAME_REX, PWD_REX, NICKNAME_REX *regexp.Regexp

func init() {

	USERNAME_REX = regexp.MustCompile(`^[0-9a-zA-Z]+@(([0-9a-zA-Z]+)[.])+[a-z]{2,4}$`)
	PWD_REX = regexp.MustCompile(`^[\x01-\xfe]{8,20}$`)

	NICKNAME_REX = regexp.MustCompile(`^[a-z0-9_]{6,16}$`)
}

///////////////////////////////////////////////////////////////////////////
// 注册
func (c *Account) Register() revel.Result {
	return c.Render()
}

func (c *Account) RegisterSuccessful() revel.Result {
	return c.Render()
}

func (c *Account) PostRegister(user *models.MockUser) revel.Result {

	c.Validation.Required(user.Email).Message("邮箱地址未填写！")
	c.Validation.Check(user.Email, revel.Match{USERNAME_REX}).Message("邮箱地址格式错误！")

	c.Validation.Required(user.Nickname).Message("昵称未填写！")
	c.Validation.Check(user.Nickname, revel.Match{NICKNAME_REX}).Message("昵称格式错误！")

	c.Validation.Required(user.Password).Message("密码未填写！")
	c.Validation.Check(user.Password, revel.Match{PWD_REX}).Message("密码格式错误！")

	c.Validation.Required(user.ConfirmPassword == user.Password).Message("两次输入密码不一致！")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect((*Account).Register) // 用 c.Register 会提示错误
	}

	err := user.Register()

	if err != nil {
		c.Flash.Error(err.Error())

		c.FlashParams()
		return c.Redirect((*Account).Register)
	}

	// 在客户端 的存储是 签名明文
	// 直接这样写会有问题
	// 应该咋写捏？？
	c.Session["email"] = user.Email
	c.Session["password"] = user.Password

	c.Flash.Success("注册并登录成功！")

	return c.Redirect((*Home).Index)
}

///////////////////////////////////////////////////////////////////////////
// 登录
func (c *Account) Login() revel.Result {
	return c.Render()
}

func (c *Account) PostLogin(loginuser *models.LoginUser) revel.Result {

	c.Validation.Required(loginuser.Email).Message("邮箱地址未填写！")
	c.Validation.Check(loginuser.Email, revel.Match{USERNAME_REX}).Message("邮箱地址格式错误！")

	c.Validation.Required(loginuser.Password).Message("密码未填写！")
	c.Validation.Check(loginuser.Password, revel.Match{PWD_REX}).Message("密码格式错误！")

	ok, _ := loginuser.IsRegistered()
	if !ok {
		c.Flash.Error("Email或密码错误！")

		c.FlashParams()
		return c.Redirect((*Account).Login)
	}

	c.Session["email"] = loginuser.Email
	c.Session["password"] = loginuser.Password

	c.Flash.Success("登录成功！")

	return c.Redirect((*Home).Index)
}

///////////////////////////////////////////////////////////////////////////
// 退出
func (c *Account) Logout() revel.Result {

	for k := range c.Session {
		delete(c.Session, k)
	}

	return c.Redirect((*Account).Login)
}
