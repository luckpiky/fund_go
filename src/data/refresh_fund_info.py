import urllib.request
import json
import _thread
import time
import threading
import csv

printLock = threading.Lock()

def log(s):
    printLock.acquire()
    print(s)
    printLock.release()

def getUrlData(url):
    header = {
       'User-Agent': 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.96 Safari/537.36'
    }

    #print(time.strftime("%Y-%m-%d-%H_%M_%S", time.localtime()))
    request = urllib.request.Request(url, headers=header)
    reponse = urllib.request.urlopen(request).read()
    #print(time.strftime("%Y-%m-%d-%H_%M_%S", time.localtime()))
    #print("get url", url)
    return reponse


def parseFundData(data, fundCode):
    #print(data)
    json_str = json.loads(data)
    #print("parse end")
    #print(json_str['data']['list'][fundCode])
    fundName = json_str['data']['list'][fundCode]['fund_name']
    fundType = json_str['data']['list'][fundCode]['fundtype']
    fundRisk = json_str['data']['list'][fundCode]['riskmatch']['fundriskStr']
    companyId = json_str['data']['list'][fundCode]['company_id']
    #print(fundName,fundRisk,fundType,companyId)
    return fundName,fundType,fundRisk,companyId

def getFundCodeStr(fundCode):
    s = str(fundCode)
    padingLen = 6 - len(s)
    padingStr = ''
    for i in range(0, padingLen):
        padingStr = padingStr + '0'
    return padingStr + s


def getFundData(fundCodeStr):
    url = 'https://trade.xincai.com/api/getFundBuyStatus2?fund_code=' + fundCodeStr
    data = getUrlData(url)
    return parseFundData(data, fundCodeStr)


fundList = []
fundListLock = threading.Lock()

def addFundToList(fund):
    fundListLock.acquire()
    log("add fundinfo:"+fund['code'] + " " + fund['name'] + " " + fund['type'] + " " + fund['risk'] + " " + fund['company'])
    fundList.append(fund)
    fundListLock.release()

def readAndClearFundList(allFundInfo, fundList):
    fundListLock.acquire()
    for item in fundList:
        line = [item['code'], item['name'], item['type'], item['risk'], item['company']]
        allFundInfo[item['code']] = line
        #print(line)
    fundList.clear()
    fundListLock.release()

def printFundInfo(name, fundCode):
    #item = dict()
    #item['name'], item['type'], item['risk'], item['company'] = getFundData(fundCode)
    try:
        item = dict()
        item['code'] = fundCode
        item['name'], item['type'], item['risk'], item['company'] = getFundData(fundCode)
        #print(item)
        addFundToList(item)
    except:
        log("get info error:"+fundCode)
        pass
    
    #print(name,code,a,b,c)


lock = threading.Lock()
fundCodeList = []

def readFundInfo(csvFileName, allFundInfo):
    try:
        csvFile = open(csvFileName, 'r')
        readCSV = csv.reader(csvFile)
        for line in readCSV:
            if len(line) < 1:
                continue
            allFundInfo[line[0]] = line
        csvFile.close()
    except:
        log("open " + csvFileName + " fail")
        pass
    

firstWrite = True
allFundInfo = {}
def writeFundInfoToCsv():
    global firstWrite
    time.sleep(60)
    log("write file")
    csvFileName = 'fund_info.csv'
    if firstWrite:
        firstWrite = False
        readFundInfo(csvFileName, allFundInfo)
        #try:
        #    csvFile = open(csvFileName, 'r')
        #    readCSV = csv.reader(csvFile)
        #    for line in readCSV:
        #        if len(line) < 1:
        #            continue
        #        #print("test:",line)
        #        #print("ttt:",line[0])
        #        allFundInfo[line[0]] = line
        #    csvFile.close()
        #except:
        #    print("open " + csvFileName + " fail")
        #    pass
    readAndClearFundList(allFundInfo, fundList)
    #print(allFundInfo)
    allFundList = list(allFundInfo.values())
    with open(csvFileName,'w', newline = '') as file:
        writer = csv.writer(file, delimiter=',')
        #for line in allFundList:
        #    print("write:", line)
        writer.writerows(allFundList)
    log("write file end")

def writeCsvWorker(a, b):
    while True:
        writeFundInfoToCsv()

def GetFundCode():
    global fundCodeList
    global lock
    lock.acquire()
    
    if len(fundCodeList) <= 0:
        lock.release()
        return ""
    code = fundCodeList[0]
    del fundCodeList[0]
    #print(code)
    lock.release()
    return code

def initCode():
    #read fundinfo
    tmpFundInfo = {}
    readFundInfo('fund_info.csv', tmpFundInfo)

    # read all fund code
    fp = open("fund_code.csv", "r")
    for line in fp:
        code = line.strip()
        if tmpFundInfo.__contains__(code):
            #print("pass", code, "exsit:", tmpFundInfo[code])
            continue
        #print(code, "need parse")
        fundCodeList.append(code)
    print(len(fundCodeList), "codes need parse")
    return

def readFundWorker(name, i):
    while True:
        code = GetFundCode()
        #print("get code:", code)
        if code == "":
            print("get none, end")
            return
        printFundInfo(name, code)


def startThreads():
    catchThreadNum = 10
    for i in range(1, catchThreadNum):
        _thread.start_new_thread( readFundWorker, ("CatchThread-"+str(i), i) )
    _thread.start_new_thread( writeCsvWorker, ("Thread-6", 1) )


#printFundInfo("aaa", "000031")
initCode()
startThreads()

while True:
    time.sleep(1)
