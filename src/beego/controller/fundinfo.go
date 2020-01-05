package controller

import (
    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"../../analyze"
	"../../util"
)


type GetFundInfo struct {
    beego.Controller
}

func (p *GetFundInfo) Get() {
	code := p.GetString("code")

	logs.Debug("GetFundInfo process")
	logs.Debug("get input data, code=", code)

	if code == "" {
		return
	}

	basicInfo := analyze.GetFundBasicInfoByCode(code)
	logs.Debug("get fund basic info:", basicInfo)

	p.Data["code"] = code
	p.Data["name"] = basicInfo[0]
	p.Data["type"] = basicInfo[1]

	growth, startDate := analyze.GetMyGrowth(code)  //交易数据
	price := analyze.GetFundPriceByCode(code)  //价格趋势
	p.Data["price"] = price
	p.Data["growth"] = analyze.GetGrowthRateByCode(code)  //增长趋势
	p.Data["growth2"] = analyze.GetGrowthRateFromBeginByCode(code, startDate) 
	
	p.Data["transGrowth"] = growth

	// 每月的收益情况
	p.Data["monthIncome"] = analyze.GetFundIncomeByMonthInRecentYear(code)

	// 累计收益
	p.Data["accumulatedIncome"],p.Data["accumulatedIncomePercent"],p.Data["handlingUnits"], p.Data["cost"] = analyze.GetFundAccumulatedIncome(code)
	p.Data["handlingIncome"], p.Data["handlingIncomePercent"] = analyze.GetFundHandlingIncome(code)

	// 最新净值，累计价值
	index := len(price)
	if index > 0 {
		index = index - 1
	}

	// 最新增长率
	if index > 0 {
		rate := (price[index].Ljjz - price[index - 1].Ljjz) * 100 / price[index - 1].Jjjz
		p.Data["curRate"] = util.GetFloatFormat(rate, 2)

	} else {
		p.Data["curRate"] = 0.0
	}
	p.Data["jjjz"] = price[index].Jjjz
	p.Data["Ljjz"] = price[index].Ljjz
	p.Data["date"] = price[index].Date
	
	p.Data["costList"] = analyze.GetTransCost(code)

	p.TplName = "fundinfo.html"
}