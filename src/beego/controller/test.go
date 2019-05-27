package controller

import (
	"fmt"
    "github.com/astaxie/beego"
	//"github.com/astaxie/beego/logs"
	"../../util"
)

type GetTest001 struct {
    beego.Controller
}

func (p *GetTest001) Get() {
	d0 := "2013-03-12 00:00:00"
	d1 := util.GetCurMonthFirstDay(d0)
	d2 := util.GetCurMonthLastDay(d0)
	d3 := util.GetNextMonthFirstDay(d0)
	d4 := util.GetNextMonthLastDay(d0)
	fmt.Println("d0:", d0)
	fmt.Println("GetCurMonthFirstDay:", d1)
	fmt.Println("GetCurMonthLastDay:", d2)
	fmt.Println("GetNextMonthFirstDay:", d3)
	fmt.Println("GetNextMonthLastDay:", d4)

	d0 = "2014-02-12 00:00:00"
	d1 = util.GetCurMonthFirstDay(d0)
	d2 = util.GetCurMonthLastDay(d0)
	d3 = util.GetNextMonthFirstDay(d0)
	d4 = util.GetNextMonthLastDay(d0)
	fmt.Println("d0:", d0)
	fmt.Println("GetCurMonthFirstDay:", d1)
	fmt.Println("GetCurMonthLastDay:", d2)
	fmt.Println("GetNextMonthFirstDay:", d3)
	fmt.Println("GetNextMonthLastDay:", d4)

	d0 = "2012-01-12 00:00:00"
	d1 = util.GetCurMonthFirstDay(d0)
	d2 = util.GetCurMonthLastDay(d0)
	d3 = util.GetNextMonthFirstDay(d0)
	d4 = util.GetNextMonthLastDay(d0)
	fmt.Println("d0:", d0)
	fmt.Println("GetCurMonthFirstDay:", d1)
	fmt.Println("GetCurMonthLastDay:", d2)
	fmt.Println("GetNextMonthFirstDay:", d3)
	fmt.Println("GetNextMonthLastDay:", d4)

	d0 = "2014-01-12 00:00:00"
	d1 = util.GetCurMonthFirstDay(d0)
	d2 = util.GetCurMonthLastDay(d0)
	d3 = util.GetNextMonthFirstDay(d0)
	d4 = util.GetNextMonthLastDay(d0)
	fmt.Println("d0:", d0)
	fmt.Println("GetCurMonthFirstDay:", d1)
	fmt.Println("GetCurMonthLastDay:", d2)
	fmt.Println("GetNextMonthFirstDay:", d3)
	fmt.Println("GetNextMonthLastDay:", d4)

	d0 = "2014-12-12 00:00:00"
	d1 = util.GetCurMonthFirstDay(d0)
	d2 = util.GetCurMonthLastDay(d0)
	d3 = util.GetNextMonthFirstDay(d0)
	d4 = util.GetNextMonthLastDay(d0)
	fmt.Println("d0:", d0)
	fmt.Println("GetCurMonthFirstDay:", d1)
	fmt.Println("GetCurMonthLastDay:", d2)
	fmt.Println("GetNextMonthFirstDay:", d3)
	fmt.Println("GetNextMonthLastDay:", d4)

	p.TplName = "index.html"
}