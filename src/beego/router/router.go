package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
    Index "../controller"
)

func init(){
	logs.Debug("router init...")

    //index交给 这个方法处理
    //beego默认需要模板文件
    beego.Router("/index",&Index.IndexController{},"*:Index")//"get:"
	beego.Router("/fundinfo",&Index.GetFundInfo{},"*:Get")//"get:"
	beego.Router("/update.html",&Index.Update{},"*:Get")//"get:"
	
	beego.Router("/income",&Index.GetFundIncome{},"*:Get")//"get:"



	beego.Router("/test001",&Index.GetTest001{},"*:Get")//"get:"


	
}

func test() {

}