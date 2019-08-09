package controller

import (
    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"../../analyze"
	"../../catch"
)


type Update struct {
    beego.Controller
}

func (p *Update) Get() {

	for code := range analyze.MyFundsList {

		logs.Debug("update", code)

		threadFlag := 0
		catch.ReadOneFundData(code, 1, analyze.GetDataPath(), &threadFlag)
	}

	p.Data["result"] = "更新完成"
	p.TplName = "update.html"
}