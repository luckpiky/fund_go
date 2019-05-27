package analyze

import (
	"fmt"
	"io/ioutil"
	"strings"
	"encoding/csv"
	"log"
	"strconv"
	"container/list"
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

		if (transData[i][0] != analyzeData.FundCode) {
			continue
		}

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

		fmt.Println(transData)
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