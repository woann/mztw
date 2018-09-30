package models

import (
	"github.com/astaxie/goredis"
)

const (
	URL_QUEUE = "url_queue"
	URL_VISIT_SET = "url_visit_set"
	URL_QUEUE_SET = "url_queue_set"
)

//定义一个链接实例
var client goredis.Client

//连接redis
func ConnectRedis(addr string){
	client.Addr = addr
	client.Db = 0
}
//加入队列
func PutinQueue(url string){
	client.Lpush(URL_QUEUE,[]byte(url))
}
//从队列取
func PopfromQueue() string{
	res,err := client.Rpop(URL_QUEUE)
	if err != nil{
		panic(err)
	}
	return string(res)
}

func AddToSet(url string){
	//将访问过的url放到集合中
	client.Sadd(URL_VISIT_SET,[]byte(url))
}
//获取队列数量
func GetQueueLength() int{
	len,err := client.Llen(URL_QUEUE)
	if err != nil{
		return 0
	}
	return len
}

//判断是否已经访问过
func IsVisit(url string)bool{
	bIsVisit,err := client.Sismember(URL_VISIT_SET,[]byte(url))
	if err != nil{
		return false
	}
	return bIsVisit
}

//判断是否在队列里
func IsQueue(url string)bool{
	bIsVisit,err := client.Sismember(URL_QUEUE_SET,[]byte(url))
	if err != nil{
		return false
	}
	return bIsVisit
}

