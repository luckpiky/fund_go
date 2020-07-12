.\fundgo.exe -catch -s data/csv/funds.csv -d data/csv/ -p 5
python py/calc.py data/csv/ all
python fundpy/manage.py runserver