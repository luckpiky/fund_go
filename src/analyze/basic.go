package analyze

import (
	"fmt"
	"io/ioutil"
	"strings"
	"encoding/csv"
	"log"
	"strconv"
	"../util"
	//"container/list"
)

var dataPath = "data/csv"

var AllFundsList map[string] []string
var MyFundsList map[string] []string


type FundPriceInfo struct {
	Date string
	Jjjz float64
	Ljjz float64
}

type FundGrowthRate struct {
	Date string
	Rate float64
}

type FundTrans struct {
	Date string
	Code string
	Units float64 //交易份额
	Amount float64  //交易金额
}

func SetDataPath(path string) {
	dataPath = path
}

func GetDataPath() string {
	return dataPath
}

func GetAllFundsList() {
	var fundsData [][] string

	fileName := dataPath + "/funds.csv"

	readContent,err := ioutil.ReadFile(fileName)
    if err == nil {
        csvReadFp := csv.NewReader(strings.NewReader(string(readContent)))
        fundsData,_ = csvReadFp.ReadAll()
        log.Println(fileName, "读取到基金数据条目数：",len(fundsData))
    } else {
		log.Println(err)
	}

	//第一行数据跳过
	for i := 1; i < len(fundsData); i++ {
		AllFundsList[fundsData[i][0]] = []string{fundsData[i][1], fundsData[i][2]}
	}
}

func GetFundBasicInfoByCode(code string) (info []string) {
	if _,ok := AllFundsList[code]; ok {
		return AllFundsList[code]
	}

	return []string{"", ""}
}

func GetMyFundsList() {
	var transData [][] string

	fileName := dataPath + "/fund_transaction.csv"

	readContent,err := ioutil.ReadFile(fileName)
    if err == nil {
        csvReadFp := csv.NewReader(strings.NewReader(string(readContent)))
        transData,_ = csvReadFp.ReadAll()
        log.Println(fileName, "读取到基金交易数据条目数：",len(transData))
    } else {
		log.Println(err)
	}

	//第一行数据跳过
	for i := 1; i < len(transData); i++ {
		MyFundsList[transData[i][0]] = GetFundBasicInfoByCode(transData[i][0])
	}
}

func GetFundPriceByCode(code string) (price []FundPriceInfo) {
	var fundPrice [][] string
	var fundPriceData []FundPriceInfo

	fileName := dataPath + "/" + code + ".csv"

	readContent,err := ioutil.ReadFile(fileName)
    if err == nil {
        csvReadFp := csv.NewReader(strings.NewReader(string(readContent)))
        fundPrice,_ = csvReadFp.ReadAll()
        log.Println(fileName, "读取到基金价格数据条目数：",len(fundPrice))
    } else {
		log.Println(err)
	}

	for i := 0; i < len(fundPrice); i++ {
		var price FundPriceInfo
		price.Date = fundPrice[i][0]
		price.Jjjz,_ = strconv.ParseFloat(fundPrice[i][1], 64)
		price.Ljjz,_ = strconv.ParseFloat(fundPrice[i][2], 64)
		fundPriceData = append(fundPriceData, price)
	}

	return fundPriceData
}

func GetGrowthRateByCode(code string) (rate []FundGrowthRate) {
	var growthRate []FundGrowthRate

	priceList := GetFundPriceByCode(code)

	lastJjjz := 0.0
	lastLjjz := 0.0

	for i := 0; i < len(priceList); i++ {
		if i == 0 {
			var rate FundGrowthRate
			rate.Date = priceList[i].Date
			rate.Rate = 0
			growthRate = append(growthRate, rate)
		} else {
			var rate FundGrowthRate
			rate.Date = priceList[i].Date
			rate.Rate = (priceList[i].Ljjz - lastLjjz) / lastJjjz * 100
			rate.Rate = util.GetFloatFormat(rate.Rate, 3)
			growthRate = append(growthRate, rate)
		}

		lastJjjz = priceList[i].Jjjz
		lastLjjz = priceList[i].Ljjz
	}

	return growthRate
}

