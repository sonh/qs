package qs

import (
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
	timeType    = reflect.TypeOf(time.Time{})
	encoderType = reflect.TypeOf(new(QueryParamEncoder)).Elem()
	zeroerType  = reflect.TypeOf(new(Zeroer)).Elem()
)

// EncoderOption provides option for Encoder
type EncoderOption func(encoder *Encoder)

// Encoder is the main instance
// Apply options by using WithTagAlias, WithCustomType
type Encoder struct {
	tagAlias string
	cache    *cacheStore
	dataPool *sync.Pool
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

// NewEncoder init new *Encoder instance
// Use EncoderOption to apply options
func NewEncoder(options ...EncoderOption) *Encoder {
	e := &Encoder{
		tagAlias: "qs",
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
			tags = append(tags, make([]byte, 0, 56))
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
			return nil, ErrInvalidInput{inputType: val.Kind().String()}
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Invalid:
		return nil, ErrInvalidInput{inputType: val.Kind().String()}
	case reflect.Struct:
		enc := e.dataPool.Get().(*encoder)
		enc.values = make(url.Values)
		err := enc.encodeStruct(val, enc.values, nil)
		if err != nil {
			return nil, err
		}
		values := enc.values
		e.dataPool.Put(enc)
		return values, nil
	default:
		return nil, ErrInvalidInput{inputType: val.Kind().String()}
	}
}

// Encode encodes a struct into the given url.Values
// v must be struct data type
func (e *Encoder) Encode(v interface{}, values url.Values) error {
	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ErrInvalidInput{inputType: val.Kind().String()}
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Invalid:
		return ErrInvalidInput{inputType: val.Kind().String()}
	case reflect.Struct:
		enc := e.dataPool.Get().(*encoder)
		err := enc.encodeStruct(val, values, nil)
		if err != nil {
			return err
		}
		return nil
	default:
		return ErrInvalidInput{inputType: val.Kind().String()}
	}
}

func (e *encoder) encodeStruct(stVal reflect.Value, values url.Values, scope []byte) error {
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
		case *mapField:
			if cachedFld.cachedKeyField == nil || cachedFld.cachedValueField == nil {
				// data type is not supported
				continue
			}
			for stFldVal.Kind() == reflect.Ptr {
				stFldVal = stFldVal.Elem()
			}
			if !stFldVal.IsValid() {
				continue
			}
			if stFldVal.Len() == 0 {
				continue
			}
		case *listField:
			if cachedFld.cachedField == nil {
				// data type is not supported
				continue
			}
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
				if count := countElem(stFldVal); count > 0 {
					// preallocate slice
					values[cachedFld.name] = make([]string, 0, countElem(stFldVal))
				} else {
					continue
				}
			}
		}

		// format value
		err := cachedFld.formatFnc(stFldVal, func(name string, val string) {
			values[name] = append(values[name], val)
		})
		if err != nil {
			return err
		}
	}
	return nil
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

		if fieldVal.Type().Implements(encoderType) {
			*fields = append(*fields, newCustomField(fieldVal.Type(), e.tags[0], e.tags[1:]))
			continue
		}

		fieldTyp := getType(fieldVal)

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
			field := newEmbedField(fieldVal.NumField(), e.tags[0], e.tags[1:])
			*fields = append(*fields, field)
			// Recursive
			e.structCaching(&field.cachedFields, fieldVal, e.scope)
		case reflect.Slice, reflect.Array:
			// Slice element type
			elemType := fieldTyp.Elem()
			if elemType.Implements(encoderType) {
				*fields = append(*fields, e.newListField(elemType, e.tags[0], e.tags[1:]))
				continue
			}
			for elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}
			*fields = append(*fields, e.newListField(elemType, e.tags[0], e.tags[1:]))
		case reflect.Map:
			keyType := fieldTyp.Key()
			/*for keyType.Kind() == reflect.Ptr {
				keyType = keyType.Elem()
			}*/
			valueType := fieldTyp.Elem()
			/*for valueType.Kind() == reflect.Ptr {
				valueType = valueType.Elem()
			}*/
			*fields = append(*fields, newMapField(keyType, valueType, e.tags[0], e.tags[1:]))
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
		e.tags[0] = append(e.tags[0][:0], f.Name...)
		e.tags = e.tags[:1]
	} else {
		// Use first tag as temp
		e.tags[0] = append(e.tags[0][:0], tag...)

		splitTags := strings.Split(tag, ",")
		e.tags = e.tags[:len(splitTags)]

		for i := 0; i < len(splitTags); i++ {
			if i == 0 {
				if len(splitTags[0]) == 0 {
					e.tags[0] = append(e.tags[i][:0], f.Name...)
					continue
				}
			}
			e.tags[i] = append(e.tags[i][:0], splitTags[i]...)
		}
	}
}
