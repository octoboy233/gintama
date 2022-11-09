package goft

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

var ResponderList []Responder

func init() {
	ResponderList = []Responder{new(StringResponder),
		new(ModelResponder),
		new(ModelsResponder),
		new(ModelListResponder),
		new(ViewResponder),
	} //new返回的是指针类型
}

func Convert(handler interface{}) gin.HandlerFunc {
	h_ref := reflect.ValueOf(handler)
	for _, r := range ResponderList {
		r_ref := reflect.ValueOf(r).Elem()            //需要elem获取指针指向内容
		if h_ref.Type().ConvertibleTo(r_ref.Type()) { //两种type类型能否转换
			r_ref.Set(h_ref)
			return r_ref.Interface().(Responder).RespondTo()
			//return handler.(Responder).RespondTo() //这里很关键
			//这里错在，传进来的func(*gin.Context) string类型是没有实现RespondTo方法的 断言会panic
			//而r已经是StringResponder类型 肯定是实现了RespondTo方法的 可以断言
		}
	}
	return nil
}

type Responder interface {
	RespondTo() gin.HandlerFunc
}

type StringResponder func(*gin.Context) string

func (this StringResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(200, this(context))
	}
}

type View string

type ViewResponder func(*gin.Context) View

func (this ViewResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		//kv := context.Keys
		context.HTML(200, fmt.Sprintf("%s.html", this(context)), context.Keys) //不能在上面定义，取keys要在控制器方法之后
	}
}

type ModelResponder func(*gin.Context) Model

func (this ModelResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, this(context))
	}
}

type ModelsResponder func(*gin.Context) Models

func (this ModelsResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Content-type", "application/json")
		context.Writer.WriteString(string(this(context)))
	}
}

type ModelListResponder func(*gin.Context) ModelList

func (this ModelListResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, this(context))
	}
}
