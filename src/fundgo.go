package main

import (
	"flag"
	"github.com/astaxie/beego"
	//f "./analyze"
	catcher "./catch"
	_ "./beego/router"
)

var catch = flag.Bool("catch", false, "数据抓取")
var runserver = flag.Bool("runserver", false, "运行服务")

func main() {

	flag.Parse()

	if *catch == true {
		catcher.CatchDataMain()
	} else if *runserver == true {
		beego.SetStaticPath("js","views/js")
		beego.Run()
	}
}