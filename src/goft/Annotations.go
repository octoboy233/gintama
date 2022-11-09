package goft

import (
	"fmt"
	"reflect"
	"strings"
)

//注解处理

type Annotations interface {
	SetTag(tag reflect.StructTag)
}

func init() { //依赖包时调用
	AnnotationList = []Annotations{new(Value)}
}

var AnnotationList []Annotations

func IsAnnotation(p reflect.Type) bool {
	for _, annotation := range AnnotationList {
		if reflect.TypeOf(annotation) == p {
			return true
		}
	}
	return false
}

type Value struct {
	tag         reflect.StructTag
	Beanfactory *BeanFactory
}

func (this *Value) SetTag(tag reflect.StructTag) {
	this.tag = tag
}

//这里从配置文件里取到对应内容
func (this *Value) String() string {
	//return "21"
	get_prefix := this.tag.Get("prefix")
	if get_prefix == "" {
		return ""
	}
	prefix := strings.Split(get_prefix, ".")
	if config := this.Beanfactory.GetBean(new(SysConfig)); config != nil {
		get_value := GetConfigValue(config.(*SysConfig).Config, prefix, 0)
		if get_value != nil {
			return fmt.Sprintf("%v", get_value)
		} else {
			return ""
		}
	} else {
		return ""
	}
}
