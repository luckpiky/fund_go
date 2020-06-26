import csv
import time
import os
import sys

def getTimeStamp(timeStr):
    timeArray = time.strptime(timeStr, "%Y-%m-%d %H:%M:%S")
    timeStamp = int(time.mktime(timeArray))
    return timeStamp

def getTimeStr(timeStamp):
    timeArray = time.localtime(timeStamp)
    timeStr = time.strftime("%Y-%m-%d %H:%M:%S", timeArray)
    return timeStr

class FundPrice:
    date = ""
    jjjz = 0.0 #基金净值
    ljjz = 0.0 #累计净值
    rate = 0.0 #日增长率
    rateYear = 0.0 #最近1年的收益
    rateYear3 = 0.0 #最近3年的收益
    rateYear5 = 0.0 #最近5年的收益

class FundIncome:
    units = 0.0 #基金份额
    cost = 0.0 # 平均成本
    totalCost = 0.0 # 总成本
    incomeTotal = 0.0 # 累计收益
    income = 0.0 #基金收益
    incomePercent = 0.0 #收益率
    incomeDay = 0.0 #当日收益

class FundTransData:
    date = ""
    units = 0.0
    cost = 0.0

class FundInfo:
    baseDir = ""
    code = ""
    price = []
    transData = []
    income = []

    def __init__(self, code):
        self.code = code
        self.price = []
        self.transData = []
        self.income = []
        return

    def setBaseDir(self, dir):
        self.baseDir = dir
        return

    def readBasicData(self):
        path = self.baseDir + self.code + ".csv"
        csvReader = csv.reader(open(path, encoding='utf-8'))
        count = 0
        jjjz = 0.0
        ljjz = 0.0
        for line in csvReader:
            count = count + 1
            priceTmp = FundPrice()
            priceTmp.date = getTimeStamp(line[0])
            priceTmp.jjjz = float(line[1])
            priceTmp.ljjz = float(line[2])
            priceTmp.rate = self.calcRate(jjjz, ljjz, priceTmp.ljjz)
            self.price.append(priceTmp)
            jjjz = priceTmp.jjjz
            ljjz = priceTmp.ljjz
            #print(getTimeStr(priceTmp.date), priceTmp.jjjz, priceTmp.ljjz, priceTmp.rate)
        return

    def readTransData(self):
        path = self.baseDir + "fund_transaction.csv"
        csvReader = csv.reader(open(path, encoding='utf-8'))
        lineCount = 0
        for line in csvReader:
            lineCount = lineCount + 1
            if lineCount == 1:
                continue

            if line[0] != self.code:
                continue

            transDataTmp = FundTransData()
            transDataTmp.date = getTimeStamp(line[1])
            transDataTmp.units = float(line[2])
            transDataTmp.cost = float(line[3])
            self.transData.append(transDataTmp)
            #print(line)
        return

    def readData(self):
        self.readTransData()
        self.readBasicData()
        return

    def calcRate(self, jjjz, ljjz1, ljjz2):
        if jjjz == 0.0:
            return 0.0

        rate = (ljjz2 - ljjz1) * 100 / jjjz
        return round(rate, 2)

    # 根据交易记录计算收益
    def calcIncome(self):
        units = 0.0
        totalCost = 0.0
        cost = 0.0
        transIncome = 0.0
        lastDayLjjz = 0.0

        for price in self.price:
            incomeDay = 0.0

            # 每日收益
            if lastDayLjjz == 0.0:
                incomeDay = 0.0
            else:
                incomeDay = (price.ljjz - lastDayLjjz) * units
            lastDayLjjz = price.ljjz

            for trans in self.transData:
                if price.date == trans.date:
                    # 添加交易数据
                    units = units + trans.units
                    if trans.units >= 0:
                        totalCost = totalCost + trans.cost
                    else :
                        totalCost = cost * units
                        transIncome = transIncome + trans.units * cost - trans.cost
                        #print("transIncome=", round(transIncome, 3), "trans.cost=", round(trans.cost, 3), "trans.units=", round(trans.units, 3))
                    #print("buy units=", round(trans.units, 3), "cost=", round(cost, 3), "total units=", round(units, 3), "totalCost=", round(totalCost, 3), "transIncome", round(transIncome, 3))
                    pass

            # 成本计算:
            # 买入100*10，成本 = 100
            # 第二次买入200*10，成本 = 100*10+200*10 = 3000,平均 = 150
            # 卖出300*10，剩下的成本 = 150*10 = 1500
            income = FundIncome()
            income.units = units
            income.totalCost = totalCost
            if units == 0:
                income.cost = 0.0
            else:
                income.cost = totalCost / units
            cost = income.cost

            # 计算收益
            income.income = price.jjjz * income.units - income.totalCost
            income.incomeTotal = income.income + transIncome
            if income.totalCost > 0:
                income.incomePercent = income.income * 100 / income.totalCost
            else :
                income.incomePercent = 0.0

            # 格式化
            income.income = round(income.income, 2)
            income.totalCost = round(income.totalCost, 2)
            income.incomePercent = round(income.incomePercent, 2)
            income.cost = round(income.cost, 2)
            if income.units == 0.0:
                income.cost = 0.0
            income.units = round(income.units, 2)
            income.incomeTotal = round(income.incomeTotal, 2)
            income.incomeDay = round(incomeDay, 2)

            #print(income.totalCost, income.income, income.incomePercent, income.incomeTotal)

            self.income.append(income)

        return

    def writeIncomeData(self):
        filename = self.baseDir + "income_" + self.code + ".csv"
        with open(filename, "w", newline='') as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow(["date", "jjjz", "ljjz", "rate", "units", "totalCost", \
                "cost", "incomeTotal", "income", "incomePercent", "incomeDay"])

            index = 0
            for price in self.price:
                income = self.income[index]
                writer.writerow([getTimeStr(price.date), price.jjjz, price.ljjz, price.rate, income.units, income.totalCost,\
                    income.cost, income.incomeTotal, income.income, income.incomePercent, income.incomeDay])
                                
                index = index + 1
        return
    
    def calcYearRate(self):
        # 一年按照254天交易日算，不到一年的数据按照比例计算
        Y1 = 245
        Y3 = 245 * 3
        Y5 = 245 * 5
        count = len(self.price)

        rateYear1 = 0.0
        rateYear3 = 0.0
        rateYear5 = 0.0

        index = count - 1

        if count < Y1:
            rateYear1 = (self.price[index].ljjz - self.price[0].ljjz) * Y1 * 100 / (count * self.price[0].jjjz)
            rateYear3 = rateYear1 * 3
            rateYear5 = rateYear1 * 5
        elif count < Y3:
            rateYear1 = (self.price[index].ljjz - self.price[count - Y1].ljjz) * 100 / self.price[count - Y1].jjjz
            rateYear3 = (self.price[index].ljjz - self.price[0].ljjz) * Y3 * 100 / (count * self.price[0].jjjz)
            rateYear5 = rateYear3 * 5 / 3
        elif count < Y5:
            rateYear1 = (self.price[index].ljjz - self.price[count - Y1].ljjz) * 100 / self.price[count - Y1].jjjz
            rateYear3 = (self.price[index].ljjz - self.price[count - Y3].ljjz) * 100 / self.price[count - Y3].jjjz
            rateYear5 = (self.price[index].ljjz - self.price[0].ljjz) * Y5 * 100 / (count * self.price[0].jjjz)
        else:
            #print(count, self.price[index].ljjz, self.price[count - Y1].ljjz, self.price[count - Y1].jjjz, getTimeStr(self.price[count - Y1].date))
            rateYear1 = (self.price[index].ljjz - self.price[count - Y1].ljjz) * 100 / self.price[count - Y1].jjjz
            rateYear3 = (self.price[index].ljjz - self.price[count - Y3].ljjz) * 100 / self.price[count - Y3].jjjz
            rateYear5 = (self.price[index].ljjz - self.price[count - Y5].ljjz) * 100 / self.price[count - Y5].jjjz     

        return round(rateYear1, 2), round(rateYear3, 2), round(rateYear5, 2)

    def monitorRate(self):

        return

