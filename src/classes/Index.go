package classes

import (
	"github.com/gin-gonic/gin"
	"log"
	"mygin/src/goft"
)

type IndexClass struct {
}

func NewIndex() *IndexClass {
	return &IndexClass{}
}

func (this IndexClass) Name() string {
	return "IndexClass"
}

func (this IndexClass) Build(goft *goft.Goft) {
	goft.Handle("GET", "/index", this.Index)
}

func (this *IndexClass) Index(ctx *gin.Context) goft.View {
	ctx.Set("name", "首页")
	return "index"
}

func (this *IndexClass) Test() interface{} {
	log.Println("测试定时方法")
	return nil
}
