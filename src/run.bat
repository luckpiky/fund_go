.\fundgo.exe -catch -s data/csv/funds.csv -d data/csv/ -p 5
python py/calc.py data/csv/ all -p 3
start chrome http://127.0.0.1:8000/index
python fundpy/manage.py runserver 0.0.0.0:8000