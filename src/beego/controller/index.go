package controller

import (
	"sort"
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

type MyFundsInfo []MyFundInfo

func (s MyFundsInfo) Len() int { return len(s) }
func (s MyFundsInfo) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s MyFundsInfo) Less(i, j int) bool { return s[i].AccumulatedIncome > s[j].AccumulatedIncome }

func (p *IndexController) Index() {
	logs.Debug("enter index controller.....")

	var myFundIncom [12]MyFundIncome

	for i:= 0; i < len(myFundIncom); i++ {
		myFundIncom[i].Income = 0.0
	}

	p.TplName = "index.html"
	var myFundsInfo MyFundsInfo
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

	sort.Stable(myFundsInfo)

	p.Data["funds"] = myFundsInfo
	p.Data["num"] = len(myFundsInfo)
	p.Data["accumulatedIncome"] = util.GetFloatFormat(accumulatedIncome, 2)
	p.Data["cost"] = util.GetFloatFormat(cost, 2)
	p.Data["accumulatedIncomePercent"] = util.GetFloatFormat(accumulatedIncomePercent, 2)
	p.Data["monthIncome"] = myFundIncom
}