package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
)
var (
	db orm.Ormer
)
type Pic struct {
	Id int64
	Title string
	Pic string
	Date string
	CreatedAt string
}
func init(){
	orm.Debug = true // 是否开启调试模式 调试模式下会打印出sql语句
	//注册mysql
	orm.RegisterDataBase("default", "mysql", "root:wqg951122@tcp(127.0.0.1:3306)/girl_pic?charset=utf8", 30)
	orm.RegisterModelWithPrefix("gp_", new(Pic))
	db = orm.NewOrm()
}
//匹配标题
func GetTitle(htmls string)string{
	if htmls == "" {
		return ""
	}
	reg := regexp.MustCompile(`<h2>\s*<a\s*href="http://www.meizitu.com/a/\d*.html">(.*)</a>\s*</h2>`)
	result := reg.FindAllStringSubmatch(htmls,-1)
	if len(result) == 0{
		return ""
	}
	return result[0][1]
}
//匹配日期
func GetDate(htmls string)string{
	if htmls == "" {
		return ""
	}
	dayReg := regexp.MustCompile(`<div\s*class="day">(.*)</div>`)
	yearMonthReg := regexp.MustCompile(`<div\s*class="month_Year">(\d*)&nbsp;(\d*)</div>`)
	days := dayReg.FindAllStringSubmatch(htmls,-1)
	yearMonths := yearMonthReg.FindAllStringSubmatch(htmls,-1)
	if len(days) == 0 || len(yearMonths) == 0{
		return ""
	}
	return yearMonths[0][2]+"-"+yearMonths[0][1]+"-"+days[0][1]
}

//匹配图片列表
func GetPicture(htmls string)string{
	if htmls == ""{
		return  ""
	}
	boxReg := regexp.MustCompile(`<div\s*id="maincontent">([\s\S]*)<div\s*id="sidebar">`)
	box := boxReg.FindAllStringSubmatch(htmls,-1)
	picReg := regexp.MustCompile(`<img.*?src="(.*?)".*?/>`)
	pic := picReg.FindAllStringSubmatch(box[0][1],-1)
	if len(box) == 0 || len(pic) == 0{
		return ""
	}
	var res string
	for k,v := range pic{
		if k == 0{
			res += v[1]
		}else{
			res += ","+v[1]
		}
	}
	return res
}

//匹配其他图片链接
func GetUrls(htmls string)[]string{
	if htmls == "" {
		return []string{}
	}
	reg := regexp.MustCompile(`<a\s*href="(http://www.meizitu.com/a/\d*.html)"`)
	result := reg.FindAllStringSubmatch(htmls,-1)
	if len(result) == 0{
		return []string{}
	}
	var res []string
	for _,v := range result{
		res = append(res,v[1])
	}
	return res
}
//添加到数据库
func AddPic(pic *Pic)(int64,error){
	pic.Id = 0
	id,err := db.Insert(pic)
	return id,err
}

