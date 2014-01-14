package controllers

import (
	"github.com/robfig/revel"
)

type Home struct {
	*revel.Controller
}

func (c *Home) Index() revel.Result {

	// First to localhost/home the error text can't show
	// but  c.Flash.Out["error"] has content

	if c.Session["email"] != "" {
		c.Flash.Success("Already logged in!")
	} else {
		c.Flash.Error("Please login or register!")
	}

	return c.Render()
}
