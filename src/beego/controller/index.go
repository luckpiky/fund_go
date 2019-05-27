package controller

import (
	"fmt"
    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"../../analyze"
)

type IndexController struct {
    beego.Controller
}

func (p *IndexController) Index() {
	logs.Debug("enter index controller.....")

	p.TplName = "index.html"
	
	for item := range analyze.MyFundsList {
		fmt.Println("::", item)
	}

	p.Data["funds"] = analyze.MyFundsList
}