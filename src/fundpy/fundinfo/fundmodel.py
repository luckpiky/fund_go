class FundPrice:
    jjjz = 0.0 #基金净值
    ljjz = 0.0 #累计净值
    rate = 0.0 #日增长率

class FundIncome:
    units = 0.0 #基金份额
    cost = 0.0 # 平均成本
    totalCost = 0.0 # 总成本
    incomeTotal = 0.0 # 累计收益
    income = 0.0 #基金收益
    incomePercent = 0.0 #收益率
    incomeDay = 0.0

class FundTransData:
    date = ""
    price = None
    income = None

class FundInfo:
    code = ""
    name = ""
    transData = []

    def __init__(self):
        transData = []