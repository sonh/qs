package qs

import (
	"bytes"
	"github.com/pkg/errors"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	tagOmitEmpty = "omitempty"
)

var (
	timeType = reflect.TypeOf(time.Time{})
)

// EncoderOption provides option for Encoder
type EncoderOption func(encoder *Encoder)

// Encoder is the main instance
// Apply options by using WithTagAlias, WithCustomType
type Encoder struct {
	tagAlias   string
	cache      *cacheStore
	dataPool   *sync.Pool
	formatters map[reflect.Type]func(val interface{}, opts []string, result func(v string))
}

type encoder struct {
	e      *Encoder
	values url.Values
	tags   [][]byte
	scope  []byte
}

// WithTagAlias create a option to set custom tag alias instead of `qs`
func WithTagAlias(tagAlias string) EncoderOption {
	return func(encoder *Encoder) {
		encoder.tagAlias = tagAlias
	}
}

// WithCustomType create a option to set custom data type
func WithCustomType(i interface{}, formatter func(val interface{}, opts []string, result func(v string))) EncoderOption {
	return func(encoder *Encoder) {
		switch typ := i.(type) {
		case reflect.Type:
			encoder.formatters[typ] = formatter
		default:
			encoder.formatters[reflect.TypeOf(i)] = formatter
		}
	}
}

// NewEncoder init new *Encoder instance
// Use EncoderOption to apply options
func NewEncoder(options ...EncoderOption) *Encoder {
	e := &Encoder{
		tagAlias:   "qs",
		formatters: make(map[reflect.Type]func(val interface{}, opts []string, result func(v string))),
	}

	// Apply options
	for _, opt := range options {
		opt(e)
	}

	e.cache = newCacheStore()

	e.dataPool = &sync.Pool{New: func() interface{} {
		tagSize := 5
		tags := make([][]byte, 0, tagSize)
		for i := 0; i < tagSize; i++ {
			tags = append(tags, make([]byte, 0, 64))
		}
		return &encoder{
			e:     e,
			tags:  tags,
			scope: make([]byte, 0, 64),
		}
	}}

	return e
}

// Values encodes a struct into url.Values
// v must be struct data type
func (e *Encoder) Values(v interface{}) (url.Values, error) {
	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, errors.Errorf("expects struct input, got %v", val.Kind())
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Invalid:
		return nil, errors.Errorf("expects struct input, got %v", val.Kind())
	case reflect.Struct:
		enc := e.dataPool.Get().(*encoder)
		enc.values = make(url.Values)
		enc.encodeStruct(val, enc.values, nil)
		values := enc.values
		e.dataPool.Put(enc)
		return values, nil
	default:
		return nil, errors.Errorf("expects struct input, got %v", val.Kind())
	}
}

// Encode encodes a struct into the given url.Values
// v must be struct data type
func (e *Encoder) Encode(v interface{}, values url.Values) error {
	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return errors.Errorf("expects struct input, got %v", val.Kind())
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Invalid:
		return errors.Errorf("expects struct input, got %v", val.Kind())
	case reflect.Struct:
		enc := e.dataPool.Get().(*encoder)
		enc.encodeStruct(val, values, nil)
		return nil
	default:
		return errors.Errorf("expects struct input, got %v", val.Kind())
	}
}

func (e *encoder) encodeStruct(stVal reflect.Value, values url.Values, scope []byte) {
	stTyp := stVal.Type()

	cachedFlds := e.e.cache.Retrieve(stTyp)

	if cachedFlds == nil {
		cachedFlds = make(cachedFields, 0, stTyp.NumField())
		e.structCaching(&cachedFlds, stVal, scope)
		e.e.cache.Store(stTyp, cachedFlds)
	}

	for i, cachedFld := range cachedFlds {
		stFldVal := stVal.Field(i)

		switch cachedFld := cachedFld.(type) {
		case nil:
			// skip field
			continue
		case *listField:
			if cachedFld.arrayFormat <= arrayFormatBracket {
				// With cachedFld type is slice/array, only accept non-nil value
				for stFldVal.Kind() == reflect.Ptr {
					stFldVal = stFldVal.Elem()
				}
				if !stFldVal.IsValid() {
					continue
				}
				if stFldVal.Len() == 0 {
					continue
				}
				// preallocate slice
				values[cachedFld.name] = make([]string, 0, stFldVal.Len())
			}
		}
		// format value
		cachedFld.formatFnc(stFldVal, func(name string, val string) {
			values[name] = append(values[name], val)
		})
	}
}

func (e *encoder) structCaching(fields *cachedFields, stVal reflect.Value, scope []byte) {

	structTyp := getType(stVal)

	for i := 0; i < structTyp.NumField(); i++ {

		structField := structTyp.Field(i)

		if structField.PkgPath != "" && !structField.Anonymous { // unexported field
			*fields = append(*fields, nil)
			continue
		}

		e.getTagNameAndOpts(structField)

		if string(e.tags[0]) == "-" { // ignored field
			continue
		}

		if string(scope) != "" {
			scopedName := strings.Builder{}
			scopedName.Write(scope)
			scopedName.WriteRune('[')
			scopedName.Write(e.tags[0])
			scopedName.WriteRune(']')
			e.tags[0] = e.tags[0][:0]
			e.tags[0] = append(e.tags[0], scopedName.String()...)
		}

		fieldVal := stVal.Field(i)
		fieldTyp := getType(fieldVal)

		if formatter := e.e.formatters[fieldTyp]; formatter != nil {
			*fields = append(*fields, newCustomField(e.tags[0], e.tags[1:], formatter))
			continue
		}

		if fieldTyp == timeType {
			*fields = append(*fields, newTimeField(e.tags[0], e.tags[1:]))
			continue
		}

		switch fieldTyp.Kind() {
		case reflect.Struct:
			fieldVal = reflect.Zero(fieldTyp)
			// Clear and set new scope
			e.scope = e.scope[:0]
			e.scope = append(e.scope, e.tags[0]...)
			// New embed field
			field := newEmbedField(fieldVal, e.tags[0], e.tags[1:])
			*fields = append(*fields, field)
			// Recursive
			e.structCaching(&field.cachedFields, fieldVal, e.scope)
		case reflect.Slice, reflect.Array:
			//Slice element type
			elemType := fieldTyp.Elem()
			for elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}
			*fields = append(*fields, e.newListField(elemType, e.tags[0], e.tags[1:]))
		case reflect.Ptr:
			*fields = append(*fields, nil)
		default:
			*fields = append(*fields, newCachedFieldByKind(fieldTyp.Kind(), e.tags[0], e.tags[1:]))
		}
	}
}

func (e *encoder) getTagNameAndOpts(f reflect.StructField) {
	// Get tag by alias
	tag := f.Tag.Get(e.e.tagAlias)

	// Clear first tag in slice
	e.tags[0] = e.tags[0][:0]

	if len(tag) == 0 {
		// no tag, using struct field name
		e.tags[0] = append(e.tags[0], f.Name...)
		e.tags = e.tags[:1]
	} else {
		// Use first tag as temp
		e.tags[0] = append(e.tags[0], tag...)
		splitTags := bytes.Split(e.tags[0], []byte{','})

		e.tags = e.tags[:len(splitTags)]
		for i := 0; i < len(splitTags); i++ {
			// Clear this tag and set new tag
			e.tags[i] = e.tags[i][:0]
			e.tags[i] = append(e.tags[i], splitTags[i]...)
		}
	}
}