class FundRateMonitor:
    monitorDay = 10
    monitorRate = -5.0

    monitorResult = []

    def __init__(self):
        monitorResult = []
        return

    def rateMonitor(self, price):
        index = len(price) - 1
        if len(price) < self.monitorDay:
            return 0.0, False
        
        rate = (price[index].ljjz - price[index - self.monitorDay].ljjz) * 100 / price[index - self.monitorDay].jjjz
        rate = round(rate, 2)
        if rate <= self.monitorRate:
            return rate, True

        return rate, False

    def monitor(self, price, code):
        rate, result = self.rateMonitor(price)
        self.monitorResult.append([code, rate, result])
        return
    
    def writeMonitorResult(self, baseDir):
        filename = baseDir + "monitor_rate.csv"
        with open(filename, "w", newline='') as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow(["code", "rate", "result"])
            writer.writerows(self.monitorResult)
        return

class FundBasicInfo:
    code = ""
    name = ""
    type = ""
    riskLevel = ""
    incomeTotal = 0.0
    income = 0.0
    incomePercent = 0.0
    cost = 0.0

    def __init__(self, code):
        self.code = code
        incomeTotal = 0.0
        income = 0.0
        incomePercent = 0.0
        cost = 0.0
        return

    def getBasicInfo(self, dir):
        path = dir + "funds.csv"
        csvReader = csv.reader(open(path, encoding='utf-8'))
        for item in csvReader:
            if self.code == item[0]:
                self.name = item[1]
                self.type = item[2]
                self.riskLevel = item[3]
                return
        return

