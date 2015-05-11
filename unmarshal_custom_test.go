package encodingjs_test

import (
	"net/url"
	"testing"

	"github.com/martint17r/encodingjs"
	"github.com/robertkrimen/otto"
)

type DataStructure struct {
	field1 string
	link   *url.URL
}

func (ds *DataStructure) UnmarshalJS(d otto.Value) error {
	tmp := struct {
		Field1 string
		Link   string `js:"Field2"`
	}{}
	encodingjs.Unmarshal(d, &tmp)

	ds.field1 = tmp.Field1
	var err error
	if tmp.Link != "" {
		ds.link, err = url.Parse(tmp.Link)
	}
	return err
}

func TestUnmarshaler(t *testing.T) {
	got := DataStructure{}
	err := evaluateJS(`a={Field1:"foo", Field2:"http://go/"}`, &got)
	if err != nil {
		t.Error(err)
	}
	wanted := DataStructure{"foo", nil}
	wanted.link, _ = url.Parse("http://go/")
	if got.link == nil || wanted.link.String() != got.link.String() ||
		wanted.field1 != got.field1 {

		t.Errorf("wanted: '%+v', got: '%+v'", wanted, got)
	}
	//fmt.Printf("%+v\n", got)
}

func TestUnmarshalerSlice(t *testing.T) {
	got := []DataStructure{}
	err := evaluateJS(`a=[{Field1:"foo", Field2:"http://go/"}]`, &got)
	if err != nil {
		t.Error(err)
	}
	wanted := DataStructure{"foo", nil}
	wanted.link, _ = url.Parse("http://go/")
	if len(got) != 1 {
		t.Errorf("wanted: '%+v', got: '%+v'", []DataStructure{wanted}, got)
	}
	if got[0].link == nil || wanted.link.String() != got[0].link.String() ||
		wanted.field1 != got[0].field1 {

		t.Errorf("wanted: '%+v', got: '%+v'", []DataStructure{wanted}, got)
	}
	//fmt.Printf("%+v\n", got)
}
