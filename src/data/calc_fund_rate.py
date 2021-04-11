#coding=utf8

import datetime
import pandas as pd
import numpy as np
import ffn
import matplotlib.pyplot as plt
import csv
import os

class FundCalc:
    def __init__(self):
        self.code = ""
        self.path = "./"
        self.df = None
        self.rate = []
        return
    
    def setFundInfo(self, code, path):
        self.code = code
        self.path = path
        return
    
    def readData(self):
        self.df = pd.read_csv(self.path + '/' + self.code + ".csv", header=None)
        self.df.columns = ["date", "jjjz", "ljjz"] #读取到的csv文件没有表头，重新设置表头
        return
    
    def getMaxDrawDown(self):
        #计算每日的收益率
        self.df['rate'] = 0.0
        for i in range(1, len(self.df)):
            self.df["rate"][i] = format((self.df["ljjz"][i] - self.df["ljjz"][i - 1]) / self.df["jjjz"][i], ".2")
            pass
        
        r = pd.Series(self.df['rate'])
        value = (1 + r).cumprod() #扩展恢复以1为基数的净值（考虑分红情况）
        return ffn.calc_max_drawdown(value) * 100
    
    def getRateByDays(self, days):
        lastIndex = len(self.df['date']) - 1
        days = days - 1
        if lastIndex >= days:
            rate = (self.df['ljjz'][lastIndex] - self.df['ljjz'][lastIndex - days]) / self.df['jjjz'][lastIndex - days]
            rate = format(rate * 100, ".2f")
            return rate
        return None
    
    def calcRate(self):
        # 1 week
        rate = self.getRateByDays(5)
        if rate == None:
            return
        self.rate.append(rate)        
        
        # 1 month
        rate = self.getRateByDays(23)
        if rate == None:
            return
        self.rate.append(rate)
        
        # 3 month
        rate = self.getRateByDays(60)
        if rate == None:
            return
        self.rate.append(rate)        
        
        # 6 month
        rate = self.getRateByDays(125)
        if rate == None:
            return
        self.rate.append(rate)           
        
        # years
        year = 1
        while True:
            days = year * 250
            rate = self.getRateByDays(days)
            if rate == None:
                return
            self.rate.append(rate)
            year = year + 1
        return
    
PATH = "D:/code/fund_go/src/data/"
CSV_PATH = "D:/code/fund_go/src/data/csv3"
    
def readAllFund():
    fundName = {}
    csvFile = open(PATH + 'fund_info.csv', 'r')
    readCSV = csv.reader(csvFile)
    for line in readCSV:
        if len(line) < 3:
            continue
        fundName[line[0]] = line[1]
    csvFile.close()

    csvFile = open(PATH + 'fund_info_important.csv', 'r')
    readCSV = csv.reader(csvFile)
    for line in readCSV:
        if len(line) < 3:
            continue
        fundName[line[0]] = line[1]
    csvFile.close()
    
    return fundName

def addNewLine(df, fstLine, newLine):
    fillLen = len(fstLine) - len(newLine)
    
    # 不够的补齐
    if fillLen > 0:
        for i in range(fillLen):
            newLine.append("")
        
    # 超过的删除
    if fillLen < 0:
        newLine = newLine[0:len(fstLine)]
        
    return df.append(pd.DataFrame([newLine],columns=fstLine))
    
rateStr = ['code', 'name', 'MaxDrawDown', '1W', '1M', '3M', '6M', '1Y', '2Y', '3Y', '4Y', '5Y']

fundNameLst = readAllFund()
df1 = pd.DataFrame([],columns=rateStr)


fundLst = ["001158", "000031", "162006", "001316", "260108", "005217", "001204", "519062", "006862", "003297", "070020"]

def readFundLst(path):
    lst = []
    for home, dirs, files in os.walk(path):
        for filename in files:
            t = filename.split(".")
            lst.append(t[0])
    return lst

fundLst = readFundLst(CSV_PATH)
print(fundLst[1])

def calcNewFund(df, fundNameLst, fundCode):
    
    print("add", fundCode, fundNameLst[fundCode])
    
    calc = FundCalc()
    calc.setFundInfo(fundCode, "D:/code/fund_go/src/data/csv3") 
    
    newLine = None
    
    try:
        calc.readData()
        calc.calcRate()
        newLine = [fundCode, fundNameLst[fundCode], calc.getMaxDrawDown()] + calc.rate
    except:
        pass

    if newLine == None:
        return df
    
    return addNewLine(df1, rateStr, newLine)

count = 0
for fundCode in fundLst:
    df1 = calcNewFund(df1, fundNameLst, fundCode)
    count = count + 1
    if count % 20 == 0:
        print("write fund_calc.csv, count:"+ str(count))
        df1.to_csv(PATH + "fund_calc.csv", encoding='gbk')
        
print("write fund_calc.csv, count:"+ str(count))
df1.to_csv(PATH + "fund_calc.csv", encoding='gbk')