def calcAll(baseDir):
    fundList = []
    path = baseDir + "fund_transaction.csv"
    csvReader = csv.reader(open(path, encoding='utf-8'))
    lineCount = 0
    rateMonitor = FundRateMonitor()
    fundInfoList = []

    for line in csvReader:
        if lineCount == 0:
            lineCount = lineCount + 1
            continue

        find = False
        fundCode = line[0]
        for code in fundList:
            if fundCode == code:
                find = True
                break
        if find == True:
            continue

        fundList.append(fundCode)
        basicInfo = FundBasicInfo(fundCode)
        basicInfo.getBasicInfo(baseDir)
        fundInfoList.append(basicInfo)

    rateYear = []

    for code in fundList:
        print("calc code =", code)
        fundInfo = FundInfo(code)
        fundInfo.setBaseDir(baseDir)
        fundInfo.readData()
        fundInfo.calcIncome()
        fundInfo.writeIncomeData()
        rateY1, rateY3, rateY5 = fundInfo.calcYearRate()
        rateYear.append([code, rateY1, rateY3, rateY5])
        rateMonitor.monitor(fundInfo.price, code)

        for item in fundInfoList:
            if item.code == code:
                index = len(fundInfo.income) - 1
                income = fundInfo.income[index]
                item.incomeTotal = income.incomeTotal
                item.income = income.income
                item.incomePercent = income.incomePercent
                item.cost = income.totalCost
                print(item.cost)

    filename = baseDir + "rate_year.csv"
    with open(filename, "w", newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(["code", "rateY1", "rateY3", "rateY5"])
        writer.writerows(rateYear)

    filename = baseDir + "index_info.csv"
    with open(filename, "w", encoding='utf-8', newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(["code", "name", "type", "riskLevel", "incomeTotal", "income", "incomePercent", "cost"])
        for item in fundInfoList:
            writer.writerow([item.code, item.name, item.type, item.riskLevel,\
                 item.incomeTotal, item.income, item.incomePercent, item.cost])


    rateMonitor.writeMonitorResult(baseDir)
    return
        

if __name__ == '__main__':
    if sys.argv == None:
        print("python calc.py basedir code")

    dir = sys.argv[1]

    if sys.argv[2] == "all":
        calcAll(dir)
    else:
        code = sys.argv[2]

        fundInfo = FundInfo(code)
        fundInfo.setBaseDir(dir)
        fundInfo.readData()
        fundInfo.calcIncome()
        fundInfo.writeIncomeData()

