package controller

import (
	//"fmt"
	//"container/list"
    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"../../analyze"
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

	p.Data["name"] = basicInfo[0]
	p.Data["type"] = basicInfo[1]

	growth, startDate := analyze.GetMyGrowth(code)  //交易数据
	p.Data["price"] = analyze.GetFundPriceByCode(code)  //价格趋势
	p.Data["growth"] = analyze.GetGrowthRateByCode(code)  //增长趋势
	p.Data["growth2"] = analyze.GetGrowthRateFromBeginByCode(code, startDate) 
	
	p.Data["transGrowth"] = growth

	p.TplName = "fundinfo.html"
}