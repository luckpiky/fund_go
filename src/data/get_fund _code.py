import urllib.request
import json
import _thread
import time
import threading
import csv
import re

def getUrlData(url):
    header = {
       'User-Agent': 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.96 Safari/537.36'
    }

    #print(time.strftime("%Y-%m-%d-%H_%M_%S", time.localtime()))
    request = urllib.request.Request(url, headers=header)
    reponse = urllib.request.urlopen(request).read()
    #print(time.strftime("%Y-%m-%d-%H_%M_%S", time.localtime()))
    return reponse

fundCodes = {}
lock = threading.Lock()

def addCodeToList(code):
    lock.acquire()
    fundCodes[code] = code
    lock.release()

def getFundList(page):
    url = 'http://vip.stock.finance.sina.com.cn/fund_center/data/jsonp.php/IO.XSRV2.CallbackList[\'6XxbX6h4CED0ATvW\']/NetValue_Service.getNetValueOpen?num=40&sort=nav_date&asc=0&ccode=&page='
    url = url + str(page)
    data = getUrlData(url)
    s = data.decode('ascii')
    s = s.replace('"', '')
    #print(s)
    lst = re.findall(r'symbol:{0,1}[\d]*' , s)
    print(lst)
    for item in lst:
        code = item.split(':')
        if len(code) > 1:
            addCodeToList(code[1])


threadNum = 10
maxPage = 284
step = int(maxPage / threadNum)
threadEndFlag = [0 for i in range(threadNum)]
print(threadEndFlag)

def catchWorker(index, start):
    end = start + step
    if end > maxPage:
        end = maxPage
    for i in range(start, end + 1):
        print("page:", i)
        getFundList(i)
    threadEndFlag[index] = 1


for i in range(threadNum):
    start = i * step + 1
    _thread.start_new_thread( catchWorker, (i, start) )

while True:
    sumThread = 0
    for i in range(threadNum):
        sumThread = sumThread + threadEndFlag[i]
    if sumThread == threadNum:
        break
    else:
        time.sleep(1)

print('end')
print(list(fundCodes))

fp = open("fund_code.csv", 'w')
for code in fundCodes:
    fp.write(code+"\n")
fp.close()

