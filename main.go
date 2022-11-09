package main

import (
	_ "github.com/go-sql-driver/mysql"
	classes2 "mygin/src/classes"
	"mygin/src/goft"
	"mygin/src/middlewares"
)

func main() {
	goft.Ignite().
		Beans(goft.NewGormAdapter()).
		Attach(middlewares.NewUserMiddleware()).
		Mount("v1", classes2.NewIndex(), classes2.NewUserClass()).
		Task("0/3 * * * * *", goft.Expr(".IndexClass.Test")).
		Launch()
}