func GetGrowthRateFromBeginByCode(code string, startDate string) (rate []FundGrowthRate) {
	var growthRate []FundGrowthRate

	priceList := GetFundPriceByCode(code)
	startIndex := 0
	start := 0

	for i := 0; i < len(priceList); i++ {
		if priceList[i].Date == startDate {
			var rate FundGrowthRate
			rate.Date = priceList[i].Date
			rate.Rate = 0
			growthRate = append(growthRate, rate)
			startIndex = i
			start = 1
			continue
		} 
		
		if start == 0 {
			continue
		}		
		
		var rate FundGrowthRate
		rate.Date = priceList[i].Date
		rate.Rate = (priceList[i].Ljjz - priceList[startIndex].Ljjz) / priceList[startIndex].Jjjz * 100
		rate.Rate = util.GetFloatFormat(rate.Rate, 3)
		growthRate = append(growthRate, rate)
	}

	return growthRate
}


func GetTransData(code string) ([]FundTrans) {
	var transList []FundTrans
	var transData [][] string

	fileName := dataPath + "/fund_transaction.csv"

	readContent,err := ioutil.ReadFile(fileName)
    if err == nil {
        csvReadFp := csv.NewReader(strings.NewReader(string(readContent)))
        transData,_ = csvReadFp.ReadAll()
        log.Println(fileName, "读取到基金交易数据条目数：",len(transData))
    } else {
		log.Println(err)
	}

	for i := 1; i < len(transData); i++ {
		if code != "" && code != transData[i][0] {
			continue
		}

		var trans FundTrans
		trans.Date = transData[i][1]
		trans.Code = transData[i][0]
		trans.Units,_ = strconv.ParseFloat(transData[i][2], 64)
		trans.Amount,_ = strconv.ParseFloat(transData[i][3], 64)
		transList = append(transList, trans)
	}

	return transList
}

func init() {

	AllFundsList = make(map[string] []string)
	MyFundsList = make(map[string] []string)

	GetAllFundsList()
	GetMyFundsList()

	fmt.Println("读取到的基金列表：")
	for item := range AllFundsList {
		fmt.Println(item, AllFundsList[item])
	}

	fmt.Println("读取到的基金交易列表：")
	for item := range MyFundsList {
		fmt.Println(item)
	}
}

type FundTransData2 struct {
	Date int64
	Code string
	Units float64 //交易份额
	Amount float64  //交易金额
}

type FundPriceData2 struct {
	Date int64
	Jjjz float64
	Ljjz float64
}

func GetFundPriceData(code string) (price []FundPriceData2) {
	var fundPrice [][] string
	var fundPriceData []FundPriceData2

	fileName := dataPath + "/" + code + ".csv"

	readContent,err := ioutil.ReadFile(fileName)
    if err == nil {
        csvReadFp := csv.NewReader(strings.NewReader(string(readContent)))
        fundPrice,_ = csvReadFp.ReadAll()
        log.Println(fileName, "读取到基金价格数据条目数：",len(fundPrice))
    } else {
		log.Println(err)
	}

	for i := 0; i < len(fundPrice); i++ {
		var price FundPriceData2
		price.Date = util.TimeStr2Int64(fundPrice[i][0])
		price.Jjjz,_ = strconv.ParseFloat(fundPrice[i][1], 64)
		price.Ljjz,_ = strconv.ParseFloat(fundPrice[i][2], 64)
		fundPriceData = append(fundPriceData, price)
	}

	return fundPriceData
}

func GetTransData2(code string) ([]FundTransData2) {
	var transList []FundTransData2
	var transData [][] string

	fileName := dataPath + "/fund_transaction.csv"

	readContent,err := ioutil.ReadFile(fileName)
    if err == nil {
        csvReadFp := csv.NewReader(strings.NewReader(string(readContent)))
        transData,_ = csvReadFp.ReadAll()
        log.Println(fileName, "读取到基金交易数据条目数：",len(transData))
    } else {
		log.Println(err)
	}

	for i := 1; i < len(transData); i++ {
		if code != "" && code != transData[i][0] {
			continue
		}

		var trans FundTransData2
		trans.Date = util.TimeStr2Int64(transData[i][1])
		trans.Code = transData[i][0]
		trans.Units,_ = strconv.ParseFloat(transData[i][2], 64)
		trans.Amount,_ = strconv.ParseFloat(transData[i][3], 64)
		transList = append(transList, trans)
	}

	return transList
}