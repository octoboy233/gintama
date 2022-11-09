package goft

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
)

//总控制器
type Goft struct {
	//看起来有点像继承 实际上是嵌套（go只有嵌套）
	*gin.Engine                  //把engine放到主类中
	g           *gin.RouterGroup //可以做到goft.Handle的书写效果
	beanFactory *BeanFactory
	exprData    map[string]interface{} //存放了key：控制器名称 value：控制器实体
}

//构造函数 相当于newGoft
func Ignite() *Goft {
	goft := &Goft{Engine: gin.New(), beanFactory: NewBeanFactory(), exprData: map[string]interface{}{}}
	goft.beanFactory.setBeans(InitConfig()) //整个配置加载到bean中
	goft.Use(ErrorHandler())                //强制使用错误handler中间件
	goft.FuncMap = map[string]any{
		"Strong": func(txt string) template.HTML {
			return template.HTML("<strong>" + txt + "</strong>")
		},
	} //gin中封装的模版函数用法,在html文件中用管道引用即可
	goft.LoadHTMLGlob("views/*") //加载全局html路径
	return goft
}

//
func (this *Goft) Beans(bean Bean) *Goft {
	//取出bean的名称加入到exprdata中
	//this.exprData[bean.Name()] = bean
	this.beanFactory.setBeans(bean)
	return this
}

//启动
func (this *Goft) Launch() {
	var port = 8080
	if config := this.beanFactory.GetBean(new(SysConfig)); config != nil {
		port = config.(*SysConfig).Server.Port
	}
	getCronTask().Start()
	this.Run(fmt.Sprintf(":%d", port))
}

//重载engine的handle，用于在挂载时使用goft控制器中的group
func (this *Goft) Handle(httpMethod, relativePath string, h interface{}) *Goft {
	//if handler, ok := h.(func(ctx *gin.Context) string); ok {
	//	this.g.Handle(httpMethod, relativePath, func(context *gin.Context) {
	//		context.String(200, handler(context))
	//	})
	//}
	if handler := Convert(h); handler != nil {
		this.g.Handle(httpMethod, relativePath, handler)
	}
	return this
}

//分组
//路由
func (this *Goft) Mount(group string, classes ...IClass) *Goft { //链式调用
	this.g = this.Group(group)
	for _, class := range classes {
		//class_ref := reflect.ValueOf(class).Elem()
		//if class_ref.NumField() > 0 && this.dba != nil {
		//	// nil -> (*gormAdapter)(nil) -> ...
		//	// 空值 到 空指针 到 改变指针指向
		//	//class_ref.Field(0).Set(reflect.ValueOf(this.dba))
		//	class_ref.Field(0).Set(reflect.New(class_ref.Field(0).Type().Elem()))
		//	//错误写法：reflect.New  必须要得到指针 指向的对象；不能初始化一个指针
		//	// class_ref.Field(0).Set(reflect.New(class_ref.Field(0).Type()).Elem())
		//	class_ref.Field(0).Elem().Set(reflect.ValueOf(this.dba).Elem())
		//}
		//这里替代了上面一部分代码，这样可以传入多个需要注入的对象，上面是控制器只注入一个依赖的演示
		this.beanFactory.inject(class)
		////上一步实现了将控制器中的成员变量进行依赖注入 这一步实现将控制器本身将入到bean容器中 （spring也是这么做的
		this.Beans(class) //beans方法里会set exprData，在执行定时任务的时候会把里面的控制器实体传过去，暂时没发现注入控制器有啥用。。。
		class.Build(this)
	}
	return this
}

//中间件
func (this *Goft) Attach(fair Fairing) *Goft {
	this.Use(func(context *gin.Context) {
		err := fair.OnRequest(context)
		if err != nil {
			context.AbortWithStatusJSON(400, gin.H{
				"error": err.Error(),
			})
		} else {
			context.Next()
		}
	})
	return this
}

//增加定时任务 expr:0/3 * * * * *
func (this *Goft) Task(cron string, expr interface{}) *Goft {
	var err error
	if f, ok := expr.(func()); ok {
		_, err = getCronTask().AddFunc(cron, f)
	} else if exp, ok := expr.(Expr); ok {
		_, err = getCronTask().AddFunc(cron, func() {
			_, err = ExecExpr(exp, this.exprData)
		})
	}
	if err != nil {
		log.Println(err)
	}
	return this
}
