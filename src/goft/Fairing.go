package goft

import "github.com/gin-gonic/gin"

//中间件接口，用来规范接口
type Fairing interface {
	OnRequest(ctx *gin.Context) error
}
