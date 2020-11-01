from django.shortcuts import render
from django.http import HttpResponse
from django.conf import settings
import csv
import time
from fundinfo.fundmodel import *

def readTransData(code):
    path = settings.CSV_DIR + "income_" + code + ".csv"
    csvReader = None
    try:
        csvReader = csv.reader(open(path, encoding='utf-8'))
    except:
        return None
    fundinfo = FundInfo()
    fundinfo.transData = []
    fundinfo.code = code
    count = 0
    canAdd = False
    for line in csvReader:
        if count == 0:
            count = count + 1
            continue

        # 跳过没有购买的日期
        if float(line[4]) == 0.0 and canAdd != True:
            continue
        canAdd = True

        trans = FundTransData()
        trans.price = FundPrice()
        trans.income = FundIncome()
        trans.date = line[0]
        trans.price.jjjz = float(line[1])
        trans.price.ljjz = float(line[2])
        trans.price.rate = float(line[3])
        trans.income.units = float(line[4])
        trans.income.totalCost = float(line[5])
        trans.income.cost = float(line[6])
        trans.income.incomeTotal = float(line[7])
        trans.income.income = float(line[8])
        trans.income.incomePercent = float(line[9])
        trans.income.incomeDay = float(line[10])
        fundinfo.transData.append(trans)
    return fundinfo

def getRateYear(code):
    path = settings.CSV_DIR + "rate_year.csv"
    print(path)
    csvReader = None
    try:
        csvReader = csv.reader(open(path, encoding='utf-8'))
    except:
        return None
    for item in csvReader:
        if item[0] == code:
            return float(item[1]), float(item[2]), float(item[3])

def getRateMonitor(code):
    path = settings.CSV_DIR + "monitor_rate.csv"
    csvReader = None
    try:
        csvReader = csv.reader(open(path, encoding='utf-8'))
    except:
        return None
    for item in csvReader:
        if item[0] == code:
            #print("-------------", float(item[1]), item[3])
            return float(item[1]), item[3], item[2]

def getMonthStr(day):
    tmp = day.split("-")
    return tmp[0] + tmp[1]

def getMonthIncome(transData):
    lastMonth = ""
    income = 0.0
    incomeMonth = []
    startAdd = False
    for item in transData:
        day = item.date
        month = getMonthStr(day)

        if startAdd == False and item.income.incomeDay != 0.0:
            startAdd = True

        if lastMonth != month:
            if lastMonth != "" and startAdd:
                incomeMonth.append([lastMonth, round(income, 2)])
            income = item.income.incomeDay
            lastMonth = month
        else:
            income = income + item.income.incomeDay
    incomeMonth.append([lastMonth, round(income, 2)])
    return incomeMonth

def readFundBasicInfo():
    path = settings.CSV_DIR + "index_info.csv"
    csvReader = None
    try:
        csvReader = csv.reader(open(path, encoding='utf-8'))
    except:
        return None

    fundlist = []

    count = 0
    for item in csvReader:
        if count == 0:
            count = count + 1
            continue
        fundlist.append(item)
    return fundlist

def getFundBasicInfo(code):
    info = readFundBasicInfo()
    for item in info:
        if item[0] == code:
            return item
    return None

# Create your views here.
def getFundInfo(request, code):
    #code = request.GET.get("code")
    fundinfo = readTransData(code)
    if fundinfo == None:
        return HttpResponse("Not find the code " + code + " !")
    info = {}
    info['code'] = code

    incomeMonth = getMonthIncome(fundinfo.transData)

    last = len(fundinfo.transData) - 1

    info['info'] = getFundBasicInfo(code)

    info['incomeTotal'] = fundinfo.transData[last].income.incomeTotal
    info['income'] = fundinfo.transData[last].income.income
    info['incomePercent'] = fundinfo.transData[last].income.incomePercent
    info['units'] = fundinfo.transData[last].income.units
    info['cost'] = fundinfo.transData[last].income.cost
    info['totalCost'] = fundinfo.transData[last].income.totalCost
    info['price'] = fundinfo.transData[last].price
    info['trans'] = fundinfo.transData[last]
    info['data'] = fundinfo.transData
    info['incomeMonth'] = incomeMonth
    print(info['info'], info['info'][7],info['totalCost'])
    info['fundCost'] = round(info['totalCost'] / fundinfo.transData[last].income.units, 3)
    
    rateY1, rateY3, rateY5 = getRateYear(code)
    info['rateY1'] = rateY1
    info['rateY3'] = rateY3
    info['rateY5'] = rateY5

    rate, result, days = getRateMonitor(code)
    info['rateMonitor'] = rate
    info['rateMonitorResult'] = result
    info['rateMonitorDays'] = days

    return render(request,"fundinfo.html", info)

