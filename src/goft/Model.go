package goft

import (
	"encoding/json"
	"log"
)

type Model interface {
	//可以规定模型要实现的方法 比如 string
	String() string
}

type Models string

func MakeModels(v interface{}) Models {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
	}
	return Models(b)
}

type ModelList []Model
