package main

const HANDLE_DIG = " /dig?"
const HANDLE_MOVIE = "/movie/"
const HANDLE_LIST = "/list/"
const HANDLE_HTML = ".html"

type cmdParams struct {
	logFilePath string
	routineNum  int
}

type digData struct {
	time  string
	url   string
	refer string
	ua    string
}

type urlNode struct {
	unType string // 详情页 或者 列表页 或者 首页
	unRid  int    // Resource ID 资源ID
	unUrl  string // 当前这个页面的url
	unTime string // 当前访问这个页面的时间
}

type urlData struct {
	uid   string
	data  digData
	unode urlNode
}

type storageBlock struct {
	counterType  string
	storageModel string
	unode        urlNode
}
