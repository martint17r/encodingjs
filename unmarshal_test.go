package encodingjs_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/martint17r/encodingjs"
	"github.com/robertkrimen/otto"
)

func ExampleUnmarshal() {
	vm := otto.New()
	jsresult, err := vm.Run(`a={Foo: "bar", Fubar:"bob"}; a`)
	if err != nil {
		panic(err)
	}
	result := struct {
		Foo   string
		Fubar string
	}{}
	err = encodingjs.Unmarshal(jsresult, &result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", result)
	// Output: {Foo:bar Fubar:bob}
}

func TestUnmarshal(t *testing.T) {
	var validTests = []struct {
		script string
		want   interface{}
	}{
		{
			script: `12345; `,
			want:   int(12345),
		},
		{
			script: `54321.0; `,
			want:   float64(54321),
		},
		{
			script: `54321.0; `,
			want:   float32(54321),
		},
		{
			script: `"foo"`,
			want:   "foo",
		},
		{
			script: `[];`,
			want:   []interface{}{},
		},
		{
			script: `var a = {}; a['A'] = "fubar"; a; `,
			want:   struct{ A string }{A: "fubar"},
		},
		{
			script: `["foo", "bar"];`,
			want:   []string{"foo", "bar"},
		},
		{
			script: `var a = []; a[0]={}; a[0]['A'] = "fubar"; a.push({A:"bar"}); a; `,
			want:   []struct{ A string }{{A: "fubar"}, {A: "bar"}},
		},
		{
			script: `[{A:{B:"fubar"}}]`,
			want:   []struct{ A struct{ B string } }{{A: struct{ B string }{B: "fubar"}}},
		},
		{
			script: `a={ A:1, B:2, C:3 }; a`,
			want:   map[string]int{"A": 1, "B": 2, "C": 3},
		},
	}

	for _, test := range validTests {
		rt := reflect.TypeOf(test.want)
		v := reflect.New(rt).Interface()

		err := evaluateJS(test.script, v)
		if err != nil {
			t.Error(err)
			return
		}

		vv := reflect.ValueOf(v).Elem().Interface()
		if !reflect.DeepEqual(vv, test.want) {
			t.Errorf("\nwant: `%#v`,\n got: `%#v`\n", test.want, vv)
		}
	}
}

func TestInvalid(t *testing.T) {
	var errTests = []struct {
		script string
		arg    interface{}
		err    error
	}{
		{
			script: `{A:"b"};`,
			arg:    []string{},
			err:    encodingjs.ErrArrayExpected,
		},
		{
			script: `"b";`,
			arg:    struct{}{},
			err:    encodingjs.ErrObjectExpected,
		},
		{
			script: `"b";`,
			arg:    map[string]string{},
			err:    encodingjs.ErrObjectExpected,
		},
	}

	for _, test := range errTests {
		rt := reflect.TypeOf(test.arg)
		v := reflect.New(rt).Interface()
		err := evaluateJS(test.script, v)
		//fmt.Printf("v=%#v\n", reflect.ValueOf(v).Elem().Interface())
		if err != test.err {
			t.Errorf("%s\nwant: `%#v`,\n got: `%#v`\n", test.script, test.err, err)
		}
	}
}

func TestOttoError(t *testing.T) {
	var fubar interface{}
	err := evaluateJS("a", fubar)
	if err.Error() != "ReferenceError: 'a' is not defined" {
		t.Errorf("wanted ReferenceError, got: `%#v`\n", err)
	}
}

func TestUnknownType(t *testing.T) {
	var f func()
	err := evaluateJS("", f)
	if _, ok := err.(*encodingjs.UnsupportedTypeError); !ok {
		t.Errorf("wanted UnknownTypeError, got: `%#v`\n", err)
	}

	// to fix coverage for err.Error()
	_ = fmt.Sprintf("%s", err)
}

func TestInvalidTypes(t *testing.T) {
	var invalidTypeTests = []struct {
		script string
		arg    interface{}
		wanted string
	}{
		{
			script: `a={ A:"fubar" }; a`,
			arg:    map[string]int{},
			wanted: "int",
		},
		{
			script: `a={ A:"fubar" }; a`,
			arg:    map[string]float64{},
			wanted: "float64",
		},
		{
			script: `a={ A:"fubar" }; a`,
			arg:    map[string]float32{},
			wanted: "float32",
		},
	}

	for _, test := range invalidTypeTests {
		rt := reflect.TypeOf(test.arg)
		v := reflect.New(rt).Interface()
		err := evaluateJS(test.script, v)
		if ive, ok := err.(*encodingjs.InvalidValueError); !ok {
			t.Errorf("wanted InvalidValueError, got: `%#v`\n", err)
		} else if ive.Wanted != test.wanted {
			t.Errorf("wanted InvalidValueError %s, got: %s\n", test.wanted, ive.Wanted)
		}

		// to fix coverage for err.Error()
		_ = fmt.Sprintf("%s", err)
	}
}

func evaluateJS(script string, result interface{}) error {
	vm := otto.New()
	jsresult, err := vm.Run(script)
	if err != nil {
		return err
	}
	return encodingjs.Unmarshal(jsresult, result)
}
