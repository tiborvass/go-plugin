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

// Plugin is a struct that must be embedded in a user-defined Plugin struct.
// It ensures proper closing of the shared library.
type Plugin struct {
	dl *dl.DL
}

// Close closes the shared library resource.
// Typically called in a defer statement after Open().
func (p Plugin) Close() error {
	if p.dl != nil {
		return p.dl.Close()
	}
	return nil
}

var _plugin = reflect.TypeOf(Plugin{}).Name()
var nopFn = func([]reflect.Value) []reflect.Value { return nil }

// Open retrieves the symbols defined in plugin, from the shared library at path.
// path should omit the file extension (e.g. "plugin" instead of "plugin.so").
// plugin should be a pointer to a struct embedding the Plugin struct.
func Open(plugin interface{}, path string) error {
	v := reflect.ValueOf(plugin)
	t := v.Type()
	if t.Kind() != reflect.Ptr {
		return errors.New("Open expects a plugin to be a pointer to a struct")
	}
	v = v.Elem()
	t = v.Type()
	if t.Kind() != reflect.Struct {
		return errors.New("Open expects a plugin to be a pointer to a struct")
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
