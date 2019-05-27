package analyze

import (
	"../util"
	//"fmt"
)

type FundTransGrowth struct {
	Date string
	Code string
	Units float64 //交易份额
	Amount float64  //交易金额
	CurUnits float64 //当日份额
	CurAmount float64 //当日总成本
	CurValue float64 //当日总价值
	CurIncome float64 //当日收入
	CurIncomeRate float64 //当日收入增长率
}

func GetMyGrowth(code string) ([]FundTransGrowth, string) {
	var fundTransGrowth []FundTransGrowth
	var firstTransDate string

	trans := GetTransData(code)
	price := GetFundPriceByCode(code)

	lastDayAmount := 0.0 //上一日的总成本
	lastDayUnits := 0.0 //上一日的总份额

	j := 0
	start := 0

	for i := 0; i < len(price); i++ {

		//记录一个交易数据
		var growth FundTransGrowth

		growth.Date = price[i].Date
		growth.Code = code

		//找到一次交易，变更平均价格和交易数据
		if price[i].Date != trans[j].Date {
			growth.Units = 	0
			growth.Amount = 0
		} else {

			if j == 0 {
				firstTransDate = trans[j].Date
			}

			growth.Units = 	trans[j].Units
			growth.Amount = trans[j].Amount

			if j < len(trans) - 1 {
				j++
			}

			start = 1
		}

		if start == 0 {
			continue
		}

		growth.CurAmount = lastDayAmount + growth.Amount
		growth.CurUnits = lastDayUnits + growth.Units
		growth.CurValue = growth.CurUnits * price[i].Jjjz
		growth.CurIncome = lastDayUnits * price[i].Jjjz - lastDayAmount //要用上一日的成本和份额算收入
		growth.CurIncomeRate = growth.CurIncome / lastDayAmount * 100 //要用上一日的成本和今天的收入算收益率
		growth.CurIncomeRate = util.GetFloatFormat(growth.CurIncomeRate, 3)

		lastDayAmount = growth.CurAmount
		lastDayUnits = growth.CurUnits

		//fmt.Println(growth)

		fundTransGrowth = append(fundTransGrowth, growth)
	}

	return fundTransGrowth, firstTransDate
}
