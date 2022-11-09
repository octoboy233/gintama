package main

import (
	"bytes"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
)

type Me struct {
	Age int
}

func (this *Me) GetAge() int {
	return 19
}

func main() {
	tmp := template.New("test").Funcs(map[string]any{})

	t, err := tmp.Parse("{{ .me.Age }}")
	if err != nil {
		log.Println(err)
	}
	//以上都是模版的申明（语法、函数） 下面是套用，从原始数据 ——模版——> 套用后的数据
	buf := bytes.Buffer{}
	err = t.Execute(&buf, map[string]any{
		"me": &Me{20},
	})
	if err != nil {
		log.Println(err)
	}
	fmt.Println(buf.String())
}
