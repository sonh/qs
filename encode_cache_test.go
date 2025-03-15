package qs

import (
	"reflect"
	"testing"
)

func TestCacheStore(t *testing.T) {
	t.Parallel()

	s := &basicVal{}

	cacheStore := newCacheStore()
	if cacheStore == nil {
		t.Error("cache store should not be nil")
		t.FailNow()
	}

	fields := cachedFields{&float64Field{}}
	cacheStore.Store(reflect.TypeOf(s), fields)
	cachedFlds := cacheStore.Retrieve(reflect.TypeOf(s))

	if cachedFlds == nil {
		t.Error("cache store should not be nil")
		t.FailNow()
	}
	if len(cachedFlds) != len(fields) {
		t.Error("cache store should have the same number of fields")
		t.FailNow()
	}
	if &fields[0] != &cachedFlds[0] {
		t.Error("cache store should have the same fields")
		t.FailNow()
	}
}

func TestNewCacheField(t *testing.T) {
	t.Parallel()

	name := []byte(`abc`)
	opts := [][]byte{[]byte(`omitempty`)}

	cacheField := newCachedFieldByKind(reflect.ValueOf("").Kind(), name, opts)

	strField, ok := cacheField.(*stringField)
	if !ok {
		t.Error("strField should be stringField")
		t.FailNow()
	}
	if string(name) != strField.name {
		t.Errorf("strField.name should be %s, but %s", string(name), strField.name)
		t.FailNow()
	}
	if !strField.omitEmpty {
		t.Error("omitEmpty should be true")
		t.FailNow()
	}
	if !reflect.DeepEqual(reflect.TypeOf(new(stringField)), reflect.TypeOf(cacheField)) {
		t.Error("cache field is not of type *stringField")
		t.FailNow()
	}
}

func TestNewCacheField2(t *testing.T) {
	t.Parallel()

	var strPtr *string
	cacheField := newCachedFieldByKind(reflect.ValueOf(strPtr).Kind(), nil, nil)
	if cacheField != nil {
		t.Error("expect cacheField to be nil")
		t.FailNow()
	}
}
