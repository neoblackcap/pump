package pump

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"unsafe"
)

type Container struct {
	dependencies map[string]interface{}
}

func (c Container) FindDependency(name string) (interface{}, error) {
	d := c.dependencies[name]
	if d != nil {
		return d, nil
	} else {
		return nil, errors.New("found no dependency")
	}
}

func NewContainer() *Container {
	m := make(map[string]interface{})
	return &Container{m}
}

func (c Container) tearDown(s string) {
	fmt.Fprint(os.Stderr, s)
}

func (c Container) wireMap(field *reflect.Value, data interface{}) {
	//m := reflect.MakeMap(field.Type())
	v := reflect.ValueOf(data)
	//for _, key := range v.MapKeys() {
	//	m.SetMapIndex(key, v.MapIndex(key))
	//}
	field.Set(v)
}

func (c Container) wireArray(field *reflect.Value, data interface{}) {
	v := reflect.ValueOf(data)
	field.Set(v)
}

func (c Container) wireChan(field *reflect.Value, data interface{}) {
	v := reflect.ValueOf(data)
	field.Set(v)
}

func (c Container) wireInterface(field *reflect.Value, data interface{}) {
	v := reflect.ValueOf(data)
	field.Set(v)
}

func (c Container) wireStruct(field *reflect.Value, data interface{}) {
	v := reflect.ValueOf(data)
	field.Set(v)
}

func (c Container) wireUnsafePointer(field *reflect.Value, data interface{}) {
	v := reflect.ValueOf(data)
	field.Set(v)
}

func (c Container) wire(src interface{}) {
	v := reflect.ValueOf(src).Elem().Elem()
	if v.CanSet() {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := t.Field(i)
			tag := fieldType.Tag.Get("inject")
			instance, err := c.FindDependency(tag)
			if err != nil {
				panic(err)
			}
			if !field.CanSet() {
				field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
			}

			switch reflect.TypeOf(instance).Kind() {
			case reflect.Bool:
				field.SetBool(instance.(bool))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				field.SetInt(instance.(int64))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				field.SetUint(instance.(uint64))
			case reflect.Float32, reflect.Float64:
				field.SetFloat(instance.(float64))
			case reflect.Complex64, reflect.Complex128:
				field.SetComplex(instance.(complex128))
			case reflect.Ptr:
				field.SetPointer(unsafe.Pointer(&instance))
			case reflect.String:
				field.SetString(instance.(string))
			case reflect.Map:
				c.wireMap(&field, instance)
			case reflect.Slice, reflect.Array:
				c.wireArray(&field, instance)
			case reflect.Chan:
				c.wireChan(&field, instance)
			case reflect.Interface:
				c.wireInterface(&field, instance)
			case reflect.Struct:
				c.wireStruct(&field, instance)
			case reflect.UnsafePointer:
				c.wireUnsafePointer(&field, instance)
			}
		}
	}
}

func (c Container) Wire(src interface{}) {
	c.wire(src)
}

func (c *Container) Register(name string, provider interface{}) {
	c.dependencies[name] = provider
}
