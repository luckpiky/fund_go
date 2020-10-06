from django.db import models
import csv
import os

# Create your models here.

def SaveTransationDataToFile():
    objs = Transation.objects.order_by('code', 'date')
    csvFile = open("data/csv/fund_transaction.csv", 'w', newline="")
    csvWriter = csv.writer(csvFile)
    csvWriter.writerow(['code', 'time', 'units', 'cost'])
    for item in objs:
        #print(item.code, item.date)
        dateStr = item.date.strftime('%Y-%m-%d 00:00:00')
        csvWriter.writerow([item.code, dateStr, item.units, item.cost])
    csvFile.close()
    os.system("python py/calc.py data/csv/ all")


class Transation(models.Model):
    code = models.CharField(verbose_name="基金编码", max_length=20, null=False)
    date = models.DateTimeField(verbose_name="交易日期", auto_now=False)
    cost = models.FloatField(verbose_name="交易金额")
    units = models.FloatField(verbose_name="交易份额")

    def save(self, *args, **kwargs):
        super().save(*args, **kwargs)
        SaveTransationDataToFile()
        return

