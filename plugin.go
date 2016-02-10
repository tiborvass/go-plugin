// Package plugin allows you to easily define a plugin for your Go application
// and have it call out at runtime, to C shared libraries fully or partially
// implementing the user-defined plugin.
//
// The advantage of this is that the implementation of the plugin is language-agnostic.
//
// Tested only on 64bit Linux.
package plugin

import (
	"errors"
	"reflect"

	"github.com/tiborvass/dl"
)

type Plugin struct {
	dl *dl.DL
}

var _plugin = reflect.TypeOf(Plugin{}).Name()

func (p Plugin) Close() error {
	if p.dl != nil {
		return p.dl.Close()
	}
	return nil
}

var nopFn = func([]reflect.Value) []reflect.Value { return nil }

func Open(plugin interface{}, path string) error {
	v := reflect.ValueOf(plugin)
	t := v.Type()
	if t.Kind() != reflect.Ptr {
		return errors.New("plugin needs to be a pointer to a struct")
	}
	v = v.Elem()
	t = v.Type()
	if t.Kind() != reflect.Struct {
		return errors.New("Open expects a plugin of type struct{ interface{ Method() } }")
	}
	lib, err := dl.Open(path, 0)
	if err != nil {
		return err
	}
	for i := 0; i < v.NumField(); i++ {
		tf := t.Field(i)
		if tf.Name != _plugin {
			sym := v.Field(i).Interface()
			if err := lib.Sym(tf.Name, &sym); err != nil && tf.Type.Kind() == reflect.Func {
				fn := reflect.MakeFunc(tf.Type, nopFn)
				v.Field(i).Set(fn)
			} else {
				v.Field(i).Set(reflect.ValueOf(sym))
			}
		} else {
			p := Plugin{lib}
			v.Field(i).Set(reflect.ValueOf(p))
		}
	}
	return nil
}
