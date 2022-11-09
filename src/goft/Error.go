package goft

import (
	"github.com/gin-gonic/gin"
)

//error中间件
func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		//使用recover（）捕捉panic异常的时候，需要defer来读取一个匿名函数，
		defer func() {
			//不使用recover，会使panic时整个程序挂掉
			//defer的作用是当程序遇到panic的时候，系统将跳过后面的代码，进入defer，
			//使用recover()则可以返回捕获到的panic的值
			if e := recover(); e != any(nil) {
				context.JSON(200, gin.H{"error": e})
			}
		}()
		context.Next()
	}
}

func Error(err error, msg ...string) {
	if err == nil {
		return
	} else {
		errMsg := err.Error()
		if len(msg) > 0 {
			errMsg = msg[0]
		}
		panic(any(errMsg))
	}
}
