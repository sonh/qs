package qs

import (
	"reflect"
	"testing"
)

func TestCacheStore(t *testing.T) {
	t.Parallel()
	//test := assert.New(t)

	s := &basicVal{}

	cacheStore := newCacheStore()
	if cacheStore == nil {
		t.FailNow()
	}
	//test.NotNil(cacheStore)

	fields := cachedFields{&float64Field{}}
	cacheStore.Store(reflect.TypeOf(s), fields)
	cachedFlds := cacheStore.Retrieve(reflect.TypeOf(s))

	if cachedFlds == nil {
		t.FailNow()
	}
	if len(cachedFlds) != len(fields) {
		t.FailNow()
	}
	if &fields[0] != &cachedFlds[0] {
		t.FailNow()
	}
	//test.NotNil(cachedFlds)
	//test.Len(cachedFlds, len(fields))
	//test.True(&fields[0] == &cachedFlds[0])
}

func TestNewCacheField(t *testing.T) {
	t.Parallel()

	name := []byte(`abc`)
	opts := [][]byte{[]byte(`omitempty`)}

	cacheField := newCachedFieldByKind(reflect.ValueOf("").Kind(), name, opts)

	strField, ok := cacheField.(*stringField)
	if !ok {
		t.FailNow()
	}
	if string(name) != strField.name {
		t.FailNow()
	}
	if !strField.omitEmpty {
		t.FailNow()
	}
	if !reflect.DeepEqual(reflect.TypeOf(new(stringField)), reflect.TypeOf(cacheField)) {
		t.FailNow()
	}
}

func TestNewCacheField2(t *testing.T) {
	//test := assert.New(t)

	var strPtr *string
	cacheField := newCachedFieldByKind(reflect.ValueOf(strPtr).Kind(), nil, nil)
	if cacheField != nil {
		t.FailNow()
	}
	//test.Nil(cacheField)
}
