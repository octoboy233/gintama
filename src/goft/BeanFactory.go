package goft

import (
	"reflect"
)

type Bean interface {
	Name() string
}

func (this *BeanFactory) Name() string { //beanfactory本身也是bean，所以要实现接口方法
	return "BeanFactory"
}

type BeanFactory struct {
	beans []Bean
}

func NewBeanFactory() *BeanFactory {
	bf := &BeanFactory{beans: make([]Bean, 0)}
	//因为在脚手架初始化时把配置文件加载了并装载进来 所以注解要是用配置有下面三部曲
	//1.将自己装载进去 也就是此步骤
	//2.注解类型中定义了beanfactory
	//3.在注入注解时有两个操作（1）把注解对象内的beanfactory注入，给注解对象注入一个指向nil的指针并将其tag获取到然后塞进去
	bf.beans = append(bf.beans, bf) //很关键的一步：将自己装载进去 这样在annotation里也能主动注入beanfactory，将通过它取配置
	return bf
}

//内部使用 在main中传入 装载bean
func (this *BeanFactory) setBeans(bean Bean) {
	this.beans = append(this.beans, bean)
}

//外部使用、
//传入一个new的实体对象即可（new返回的是一个指针） 返回装载好的对象
func (this *BeanFactory) GetBean(v interface{}) interface{} {
	return this.getBean(reflect.TypeOf(v))
}

//内部使用 传入的是从控制器中取到的需要注入的成员变量
//是否在装载列表中，是的话就返回装载的实体对象
func (this *BeanFactory) getBean(t reflect.Type) interface{} {
	for _, bean := range this.beans {
		if reflect.TypeOf(bean) == t {
			return bean
		}
	}
	return nil
}

//获得控制器中需要注入的成员变量 并实现注入
func (this *BeanFactory) inject(class IClass) {
	class_ref := reflect.ValueOf(class).Elem()
	for i := 0; i < class_ref.NumField(); i++ {
		member := class_ref.Field(i)
		if member.Kind() == reflect.Ptr && member.IsNil() {
			if prop := this.getBean(member.Type()); prop != nil {
				//这里可以替代下面两行
				//member.Set(reflect.ValueOf(prop))
				//reflect.New必须传入一个指针指向的对象而非一个值
				member.Set(reflect.New(member.Type().Elem()))   //这部是把member从nil转换成一个指向nil的指针
				member.Elem().Set(reflect.ValueOf(prop).Elem()) //这里是改变指针指向的值
			}
			//如果type在annotation列表里已经加入，则认定为annotation
			if IsAnnotation(member.Type()) {
				//这里是防止下一步调用setTag出现空指针
				member.Set(reflect.New(member.Type().Elem()))
				//这里注意，tag成员变量仅可以从结构体的type的field中取
				member.Interface().(Annotations).SetTag(class_ref.Type().Field(i).Tag)
				//把注解对象中的beanfactory自动注入
				this.Inject(member.Interface())
			}
		}
	}
}

//给外部用的 （后面还要改,这个方法不处理注解)
func (this *BeanFactory) Inject(object interface{}) {

	vObject := reflect.ValueOf(object)
	if vObject.Kind() == reflect.Ptr { //由于不是控制器 ，所以传过来的值 不一定是指针。因此要做判断
		vObject = vObject.Elem()
	}
	for i := 0; i < vObject.NumField(); i++ {
		f := vObject.Field(i)
		if f.Kind() != reflect.Ptr || !f.IsNil() {
			continue
		}
		if p := this.getBean(f.Type()); p != nil && f.CanInterface() { //判断 能够实现接口方法
			f.Set(reflect.New(f.Type().Elem()))
			f.Elem().Set(reflect.ValueOf(p).Elem())

		}

	}
}
