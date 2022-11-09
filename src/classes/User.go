package classes

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"mygin/src/goft"
	"mygin/src/models"
	"reflect"
	"time"
)

type UserClass struct {
	//删除掉了 gin.Engine
	//控制器不应该和服务类有相关联
	//这样太耦合了

	*goft.GormAdapter

	//类似java 的@Value.....注解 可以自动从application.yaml对应的条项中自动加载
	Age *goft.Value `prefix:"user.age"`
}

func NewUserClass() *UserClass {
	return &UserClass{}
}

func (this UserClass) Name() string {
	return "UserClass"
}

func (this *UserClass) Test(ctx *gin.Context) string {
	return "测试方法" + this.Age.String()
}

//两种方式 去返回一个实体类切片
//这一种 比较投机取巧 handler返回的是一个字符串 但是不同于直接返回字符串的函数
//我们拿到对象后marshall成字符串，然后定一个新的类型（其实是string），用来在responder中区别于返回string的
//然后对应类型Models的Responder中需要把返回的头设定成返回json格式的
//这样子就可以write一个string过去，但是前端接受到的是json格式的实体类切片数据
func (this *UserClass) GetUserList(ctx *gin.Context) goft.Models {
	users := []*models.UserModel{
		{UserName: "zyg", UserId: 101},
		{UserName: "lisi", UserId: 102},
	}
	return goft.MakeModels(users)
}

//比较常规的做法
//将获取到的对象切片遍历一遍返回model切片
//两种 responder
func (this *UserClass) GetUserList2(ctx *gin.Context) goft.ModelList {
	users := []*models.UserModel{
		{UserName: "boy", UserId: 103},
		{UserName: "wangwu", UserId: 104},
	}
	var list goft.ModelList
	for _, v := range users {
		//强制转换 和interface断言 都可以 个人喜欢用断言
		list = append(list, reflect.ValueOf(v).Interface().(goft.Model))
		//list = append(list, (goft.Model(v)))
	}
	return list
}

//采用依赖注入的方法去使用orm
func (this *UserClass) GetUserDetail(ctx *gin.Context) goft.Model {
	id := ctx.Query("id")
	user := &models.UserModel{}
	this.Table("user").Where("id=?", id).Find(user)
	return user
}

//统一处理error的演示
func (this *UserClass) GetUserInfo(ctx *gin.Context) goft.Model {
	user := models.NewUserModel()
	err := ctx.ShouldBindUri(user)
	//binduri 返回的是mustbinduri的值，如果失败会直接ctx.abortWithError 并且中间件无法重写
	//Headers were already written. Wanted to override status code 400 with 200
	//所以一般用shouldbind
	goft.Error(err, "用户id不符合要求")
	this.Table("user").Where("id=?", user.UserId).Find(user)

	//在这里执行异步任务，先把结果返回，再处理
	goft.Task(this.UpdateViews, this.TaskDone, user.UserId) //如果回调函数需要传参，用匿名函数把它包起来就可以实现
	return user
}

func (this *UserClass) UpdateViews(params ...interface{}) {
	time.Sleep(3000 * time.Millisecond)
	this.Table("user").Where("id=?", params[0]).
		Update("views", gorm.Expr("views+1"))
}

func (this *UserClass) TaskDone() {
	log.Println("user信息获取函数回调执行成功")
}

func (this UserClass) Build(goft *goft.Goft) { //这个穿参是关键
	goft.Handle("GET", "/test", this.Test)
	goft.Handle("GET", "/user/:id", this.GetUserInfo)
	goft.Handle("GET", "/user", this.GetUserDetail)
	goft.Handle("GET", "/userlist", this.GetUserList)
	goft.Handle("GET", "/userlist2", this.GetUserList2)
}
