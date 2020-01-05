package analyze

import (
//	"fmt"
	"io/ioutil"
	"strings"
	"encoding/csv"
	"log"
	"strconv"
	"container/list"
	"time"
	"../util"
)

type FundTransData struct {
	code string
	time string
    buy float64
	sell float64
	value float64 //当前价值
	average float64 //平均累计价格
	count float64 //当前份额
	income float64  //收益
}

type FundPriceData struct {
	time string
	jjjz float64
	ljjz float64
}

type FundAnalyzeData struct {
	//fundCode string
	fundPriceData *list.List
	fundTransData *list.List
	FundCode string
}

func ReadFundTransactionData(path string, analyzeData FundAnalyzeData) {
	var transData [][] string

	fileName := path + "/fund_transaction.csv"

	readContent,err := ioutil.ReadFile(fileName)
    if err == nil {
        csvReadFp := csv.NewReader(strings.NewReader(string(readContent)))
        transData,_ = csvReadFp.ReadAll()
        log.Println(fileName, "读取到基金交易数据条目数：",len(transData))
    } else {
		log.Println(err)
	}

	for i := 1; i < len(transData); i++ {
		var data FundTransData

		log.Println(transData[i][0], transData[i][1], transData[i][2], analyzeData.FundCode)

		if (transData[i][0] != analyzeData.FundCode) {
			continue
		}

		//log.Println(transData[i][0])

		data.code = transData[i][0]
		data.time = transData[i][1]
		data.buy,_  = strconv.ParseFloat(transData[i][2],64)
		data.sell,_ = strconv.ParseFloat(transData[i][3],64)

		analyzeData.fundTransData.PushBack(data)
		//analyzeData.fundTransData = append(analyzeData.fundTransData, data)
	}
}


func ReadFundPriceData(fundCode string, path string, analyzeData FundAnalyzeData) {
	var funData [][]string

	fileName := path + "/" + fundCode + ".csv"

    //读取老数据
    readContent,err := ioutil.ReadFile(fileName)
    if err == nil {
        csvReadFp := csv.NewReader(strings.NewReader(string(readContent)))
        funData,_ = csvReadFp.ReadAll()
        log.Println(fileName, "读取到基金数据条目数：",len(funData))
    } else {
		log.Println(err)
	}

	for i := 0; i < len(funData); i++ {
		var data FundPriceData

		data.time = funData[i][0]
		data.jjjz,_  = strconv.ParseFloat(funData[i][1],64)
		data.ljjz,_  = strconv.ParseFloat(funData[i][2],64)

		analyzeData.fundPriceData.PushBack(data)
	}
}

func GetFundPrice(fundPriceData *list.List, time string) (jjjz float64, ljjz float64) {

	for item := fundPriceData.Front(); item != nil; item = item.Next() {
		priceData := item.Value.(FundPriceData)

		if priceData.time == time {
			return priceData.jjjz, priceData.ljjz
		}
	}

	return 0,0
}

func AnalyzeGain(analyzeData FundAnalyzeData, path string) {
	ReadFundPriceData(analyzeData.FundCode, path, analyzeData)
	ReadFundTransactionData(path, analyzeData)

	PrevAverage := 0.0
	prevCount := 0.0
	preIncome := 0.0

	//遍历所有的买卖记录，对数据计算收益情况，数据必须按照时间顺序排列好，否则会计算出错
	for item := analyzeData.fundTransData.Front(); item != nil; item = item.Next() {
		transData := item.Value.(FundTransData)

		jjjz, ljjz := GetFundPrice(analyzeData.fundPriceData, transData.time)

		// 第一次购买时，总价值就是购买的数量×基金净值
		if (item == analyzeData.fundTransData.Front()) {
			transData.value = transData.buy * jjjz
			transData.average = ljjz
			transData.count = transData.buy
			transData.income = 0
		} else {

			//购买：
			//当前的总价值=（上一次的份额+本次购买的份额） × 基金净值
			//份额=上一次的份额+本次购买的份额
			//平均累计净值=（上一次的平均累计净值×上次的份额+这次的份额×累计净值） / 总份额
			//收入就是（累计净值-上次的平均累计净值） × 上次份额
			// 这里平均累计净值就是成本价格，但是为了与累计净值作比较，按照累计净值算
			if (transData.buy > 0) {
				transData.value = (prevCount + transData.buy) * jjjz
				transData.count = prevCount + transData.buy
				transData.average = (PrevAverage * prevCount + transData.buy * ljjz) / transData.count
				transData.income = (ljjz - PrevAverage) * prevCount

				PrevAverage = transData.average
				prevCount = transData.count
				preIncome = transData.income	
			}

			if (transData.sell > 0) {
				transData.value = (prevCount - transData.sell) * jjjz
				transData.count = prevCount - transData.sell
				transData.average = (PrevAverage * prevCount - transData.buy * jjjz + preIncome) / transData.count
				transData.income = (ljjz - PrevAverage) * prevCount
			}
		}

		PrevAverage = transData.average
		prevCount = transData.count
		preIncome = transData.income

		//fmt.Println(transData)
	}
}

