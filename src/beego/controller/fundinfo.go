package controller

import (
    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"../../analyze"
	//"../../util"
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
	p.Data["price"] = analyze.GetFundPriceByCode(code)  //价格趋势
	p.Data["growth"] = analyze.GetGrowthRateByCode(code)  //增长趋势
	p.Data["growth2"] = analyze.GetGrowthRateFromBeginByCode(code, startDate) 
	
	p.Data["transGrowth"] = growth

	// 每月的收益情况
	p.Data["monthIncome"] = analyze.GetFundIncomeByMonthInRecentYear(code)

	// 累计收益
	p.Data["accumulatedIncome"],p.Data["accumulatedIncomePercent"],p.Data["cost"] = analyze.GetFundAccumulatedIncome(code)
	p.Data["handlingIncome"], p.Data["handlingIncomePercent"] = analyze.GetFundHandlingIncome(code)

	p.TplName = "fundinfo.html"
}