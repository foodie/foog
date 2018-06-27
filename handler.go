package foog

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

//空对象，可以替代任意的类型
type IObject interface {
}

type handlerEntity struct {
	object   IObject       //空对象
	method   reflect.Value //方法的值
	argType  reflect.Type  //方法接收参数的类型
	argIsRaw bool          //是否是原生的
}

//handler管理器，可以注册多个handler
type handlerManager struct {
	handlers map[string]*handlerEntity
}

var (
	//空类型
	//Elem返回v持有的接口保管的值的Value封装
	//或者v持有的指针指向的值的Value封装
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()

	//字节类型
	typeOfBytes = reflect.TypeOf(([]byte)(nil))

	typeOfString = reflect.TypeOf((*string)(nil))
)

//有三个类型
func isHandlerMethod(method reflect.Method) bool {
	//获取类型
	mt := method.Type
	//返回func类型的参数个数，如果不是函数，将会panic
	//返回func类型的参数个数
	log.Println("mt.NumIn()=", mt.NumIn(), method.Name, method.Type)
	if mt.NumIn() != 3 {
		return false
	}

	return true
}

//注册一个对象和方法
func (this *handlerManager) register(obj IObject) {
	//获取type
	t := reflect.TypeOf(obj)
	//获取value
	v := reflect.ValueOf(obj)

	//得到type的名字
	//返回持有v持有的指针指向的值的Value。
	//如果v持有nil指针，会返回Value零值；如果v不持有指针，会返回v。
	//返回类型的名称
	name := reflect.Indirect(v).Type().Name()

	//如果handlers为空，创建对象
	if this.handlers == nil {
		this.handlers = make(map[string]*handlerEntity)
	}

	//定义类型的方法的数量
	for m := 0; m < t.NumMethod(); m++ {
		//获取方法的Value
		method := t.Method(m)
		//类型，和名字
		mt := method.Type
		mn := method.Name
		log.Println("method-name=", mn)
		if isHandlerMethod(method) {
			//定义raw是否为false
			raw := false
			//返回func类型的第2个参数的类型
			//如果是字节类型，返回是原生的
			if mt.In(2) == typeOfBytes {
				raw = true
			}
			//对象名.方法名=》返回对象，
			keyName := strings.ToLower(fmt.Sprintf("%s.%s", name, mn))
			log.Printf("keyName=%s", keyName)
			log.Println(mt.In(0), mt.In(1), mt.In(2))
			this.handlers[keyName] = &handlerEntity{
				object:   obj,
				method:   v.Method(m), //第m个方法的值
				argType:  mt.In(2),    //第二个参数的类型，返回func类型的第i个参数的类型
				argIsRaw: raw,         //是否是原生的
			}
		} else { //不是isHandlerMethod，返回错误
			log.Printf("%s.%s register failed, argc=%d\n", name, mn, mt.NumIn())
		}
	}
}

//分发数据
func (this *handlerManager) dispatch(name string,
	sess *Session,
	data interface{}) {
	//获取对应的处理器
	h, ok := this.handlers[strings.ToLower(name)]
	if !ok {
		log.Println("not found handle by", name)
		return
	}

	//定义后续的处理函数
	defer func() {
		//错误恢复
		if err := recover(); err != nil {
			log.Println("dispatch error", name, err)
		}
	}()
	log.Printf("handle data=%v", data)

	//是否序列化后
	var serialized bool
	var argv reflect.Value

	//不是原生的，序列化不为空
	if !h.argIsRaw && sess.serializer != nil {
		log.Println("parsering")
		//获取byte类型
		if bytes, ok := data.([]byte); ok {
			//新建一个value存储值
			argv = reflect.New(h.argType.Elem())
			//接下bytes，把数据存储到argv里面
			err := sess.serializer.Decode(bytes, argv.Interface())
			if err != nil {
				log.Println("deserialize error", err.Error())
				return
			}
			//可序列化
			serialized = true
		}
	}
	//不可序列化，直接返回值
	if !serialized {
		//不可序列化，直接返回valueof的data
		argv = reflect.ValueOf(data)
	}
	log.Println("argv=%v", argv)
	//参数返回reflect的值
	args := []reflect.Value{
		reflect.ValueOf(sess), //参数1 是sess
		argv} //参数2是传递的参数
	//方法调用参数
	h.method.Call(args)
}
