import urllib.request
import json
import _thread
import time
import threading
import csv

MAX_PAGE = 130
DELAY_PAGE = 5
DELAY_TIME = 5
DELAY_FUND_COUNT = 5
DELAY_FUND_TIME = 10

printLock = threading.Lock()

def log(s):
    printLock.acquire()
    print(s)
    printLock.release()

def getFilePath(code, path):
    return path + '/' + code + '.csv'

def getSortData(item):
    return item[0]

def writeData(code, path, priceLst):
    fileName = getFilePath(code, path)
    data = list(priceLst.values())
    data.sort(key=getSortData)
    with open(fileName,'w', newline = '') as file:
        writer = csv.writer(file, delimiter=',')
        writer.writerows(data)
        log("write " + fileName + " success")
    return

def readDataFromCsv(code, path, priceLst):
    fileName = getFilePath(code, path)
    #print('read', fileName)
    try:
        csvFile = open(fileName, 'r')
        readCSV = csv.reader(csvFile)
        for line in readCSV:
            if len(line) < 1:
                continue
            priceLst[line[0]] = line
            #print('read data:', priceLst)
        csvFile.close()
    except:
        log("open " + fileName + " fail")
        pass

def getUrlData(url):
    header = {
       'User-Agent': 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.96 Safari/537.36'
    }

    reponse = ""
    try:
        request = urllib.request.Request(url, headers=header)
        reponse = urllib.request.urlopen(request).read()
    except:
        log("read url fail, wait 5 mins:" + url)
        time.sleep(300)
        request = urllib.request.Request(url, headers=header)
        reponse = urllib.request.urlopen(request).read()        
        pass
    return reponse


def parseData(data, priceLst):
    count = 0
    json_str = json.loads(data)
    data = json_str['result']['data']
    try:
        for item in data['data']:
            if priceLst.__contains__(item['fbrq']):
                count = count + 1
                #print("find data")
            priceLst[item['fbrq']] = list(item.values())
            #print(item)
    except:
        log("parse fail for " + code + ":" + data)
        return True
    if count > 10: #if 10 items of one page are found same, it means this page has been read
        return True
    return False


def getData(code, page, priceLst):
    url = 'http://stock.finance.sina.com.cn/fundInfo/api/openapi.php/CaihuiFundInfoService.getNav?symbol=CODE&page=PAGE'
    url = url.replace('CODE', code)
    url = url.replace('PAGE', str(page))
    data = getUrlData(url)
    return parseData(data, priceLst)


def getNewData(code, path):
    log("refresh " + code)
    priceLst = {}
    readDataFromCsv(code, path, priceLst)
    for i in range(1, MAX_PAGE):
        isExsit = getData(code, i, priceLst)
        if isExsit:
            log("read same for " + code + ", break")
            break
        if i % DELAY_PAGE == 0:
            time.sleep(DELAY_TIME)
            log("sleep a moment for " + code)
    writeData(code, path, priceLst)
    return

def getAllData(code, path):
    priceLst = {}
    for i in range(1, MAX_PAGE):
        getFundData(code, i, priceLst)
    writeData(code, path, priceLst)
    return


allFundCode = []
def readAllCode():
    fileName = "fund_info.csv"
    try:
        csvFile = open(fileName, 'r')
        readCSV = csv.reader(csvFile)
        for line in readCSV:
            if len(line) < 3:
                continue
            if line[3] == '高风险' or line[3] == '低风险':
                #log("skip:" + str(line))
                continue
            if line[2] == 'U':
                continue

            allFundCode.append(line[0])
        csvFile.close()

        csvFile = open(fileName2, 'r')
        readCSV = csv.reader(csvFile)
        for line in readCSV:
            if len(line) < 1:
                continue
            allFundCode.append(line[0])
        csvFile.close()
    except:
        log("read " + fileName + " fail")
        pass

    fileName = "fund_info_important.csv"
    try:
        csvFile = open(fileName, 'r')
        readCSV = csv.reader(csvFile)
        for line in readCSV:
            if len(line) < 1:
                continue
            allFundCode.append(line[0])
        csvFile.close()
    except:
        log("read " + fileName + " fail")
        pass

    log("get code counts:" + str(len(allFundCode)))
    return allFundCode

lock = threading.Lock()
def readCode():
    lock.acquire()
    if len(allFundCode) <= 0:
        lock.release()
        return ""
    index = len(allFundCode) - 1
    code = allFundCode[index]
    del allFundCode[index]
    lock.release()
    return code

def refreshWoker(name, path):
    count = 0
    while True:
        code = readCode()
        if code == "":
            return
        count = count + 1
        if count % DELAY_FUND_COUNT == 0:
            time.sleep(DELAY_FUND_TIME)
        getNewData(code, path)
    return

def startThreads(threadNum, path):
    for i in range(threadNum):
        _thread.start_new_thread( refreshWoker, ("refreshThread-"+str(i), path) )
    return

readAllCode()
startThreads(5, './csv3/')