def getFundTypesInfo(fundList, totalCost):
    types = []
    for item in fundList:
        canAdd = True
        for type in types:
            if type[0] == item[2]:
                canAdd = False
                type[1] = type[1] + float(item[7])
                type[2] = type[2] + float(item[4])
                break
        if canAdd:
            types.append([item[2], float(item[7]), float(item[4]), 0.0])

    for item in types:
        item[1] = round(item[1], 2)
        item[2] = round(item[2], 2)
        item[3] = round(item[2] * 100 / item[1], 2)
    return types

def getFundRiskInfo(fundList, totalCost):
    risks = []
    for item in fundList:
        canAdd = True
        for risk in risks:
            if risk[0] == item[3]:
                canAdd = False
                risk[1] = risk[1] + float(item[7])
                risk[2] = risk[2] + float(item[4])
                break
        if canAdd:
            risks.append([item[3], float(item[7]), float(item[4]), 0.0])

    for item in risks:
        item[1] = round(item[1], 2)
        item[2] = round(item[2], 2)
        item[3] = round(item[2] * 100 / item[1], 2)
    return risks

def getFundCostByTypeOrder(types, fundList):
    costs = []
    for type in types:
        for fund in fundList:
            if type[0] == fund[2] and float(fund[7]) != 0.0:
                costs.append([fund[1], fund[7]])
    return costs

def getFundName(code, fundList):
    for item in fundList:
        if item[0] == code:
            return item[1]
    return ""

def getWarning(fundList):
    info = ""
    path = settings.CSV_DIR + "monitor_rate.csv"
    csvReader = None
    try:
        csvReader = csv.reader(open(path, encoding='utf-8'))
    except:
        return None
    for item in csvReader:
        if item[3] == 'True':
            info = info + getFundName(item[0], fundList) + " " + item[1] + " " + item[2] + "days;"
    return info

def fundListOrderKey(fundList):
    return float(fundList[7])

def getMonthInome():
    income = []
    path = settings.CSV_DIR + "month_income.csv"
    csvReader = None
    try:
        csvReader = csv.reader(open(path, encoding='utf-8'))
    except:
        return None

    for item in csvReader:
        income.append([item[0], float(item[1])])
    return income

def getIndex(request):

    fundList = readFundBasicInfo()
    fundList.sort(key=fundListOrderKey, reverse=True)

    print(len(fundList))

    info = {}
    info['funds'] = fundList

    # 基金数量
    info['fundNum'] = len(fundList)

    # 成本
    cost = 0.0
    for item in fundList:
        cost = cost + float(item[7])
    info['cost'] = round(cost, 2)  

    # 累计收益
    incomeTotal = 0.0
    for item in fundList:
        incomeTotal = incomeTotal + float(item[4])
    info['incomeTotal'] = round(incomeTotal, 2)
    info['incomeTotalPercent'] = round(info['incomeTotal'] * 100 / info['cost'], 2)

    # 持有收益
    income = 0.0
    for item in fundList:
        income = income + float(item[5])
    info['income'] = round(income, 2)
    info['incomePercent'] = round(info['income'] * 100 / info['cost'], 2)

    info['types'] = getFundTypesInfo(fundList, info['cost'])
    info['risks'] = getFundRiskInfo(fundList, info['cost'])
    info['typeItems'] = getFundCostByTypeOrder(info['types'], fundList)

    info['warning'] = getWarning(fundList)

    info['monthIncome'] = getMonthInome()

    info['lastYearIncome'] = 0.0
    cnt = len(info['monthIncome'])
    for i in range(cnt- 12, cnt):
        if i < 0:
            continue
        info['lastYearIncome'] = info['lastYearIncome'] + info['monthIncome'][i][1]
    info['lastYearIncome'] = round(info['lastYearIncome'], 2)

    info['last2YearIncome'] = 0.0
    for i in range(cnt- 24, cnt):
        if i < 0:
            continue
        info['last2YearIncome'] = info['last2YearIncome'] + info['monthIncome'][i][1]
    info['last2YearIncome'] = round(info['last2YearIncome'], 2)

    return render(request,"index.html", info)