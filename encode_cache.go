package qs

import (
	"reflect"
	"sync"
)

type cacheStore struct {
	m     map[reflect.Type]cachedFields
	mutex sync.RWMutex
}

func newCacheStore() *cacheStore {
	return &cacheStore{
		m: make(map[reflect.Type]cachedFields),
	}
}

// Retrieve cachedFields corresponding to reflect.Type
func (cacheStore *cacheStore) Retrieve(typ reflect.Type) cachedFields {
	return cacheStore.m[typ]
}

// cacheStore stores cachedFields that corresponds to reflect.Type
func (cacheStore *cacheStore) Store(typ reflect.Type, cachedFields cachedFields) {
	cacheStore.mutex.Lock()
	defer cacheStore.mutex.Unlock()
	if _, ok := cacheStore.m[typ]; !ok {
		cacheStore.m[typ] = cachedFields
	}
}

type (
	resultFunc func(name string, val string)

	// cachedField
	cachedField interface {
		formatFnc(value reflect.Value, result resultFunc) error
	}

	cachedFields []cachedField
)

func newCacheFieldByType(typ reflect.Type, tagName []byte, tagOptions [][]byte) cachedField {
	if typ.Implements(encoderType) {
		return newCustomField(typ, tagName, tagOptions)
	}
	switch typ {
	case timeType:
		return newTimeField(tagName, tagOptions)
	default:
		return newCachedFieldByKind(typ.Kind(), tagName, tagOptions)
	}
}

func newCachedFieldByKind(kind reflect.Kind, tagName []byte, tagOptions [][]byte) cachedField {
	switch kind {
	case reflect.String:
		return newStringField(tagName, tagOptions)
	case reflect.Bool:
		return newBoolField(tagName, tagOptions)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return newIntField(tagName, tagOptions)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return newUintField(tagName, tagOptions)
	case reflect.Float32:
		return newFloat32Field(tagName, tagOptions)
	case reflect.Float64:
		return newFloat64Field(tagName, tagOptions)
	case reflect.Complex64:
		return newComplex64Field(tagName, tagOptions)
	case reflect.Complex128:
		return newComplex128Field(tagName, tagOptions)
	case reflect.Struct:
		return newEmbedField(0, tagName, tagOptions)
	case reflect.Interface:
		return newInterfaceField(tagName, tagOptions)
	default:
		return nil
	}
}

func getType(fieldVal reflect.Value) reflect.Type {
	stFieldTyp := fieldVal.Type()
	for fieldVal.Kind() == reflect.Ptr {
		fieldVal = fieldVal.Elem()
		stFieldTyp = stFieldTyp.Elem()
	}
	return stFieldTyp
}

func countElem(value reflect.Value) int {
	count := 0
	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)
		for elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		if elem.IsValid() {
			count++
		}
	}
	return count
}
