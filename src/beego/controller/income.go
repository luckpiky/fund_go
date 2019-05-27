package controller

import (
	"container/list"
    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"../../analyze"
)

type GetFundIncome struct {
    beego.Controller
}

type FundIncomeShowData struct {
	Code string
	Income []analyze.FundTransData
}

func (p *GetFundIncome) Get() {
	logs.Debug("enter income controller.....")

	allData := list.New()

	analyze.GetAllFundIncomeData("../test/csv", allData)

	data := []FundIncomeShowData{}

	for item := allData.Front(); item != nil; item = item.Next() {
		itemData := item.Value.(analyze.FundAnalyzeData)

		logs.Debug(itemData.FundCode)

		var d1 FundIncomeShowData
		d1.Code = itemData.FundCode
		data = append(data, d1)
	}

	p.Data["fund"] = data

    p.TplName = "index.html"
}