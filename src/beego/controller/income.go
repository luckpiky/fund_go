package controller

import (
	"container/list"
	//"log"
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

	//income := analyze.GetInComeData("000311")
	//log.Println(income)

	//analyze.GetFundIncomeByMonthInRecentYear("519062")

	allData := list.New()
	analyze.GetAllFundIncomeData(analyze.GetDataPath(), allData)

	data := []FundIncomeShowData{}

	for item := allData.Front(); item != nil; item = item.Next() {
		itemData := item.Value.(analyze.FundAnalyzeData)

		logs.Debug("fund code:", itemData.FundCode)

		var d1 FundIncomeShowData
		d1.Code = itemData.FundCode
		data = append(data, d1)
	}

	p.Data["fund"] = data

    p.TplName = "index.html"
}