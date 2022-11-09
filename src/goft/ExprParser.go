package goft

import (
	"bytes"
	"fmt"
	"html/template"
)

type Expr string //表达式类型

//执行表达式
//这个方法 用于定时任务中调用控制器中的方法（因为可能会使用到控制器中的变量）
//expr：".UserClass.Test" data在脚手架主类中的exprData切片
func ExecExpr(expr Expr, data map[string]interface{}) (string, error) {
	tpl := template.New("expr")
	//这里很关键 .UserClass.Test UserClass作为data中的key，取到控制器（一个结构体）,Test既可以取成员变量 也可以调用指针函数
	//利用了template库 比较取巧的一个实现
	t, err := tpl.Parse(fmt.Sprintf("{{%s}}", expr))
	if err != nil {
		return "", err
	}
	//以上是模版的声明（语法，函数也可） 下面是套用，从原始数据 ——模版——> 套用后的数据
	// data 是 形如 key "User" value控制器实体
	var buf = &bytes.Buffer{}
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
