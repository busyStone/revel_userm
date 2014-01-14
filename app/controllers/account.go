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

	USERNAME_REX = regexp.MustCompile(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`)
	PWD_REX = regexp.MustCompile(`^[\x01-\xfe]{8,20}$`)

	NICKNAME_REX = regexp.MustCompile(`^[a-z0-9_]{6,16}$`)
}

///////////////////////////////////////////////////////////////////////////
// Register
func (c *Account) Register() revel.Result {
	c.Flash.Error("")
	return c.Render()
}

func (c *Account) PostRegister(user *models.MockUser) revel.Result {

	c.Validation.Required(user.Email).Message("Email not filled!")
	c.Validation.Check(user.Email, revel.Match{USERNAME_REX}).Message("Email format error!")

	c.Validation.Required(user.Nickname).Message("Nickname not filled!")
	c.Validation.Check(user.Nickname, revel.Match{NICKNAME_REX}).Message("Nickname format error!")

	c.Validation.Required(user.Password).Message("Password is not filled!")
	c.Validation.Check(user.Password, revel.Match{PWD_REX}).Message("Password format error!")

	c.Validation.Required(user.ConfirmPassword == user.Password).Message("Please confirm the password!")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect((*Account).Register)
	}

	err := user.Register()

	if err != nil {
		c.Flash.Error(err.Error())

		c.FlashParams()
		return c.Redirect((*Account).Register)
	}

	c.Session["email"] = user.Email

	c.Flash.Success("Register & Login successfully!")

	return c.Redirect((*Home).Index)
}

///////////////////////////////////////////////////////////////////////////
// Login
func (c *Account) Login() revel.Result {
	return c.Render()
}

func (c *Account) PostLogin(loginuser *models.LoginUser) revel.Result {

	c.Validation.Required(loginuser.Email).Message("Email not filled!")
	c.Validation.Check(loginuser.Email, revel.Match{USERNAME_REX}).Message("Email format error!")

	c.Validation.Required(loginuser.Password).Message("Password is not filled!")
	c.Validation.Check(loginuser.Password, revel.Match{PWD_REX}).Message("Password format error!")

	ok, _ := loginuser.IsRegistered()
	if !ok {
		c.Flash.Error("Email or Password is incorrect!")

		c.FlashParams()
		return c.Redirect((*Account).Login)
	}

	c.Session["email"] = loginuser.Email

	c.Flash.Success("Login successfully!")

	return c.Redirect((*Home).Index)
}

///////////////////////////////////////////////////////////////////////////
// Logout
func (c *Account) Logout() revel.Result {

	for k := range c.Session {
		delete(c.Session, k)
	}

	return c.Redirect((*Account).Login)
}