func GetAllFundIncomeData(path string, analyzeData *list.List) {
	var fundData [][]string
	fundList := list.New()

	fileName := path + "/fund_transaction.csv"

    //读取老数据
    readContent,err := ioutil.ReadFile(fileName)
    if err == nil {
        csvReadFp := csv.NewReader(strings.NewReader(string(readContent)))
        fundData,_ = csvReadFp.ReadAll()
		
		for i := 1; i < len(fundData); i++ {

			find := 0
			for item := fundList.Front(); item != nil; item = item.Next() {
				fundCode := item.Value.(string)
				if fundCode == fundData[i][0] {
					find = 1
					break
				}
			}

			if find == 0 {
				fundList.PushBack(fundData[i][0])
			}
		}
	}

	for fund := fundList.Front(); fund != nil; fund = fund.Next() {
		var analyze FundAnalyzeData
		analyze.fundPriceData = list.New()
		analyze.fundTransData = list.New()
		analyze.FundCode = fund.Value.(string)
		AnalyzeGain(analyze, path)
		analyzeData.PushBack(analyze)
	}
}


func Init1() {

	allData := list.New()

	GetAllFundIncomeData("../test/csv", allData)
}

/* 获取指定时间的基金份额 */
func GetFundUnitsByTime() {

}

type FundIncomeData struct {
	Code string
	Date int64
	Income float64 // 收益
	Units float64  //基金份额
	Cost float64   // 成本
	AccumulatedIncome float64  // 累计收益
	HoldingIncome float64
}

/* 获取指定基金的历史收益 */
func GetInComeData(code string) (income []FundIncomeData) {
	transData := GetTransData2(code)
	priceData := GetFundPriceData(code)
	var incomeData []FundIncomeData

	units := 0.0
	cost := 0.0
	accumulatedIncome := 0.0
	holdingIncome := 0.0
	
	j := 0
	for i := 0; i < len(priceData); i++ {
		first := false
		pay := 0.0

		//sell := false

		// 存在交易
		if (priceData[i].Date == transData[j].Date) {
			if units == 0 {
				first = true
			}
			
			if (transData[j].Units != 0) {  // 分红，且份额不减少，暂时先这么算吧
				pay = transData[j].Amount - priceData[i].Jjjz * transData[j].Units
				pay = util.GetFloatFormat(pay, 2)
			}

			// 卖出时，成本剩余按照比例计算
			if  transData[j].Units >= 0 {
				cost += transData[j].Amount
			} else {
				cost = cost * (units + transData[j].Units) / units
			}

			//if (transData[j].Units < 0) {
				// 剩余收益 = 持有收益 × （总份额 - 卖出份额） / 总份额 - 费用
			//	holdingIncome = holdingIncome * (units + transData[j].Units) / units - pay
			//	log.Println("卖出，持有收益：", holdingIncome)
			//	sell = true
			//}

			dateStr := time.Unix(priceData[i].Date, 0).Format("2006-01-02 15:04:05") 
			log.Println(dateStr, "购买基金份额：", transData[j].Units, "交易金额：", transData[j].Amount, "交易费用：", pay, "累计成本：", cost)
		}

		// 持有收益 = 份额 * 基金净值 - 成本
		holdingIncome = priceData[i].Jjjz * units - cost
		log.Println(units, cost, holdingIncome)

		if (units != 0 || first) {
			var income FundIncomeData
			income.Cost = cost
			income.Units = units
			income.Code = code
			income.Date = priceData[i].Date
			if (first) {
				income.Income = 0 - pay
			} else {
				income.Income = units * (priceData[i].Ljjz - priceData[i-1].Ljjz) - pay
			}

			income.Income = util.GetFloatFormat(income.Income, 2)
			
			// 计算累计收益，将每个交易日的收益相加
			accumulatedIncome += income.Income
			income.AccumulatedIncome = accumulatedIncome

			//if (!sell) {
			//	holdingIncome += income.Income
			//}
			if (units + transData[j].Units == 0) {
				income.HoldingIncome = 0
			} else {
				income.HoldingIncome = holdingIncome
			}
			
			log.Println(holdingIncome)
	
			incomeData = append(incomeData, income)
	
			dateStr := time.Unix(priceData[i].Date, 0).Format("2006-01-02 15:04:05") 
			log.Println(dateStr, "基金收益：", income.Income, " 基金份额：", units, "持有收益:", income.HoldingIncome)
		}

		// 交易当天不能计算收益
		if (priceData[i].Date == transData[j].Date) {
			units += transData[j].Units
			units = util.GetFloatFormat(units, 2)
			if j < len(transData) - 1 {
				j++
			}
		}
	}

	log.Println("基金收益条目数：", len(incomeData))
	return incomeData
}


