package goft

import (
	"github.com/robfig/cron/v3"
	"sync"
)

//task异步任务组件

//量小的异步任务就用channel 别用协程池 多的就用mq
var taskList chan *TaskExecutor //任务列表
var once sync.Once
var onceCron sync.Once
var taskCron *cron.Cron //定时任务

//初始化函数 消费队列中的任务
func init() {
	chlist := getTaskList() //被引用时调用 阻塞在这里 等待任务放进来
	go func() {
		for t := range chlist {
			doTask(t)
		}
	}()
}

//执行单个任务
func doTask(executor *TaskExecutor) {
	go func() { //这里为什么要用协程 因为需要defer在每个task结束后取执行回调
		defer executor.callBack()
		executor.Exec()
	}()
}

//单例创建定时任务对象
func getCronTask() *cron.Cron {
	onceCron.Do(func() {
		taskCron = cron.New(cron.WithSeconds())
	})
	return taskCron
}

type TaskFunc func(params ...interface{})

func getTaskList() chan *TaskExecutor {
	once.Do(func() { //单例模式
		taskList = make(chan *TaskExecutor) //初始化
	})
	return taskList
}

type TaskExecutor struct {
	f        TaskFunc
	p        []interface{}
	callBack func() //回调函数
}

func NewTaskExecutor(f TaskFunc, cb func(), p []interface{}) *TaskExecutor {
	return &TaskExecutor{f: f, p: p, callBack: cb}
}

func (this *TaskExecutor) Exec() {
	this.f(this.p...)
}

func Task(f func(params ...interface{}), cb func(), params ...interface{}) {
	go func() { //这里用了协程 因为channel没有缓冲
		getTaskList() <- NewTaskExecutor(f, cb, params) //注意这里的写法
	}()
}
