package controller

import (
	//"fmt"
    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"../../analyze"
	"../../util"
)

type IndexController struct {
    beego.Controller
}

type MyFundInfo struct {
	Code string
	Name string
	FundType string
	AccumulatedIncome float64
	AccumulatedIncomePercent float64
	Cost float64
}

type MyFundIncome struct {
	Date int64
	Income float64
}

func (p *IndexController) Index() {
	logs.Debug("enter index controller.....")

	var myFundIncom [12]MyFundIncome

	for i:= 0; i < len(myFundIncom); i++ {
		myFundIncom[i].Income = 0.0
	}

	p.TplName = "index.html"
	var myFundsInfo []MyFundInfo
	accumulatedIncome := 0.0
	accumulatedIncomePercent := 0.0
	cost := 0.0
	for code := range analyze.MyFundsList {
		var fundInfo MyFundInfo
		fundInfo.Code = code
		fundInfo.Name = analyze.MyFundsList[code][0]
		fundInfo.FundType = analyze.MyFundsList[code][1]
		fundInfo.AccumulatedIncome, fundInfo.AccumulatedIncomePercent, fundInfo.Cost = analyze.GetFundAccumulatedIncome(code)
		myFundsInfo = append(myFundsInfo, fundInfo)

		accumulatedIncome += fundInfo.AccumulatedIncome
		cost += fundInfo.Cost

		monthIncome := analyze.GetFundIncomeByMonthInRecentYear(code)
		for i:= 0; i < len(monthIncome); i++ {
			myFundIncom[i].Date = monthIncome[i].Date
			myFundIncom[i].Income += monthIncome[i].Income
		}
	}

	if (cost > 0) {
		accumulatedIncomePercent = accumulatedIncome * 100 / cost
	}

	p.Data["funds"] = myFundsInfo
	p.Data["num"] = len(myFundsInfo)
	p.Data["accumulatedIncome"] = util.GetFloatFormat(accumulatedIncome, 2)
	p.Data["cost"] = cost
	p.Data["accumulatedIncomePercent"] = util.GetFloatFormat(accumulatedIncomePercent, 2)
	p.Data["monthIncome"] = myFundIncom
}