# fundgo

  ## 1 基金数据收集
  程序名：catch_data
  描述：从新浪基金数据网站获取数据，通过程序参数指定基金编码，获取的页面数，以及存储的csv文件。  
  参数：  
  -runserver 运行web服务  
  -catch 数据抓取  
  -a    读取所有的页数  
  -c string 基金编码  
  -d string csv文件输出目录 (default "./")   
  -p int 读取的页数 (default 1)  
  -s string 基金列表文件路径，采用此参数可以进行批量更新  