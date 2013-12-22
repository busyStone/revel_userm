package controllers

import (
	"github.com/robfig/revel"
	"revel_userm/app/models"
)

type Home struct {
	*revel.Controller
}

func (c *Home) Index() revel.Result {

	var loginuser models.LoginUser

	loginuser.Email = c.Session["email"]
	loginuser.Password = c.Session["password"]

	ok, _ := loginuser.IsRegistered()
	if ok {
		c.Flash.Success("已经登录！")
	} else {
		c.Flash.Error("请先登录或注册！")
	}

	c.FlashParams()

	return c.Render()
}
