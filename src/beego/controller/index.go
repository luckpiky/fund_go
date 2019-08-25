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
	HandlingIncome float64
	HandlingIncomePercent float64
	Cost float64
}

type MyFundIncome struct {
	Date int64
	Income float64
	Cost float64
}

type FundType struct {
	Name string
	Income float64
	Cost float64
	AccumulatedIncomePercent float64
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
		fundInfo.HandlingIncome, fundInfo.HandlingIncomePercent = analyze.GetFundHandlingIncome(code)
		myFundsInfo = append(myFundsInfo, fundInfo)

		accumulatedIncome += fundInfo.AccumulatedIncome
		cost += fundInfo.Cost

		monthIncome := analyze.GetFundIncomeByMonthInRecentYear(code)
		for i:= 0; i < len(monthIncome); i++ {
			myFundIncom[i].Date = monthIncome[i].Date
			myFundIncom[i].Income += monthIncome[i].Income
			myFundIncom[i].Cost += monthIncome[i].Cost
		}
	}

	if (cost > 0) {
		accumulatedIncomePercent = accumulatedIncome * 100 / cost
	}

	sort.Stable(myFundsInfo)

	var myFundIncomePercent [12]MyFundIncome
	for i := 0; i < len(myFundIncom); i++ {
		myFundIncomePercent[i].Income = util.GetFloatFormat(myFundIncom[i].Income * 100 / myFundIncom[i].Cost, 2)
		myFundIncomePercent[i].Date = myFundIncom[i].Date

		myFundIncom[i].Income = util.GetFloatFormat(myFundIncom[i].Income, 2)

		//logs.Debug(myFundIncom[i].Date, myFundIncom[i].Income, myFundIncom[i].Cost, myFundIncomePercent[i].Income)
	}

	// 根据类型分类，画饼图
	var fundTypes []FundType
	for i := 0; i < len(myFundsInfo); i++ { // 先分类
		find := false
		for j := 0; j < len(fundTypes); j++ {
			if fundTypes[j].Name == myFundsInfo[i].FundType {
				fundTypes[j].Income += myFundsInfo[i].AccumulatedIncome
				fundTypes[j].Income = util.GetFloatFormat(fundTypes[j].Income, 2)
				fundTypes[j].Cost += myFundsInfo[i].Cost
				fundTypes[j].Cost = util.GetFloatFormat(fundTypes[j].Cost, 2)
				fundTypes[j].AccumulatedIncomePercent = fundTypes[j].Income * 100 / fundTypes[j].Cost
				fundTypes[j].AccumulatedIncomePercent = util.GetFloatFormat(fundTypes[j].AccumulatedIncomePercent, 2)
				find = true
			}
		}
		if find == false {
			var fundType FundType
			fundType.Name = myFundsInfo[i].FundType
			fundType.Income += myFundsInfo[i].AccumulatedIncome
			fundType.Income = util.GetFloatFormat(fundType.Income, 2)
			fundType.Cost += myFundsInfo[i].Cost
			fundType.Cost = util.GetFloatFormat(fundType.Cost, 2)
			fundType.AccumulatedIncomePercent = util.GetFloatFormat(fundType.Income * 100 / fundType.Cost, 2)
			fundTypes = append(fundTypes, fundType)
		}
	}

	var fundTypeItems []FundType
	for i := 0; i < len(fundTypes); i++ { // 根据分类的顺序增加基金
		for j := 0; j < len(myFundsInfo); j++ {
			if fundTypes[i].Name == myFundsInfo[j].FundType {
				var fundType FundType
				fundType.Name = myFundsInfo[j].Name
				fundType.Income += myFundsInfo[j].AccumulatedIncome
				fundType.Income = util.GetFloatFormat(fundType.Income, 2)
				fundType.Cost += myFundsInfo[j].Cost
				fundType.Cost = util.GetFloatFormat(fundType.Cost, 2)
				fundTypeItems = append(fundTypeItems, fundType)
			}
		}
	}

	p.Data["funds"] = myFundsInfo
	p.Data["num"] = len(myFundsInfo)
	p.Data["accumulatedIncome"] = util.GetFloatFormat(accumulatedIncome, 2)
	p.Data["cost"] = util.GetFloatFormat(cost, 2)
	p.Data["accumulatedIncomePercent"] = util.GetFloatFormat(accumulatedIncomePercent, 2)
	p.Data["monthIncome"] = myFundIncom
	p.Data["monthIncomePercent"] = myFundIncomePercent
	p.Data["fundTypes"] = fundTypes
	p.Data["fundTypeItems"] = fundTypeItems
}