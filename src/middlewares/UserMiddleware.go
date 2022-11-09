package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type UserMiddleware struct {
}

func NewUserMiddleware() *UserMiddleware {
	return &UserMiddleware{}
}

func (this *UserMiddleware) OnRequest(ctx *gin.Context) error {
	fmt.Println("my middleware")
	fmt.Println(ctx.DefaultQuery("name", "nobody"))
	return nil
}