/* 获取指定时间范围的基金收益 */
func GetFundIncomeByTimeRange(code string, incomeData []FundIncomeData, time1 int64, time2 int64) (float64, float64) {
	begin := false
	income := 0.0
	cost := 0.0

	for i := 0; i < len(incomeData); i++ {
		if (time1 <= incomeData[i].Date && time2 >= incomeData[i].Date) {
			begin = true
		}

		if begin {
			income += incomeData[i].Income
			cost = incomeData[i].Cost
			//log.Println(incomeData[i].Date, incomeData[i].Income, incomeData[i].Cost)
		}

		if (time2 < incomeData[i].Date) {
			break
		}
	}

	return util.GetFloatFormat(income, 2), util.GetFloatFormat(cost, 2)
}

/* 获取最近一年中每月的基金收益 */
func GetFundIncomeByMonthInRecentYear(code string) (income []FundIncomeData)  {
	incomeData := GetInComeData(code)
	var incomeData2 []FundIncomeData

	year:=time.Now().Year()
	month:=time.Now().Month()

	for i := 0; i < 12; i++ {
		firstDay := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-1 00:00:00"
		//log.Println(firstDay)
		firstDayInt := util.TimeStr2Int64_2(firstDay)
		firstDayTime := time.Unix(firstDayInt, 0)
		lastDayTime := firstDayTime.AddDate(0, 1, -1)
		lastDayInt := lastDayTime.Unix()

		var income FundIncomeData
		income.Income, income.Cost = GetFundIncomeByTimeRange(code, incomeData, firstDayInt, lastDayInt)
		income.Date = firstDayInt
		income.Code = code
		incomeData2 = append(incomeData2, income)

		//log.Println(firstDay, "收益：", income)

		if (month == 1) {
			month = 12
			year = year - 1
		} else {
			month = month - 1
		}
	}

	/* 前面是倒序的，需要正过来 */
	var incomeData3 []FundIncomeData
	for i := 0; i < len(incomeData2); i++ {
		incomeData3 = append(incomeData3, incomeData2[len(incomeData2) - i - 1])
	}

	return incomeData3
}

func GetFundAccumulatedIncome(code string) (float64, float64, float64, float64) {
	incomeData := GetInComeData(code)
	if (len(incomeData) > 0) {
		
		units := incomeData[len(incomeData)-1].Units
		cost := incomeData[len(incomeData)-1].Cost

		accumulatedIncome := incomeData[len(incomeData)-1].AccumulatedIncome
		accumulatedIncomePercent := accumulatedIncome * 100 / incomeData[len(incomeData)-1].Cost

		return util.GetFloatFormat(accumulatedIncome, 2),
			   util.GetFloatFormat(accumulatedIncomePercent, 2),
			   util.GetFloatFormat(units, 2),
			   util.GetFloatFormat(cost, 2)
	}
	
	return 0.0, 0.0, 0.0, 0.0
}

func GetFundHandlingIncome(code string) (float64, float64) {
	incomeData := GetInComeData(code)
	if (len(incomeData) > 0) {
		handlingIncome := incomeData[len(incomeData)-1].HoldingIncome
		handlingIncomePercent := handlingIncome * 100 / incomeData[len(incomeData)-1].Cost

		if incomeData[len(incomeData)-1].Cost == 0 {
			handlingIncome = 0
			handlingIncomePercent = 0
		}

		return util.GetFloatFormat(handlingIncome, 2),
			   util.GetFloatFormat(handlingIncomePercent, 2)
	}
	
	return 0.0, 0.0
}