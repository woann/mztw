package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/axgle/mahonia"
	"mztw/models"
	"time"
)

type GrabdataDisat struct {
	beego.Controller
}

func (this *GrabdataDisat) Getgrabdata(){
	id := this.Input().Get("id")
	if id == ""{
		this.Ctx.WriteString("id is miss")
		return
	}
	var pic models.Pic
	//连接redis
	models.ConnectRedis("127.0.0.1:6379")
	//关闭模板渲染
	this.EnableRender = false
	//爬虫入口url
	url1 := "http://www.meizitu.com/a/more_1.html"
	models.PutinQueue(url1)
	for {
		//爬取
		//获取队列长度
		len := models.GetQueueLength()
		if len == 0{
			fmt.Println("队列已清空")
			break
		}
		//从队列获取url
		url1 = models.PopfromQueue()
		//url1 = "http://www.meizitu.com/a/496.html"
		//判断该连接是否访问过
		if models.IsVisit(url1){
			fmt.Println("该url访问过")
			continue
		}
		fmt.Println("开始获取数据"+url1)
		resps := httplib.Get(url1)
		htmls,err := resps.String()
		if err != nil{
			fmt.Println("获取数据超时")
			continue
		}
		fmt.Println("成功获取数据")
		htmls = ConvertToString(htmls, "gbk", "utf-8")
		pic.Title = models.GetTitle(htmls)
		pic.Date = models.GetDate(htmls)
		pic.Pic = models.GetPicture(htmls)
		t := time.Now()
		pic.CreatedAt = t.Format("2006-01-02 15:04:05")
		if pic.Pic != ""{
			_,err1 := models.AddPic(&pic)
			if err1 != nil{
				panic(err1)
			}
		}
		//获取该页面所有图片文章连接
		urlList := models.GetUrls(htmls)
		for _,url := range urlList{
			if models.IsVisit(url){
				fmt.Println("已经访问过"+url)
				continue
			}
			//将url添加到队列
			models.PutinQueue(url)
			fmt.Println("添加到队列"+url)
		}
		//将url标记为已访问
		models.AddToSet(url1)
		time.Sleep(time.Second)
	}
	//爬取结束
	this.Ctx.WriteString("no more")

}

//gbk转utf8
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

