package qs

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCacheStore(t *testing.T) {
	test := assert.New(t)

	s := &basicVal{}

	cacheStore := newCacheStore()
	test.NotNil(cacheStore)

	fields := cachedFields{&float64Field{}}
	cacheStore.Store(reflect.TypeOf(s), fields)
	cachedFlds := cacheStore.Retrieve(reflect.TypeOf(s))

	test.NotNil(cachedFlds)
	test.Len(cachedFlds, len(fields))
	test.True(&fields[0] == &cachedFlds[0])
}

func TestNewCacheField(t *testing.T) {
	test := assert.New(t)
	name := []byte(`abc`)
	opts := [][]byte{[]byte(`omitempty`)}

	cacheField := newCachedFieldByKind(reflect.ValueOf("").Kind(), name, opts)
	if stringField, ok := cacheField.(*stringField); ok {
		test.Equal(string(name), stringField.name)
		test.True(stringField.omitEmpty)
	} else {
		test.FailNow("")
	}
	test.IsType(&stringField{}, cacheField)
}

func TestNewCacheField2(t *testing.T) {
	test := assert.New(t)

	var strPtr *string
	cacheField := newCachedFieldByKind(reflect.ValueOf(strPtr).Kind(), nil, nil)
	test.Nil(cacheField)
}
