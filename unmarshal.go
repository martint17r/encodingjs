package encodingjs

import (
	"fmt"
	"reflect"

	"github.com/robertkrimen/otto"
)

const maxInt = int64(^uint(0) >> 1)

// UnmarshalJSVar transfers the contents of d into result
func Unmarshal(d otto.Value, result interface{}) error {
	rv := reflect.ValueOf(result)
	return unmarshal(d, rv)
}

func unmarshal(d otto.Value, rv reflect.Value) error {
	//fmt.Println(d, "<-", rv, rv.Type().Name(), rv.CanSet())
	switch rv.Kind() {
	case reflect.Ptr:
		return unmarshal(d, rv.Elem())
	case reflect.Struct:
		return unmarshalStruct(d, rv)
	case reflect.Slice:
		return unmarshalSlice(d, rv)
	case reflect.Map:
		return unmarshalMap(d, rv)
	case reflect.String:
		rv.SetString(d.String())
	case reflect.Int:
		if !d.IsNumber() {
			return &InvalidValueError{rv.Type().String(), d}
		}
		iv, err := d.ToInteger()
		if err != nil {
			return err
		}
		rv.SetInt(iv)
	case reflect.Float64:
		if !d.IsNumber() {
			return &InvalidValueError{rv.Type().String(), d}
		}
		fv, err := d.ToFloat()
		if err != nil {
			return err
		}
		rv.SetFloat(fv)
	case reflect.Float32:
		if !d.IsNumber() {
			return &InvalidValueError{rv.Type().String(), d}
		}
		fv, err := d.ToFloat()
		if err != nil {
			return err
		}
		rv.SetFloat(fv)
	default:
		return &UnsupportedTypeError{rtype: rv.Type()}
	}
	return nil
}

func unmarshalStruct(d otto.Value, rv reflect.Value) error {
	rt := rv.Type()
	if !d.IsObject() {
		return ErrObjectExpected
	}
	for i := 0; i < rt.NumField(); i++ {
		if dv, err := d.Object().Get(rt.Field(i).Name); err == nil {
			unmarshal(dv, rv.Field(i))
		} else if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalMap(d otto.Value, rv reflect.Value) error {
	rt := rv.Type()
	if !d.IsObject() {
		return ErrObjectExpected
	}
	dobj := d.Object()

	keyT := rt.Key()
	valueT := rt.Elem()
	rv.Set(reflect.MakeMap(rt))
	for _, key := range dobj.Keys() {
		kv := reflect.Indirect(reflect.New(keyT))
		kv.SetString(key)
		fv := reflect.Indirect(reflect.New(valueT))
		dval, err := dobj.Get(key)
		if err != nil {
			return err
		}
		err = unmarshal(dval, fv)
		if err != nil {
			return err
		}
		rv.SetMapIndex(kv, fv)
	}

	return nil
}

func unmarshalSlice(d otto.Value, rv reflect.Value) error {
	rt := rv.Type()
	if d.Class() != "Array" {
		return ErrArrayExpected
	}
	darr := d.Object()
	dlens, err := darr.Get("length")
	if err != nil {
		return fmt.Errorf("could not get array length: %s", err)
	}
	dlen64, err := dlens.ToInteger()
	if err != nil {
		return fmt.Errorf("non numeric array length: %s", err)
	}
	if dlen64 > maxInt {
		return fmt.Errorf("array length exceeded: %d", dlen64)
	}
	dlen := int(dlen64)
	rslice := reflect.MakeSlice(rt, dlen, dlen)
	for i := 0; i < dlen; i++ {
		delem, err := darr.Get(fmt.Sprintf("%d", i))
		if err != nil {
			return err
		}
		unmarshal(delem, rslice.Index(i))
	}
	rv.Set(rslice)
	return nil
}
