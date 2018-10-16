package main

import (
	"flag"
	"github.com/mediocregopher/radix.v2/pool"
	"os"
	"time"
)

func parseParam() (cmdParams, *string) {
	// 获取参数
	logFilePath := flag.String("logFilePath", "/Users/pangee/Public/nginx/logs/dig.log", "log file path")
	routineNum := flag.Int("routineNum", 5, "consumer numble by goroutine")
	l := flag.String("l", "/tmp/log", "this programe runtime log target file path")
	flag.Parse()

	params := cmdParams{*logFilePath, *routineNum}
	return params, l
}

func writeLog(params cmdParams, l *string) {
	logFd, err := os.OpenFile(*l, os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		log.Out = logFd
		defer logFd.Close()
	}
	log.Infof("Exec start.")
	log.Infof("Params: logFilePath=%s, routineNum=%d", params.logFilePath, params.routineNum)
}

func initRedisPool(params cmdParams) *pool.Pool {
	// Redis Pool
	redisPool, err := pool.New("tcp", "localhost:6379", 2*params.routineNum)
	if err != nil {
		log.Fatalln("Redis pool created failed.")
		panic(err)
	} else {
		go func() {
			for {
				redisPool.Cmd("PING")
				time.Sleep(3 * time.Second)
			}
		}()
	}
	return redisPool
}

func main() {
	// 获取参数
	params, l := parseParam()

	// 打日志
	writeLog(params, l)

	// 初始化一些channel，用于数据传递
	var logChannel = make(chan string, 3*params.routineNum)
	var pvChannel = make(chan urlData, params.routineNum)
	var uvChannel = make(chan urlData, params.routineNum)
	var storageChannel = make(chan storageBlock, params.routineNum)

	// Redis Pool
	redisPool := initRedisPool(params)

	// 日志消费者
	go readFileLinebyLine(params, logChannel)

	// 创建一组日志处理
	for i := 0; i < params.routineNum; i++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}

	// 创建PV UV 统计器
	go pvCounter(pvChannel, storageChannel)
	go uvCounter(uvChannel, storageChannel, redisPool)
	// 可扩展的 xxxCounter

	// 创建 存储器
	go dataStorage(storageChannel, redisPool)

	time.Sleep(1000 * time.Second)
}
