package qs

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

type timeFormat uint8

const (
	_ timeFormat = iota
	timeFormatSecond
	timeFormatMillis
)

type listFormat uint8

const (
	arrayFormatRepeat listFormat = iota
	arrayFormatBracket
	arrayFormatComma
	arrayFormatIndex
)

// other fields implement baseField
type baseField struct {
	name      string
	omitEmpty bool
}

// embedField represents for nested struct
type embedField struct {
	*baseField
	cachedFields cachedFields
}

func newEmbedField(preAlloc int, tagName []byte, tagOptions [][]byte) *embedField {
	embedField := &embedField{
		baseField: &baseField{
			name: string(tagName),
		},
		cachedFields: make(cachedFields, 0, preAlloc),
	}
	for _, tagOption := range tagOptions {
		if string(tagOption) == tagOmitEmpty {
			embedField.omitEmpty = true
		}
	}
	return embedField
}

func (embedField *embedField) formatFnc(v reflect.Value, result resultFunc) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !embedField.omitEmpty {
				result(embedField.name, "")
			}
			return
		}
		v = v.Elem()
	}
	for i, cachedField := range embedField.cachedFields {
		if cachedField == nil {
			continue
		}
		cachedField.formatFnc(v.Field(i), result)
	}
}

// Present for field with slice/array data type
type listField struct {
	*baseField
	cachedField cachedField
	arrayFormat listFormat
}

func (listField *listField) formatFnc(field reflect.Value, result resultFunc) {
	switch listField.arrayFormat {
	case arrayFormatComma:
		var str strings.Builder
		for i := 0; i < field.Len(); i++ {
			if listField.cachedField == nil {
				return
			}
			elemVal := field.Index(i)
			for elemVal.Kind() == reflect.Ptr {
				elemVal = elemVal.Elem()
			}
			if !elemVal.IsValid() {
				continue
			}
			if i != 0 {
				str.WriteByte(',')
			}
			listField.cachedField.formatFnc(elemVal, func(name string, val string) {
				str.WriteString(val)
			})
		}
		result(listField.name, str.String())
	case arrayFormatRepeat, arrayFormatBracket:
		for i := 0; i < field.Len(); i++ {
			if listField.cachedField == nil {
				return
			}
			elemVal := field.Index(i)
			for elemVal.Kind() == reflect.Ptr {
				elemVal = elemVal.Elem()
			}
			if !elemVal.IsValid() {
				continue
			}
			listField.cachedField.formatFnc(elemVal, func(name string, val string) {
				result(listField.name, val)
			})
		}
	case arrayFormatIndex:
		count := 0
		for i := 0; i < field.Len(); i++ {
			if listField.cachedField == nil {
				return
			}

			elemVal := field.Index(i)
			for elemVal.Kind() == reflect.Ptr {
				elemVal = elemVal.Elem()
			}
			if !elemVal.IsValid() {
				continue
			}
			if v, ok := listField.cachedField.(*embedField); ok {
				v.formatFnc(elemVal, func(name string, val string) {
					var str strings.Builder
					str.WriteString(listField.name)
					str.WriteString(strconv.FormatInt(int64(count), 10))
					str.WriteByte(']')
					str.WriteByte('[')
					str.WriteString(name)
					str.WriteByte(']')
					result(str.String(), val)
					count++
				})
				continue
			}
			listField.cachedField.formatFnc(elemVal, func(name string, val string) {
				var key strings.Builder
				key.WriteString(listField.name)
				key.WriteString(strconv.FormatInt(int64(count), 10))
				key.WriteString("]")
				result(key.String(), val)
				count++
			})
		}
	}
}

func (e *encoder) newListField(elemTyp reflect.Type, tagName []byte, tagOptions [][]byte) *listField {
	removeIdx := -1
	for i, tagOption := range tagOptions {
		if string(tagOption) == tagOmitEmpty {
			removeIdx = i
		}
	}
	if removeIdx > -1 {
		tagOptions = append(tagOptions[:removeIdx], tagOptions[removeIdx+1:]...)
	}

	listField := &listField{
		cachedField: e.newCacheFieldByType(elemTyp, nil, tagOptions),
	}

	for _, tagOption := range tagOptions {
		switch string(tagOption) {
		case "comma":
			listField.arrayFormat = arrayFormatComma
		case "bracket":
			listField.arrayFormat = arrayFormatBracket
		case "index":
			listField.arrayFormat = arrayFormatIndex
		}
	}

	if field, ok := listField.cachedField.(*embedField); ok {
		e.structCaching(&field.cachedFields, reflect.Zero(elemTyp), nil)
	}

	switch listField.arrayFormat {
	case arrayFormatRepeat, arrayFormatBracket:
		if listField.arrayFormat >= arrayFormatBracket {
			tagName = append(tagName, '[')
			tagName = append(tagName, ']')
		}
	case arrayFormatIndex:
		tagName = append(tagName, '[')
	}
	listField.baseField = &baseField{
		name: string(tagName),
	}
	return listField
}

//
type boolField struct {
	*baseField
	useInt bool
}

func (boolField *boolField) formatFnc(v reflect.Value, result resultFunc) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !boolField.omitEmpty {
				result(boolField.name, "")
			}
			return
		}
		v = v.Elem()
	}
	b := v.Bool()
	if !b && boolField.omitEmpty {
		return
	}
	if boolField.useInt {
		if v.Bool() {
			result(boolField.name, "1")
		} else {
			result(boolField.name, "0")
		}
	} else {
		result(boolField.name, strconv.FormatBool(b))
	}
}

func newBoolField(tagName []byte, tagOptions [][]byte) *boolField {
	field := &boolField{
		baseField: &baseField{
			name: string(tagName),
		},
	}
	for _, tagOption := range tagOptions {
		switch string(tagOption) {
		case tagOmitEmpty:
			field.omitEmpty = true
		case "int":
			field.useInt = true
		}
	}
	return field
}

// Int field
type intField struct {
	*baseField
}

func (intField *intField) formatFnc(value reflect.Value, result resultFunc) {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			if !intField.omitEmpty {
				result(intField.name, "")
			}
			return
		}
		value = value.Elem()
	}
	i := value.Int()
	if i == 0 && intField.omitEmpty {
		return
	}
	result(intField.name, strconv.FormatInt(i, 10))
}

func newIntField(tagName []byte, tagOptions [][]byte) *intField {
	field := &intField{
		baseField: &baseField{
			name: string(tagName),
		},
	}
	for _, tagOption := range tagOptions {
		if string(tagOption) == tagOmitEmpty {
			field.omitEmpty = true
		}
	}
	return field
}

// Uint field
type uintField struct {
	*baseField
}

func (uintField *uintField) formatFnc(value reflect.Value, result resultFunc) {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			if !uintField.omitEmpty {
				result(uintField.name, "")
			}
			return
		}
		value = value.Elem()
	}
	i := value.Uint()
	if i == 0 && uintField.omitEmpty {
		return
	}
	result(uintField.name, strconv.FormatUint(i, 10))
}

func newUintField(tagName []byte, tagOptions [][]byte) *uintField {
	field := &uintField{
		baseField: &baseField{
			name: string(tagName),
		},
	}
	for _, tagOption := range tagOptions {
		if string(tagOption) == tagOmitEmpty {
			field.omitEmpty = true
		}
	}
	return field
}

// String field
type stringField struct {
	*baseField
}

func (stringField *stringField) formatFnc(value reflect.Value, result resultFunc) {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			if !stringField.omitEmpty {
				result(stringField.name, "")
			}
			return
		}
		value = value.Elem()
	}
	str := value.String()
	if str == "" && stringField.omitEmpty {
		return
	}
	result(stringField.name, str)
}

func newStringField(tagName []byte, tagOptions [][]byte) *stringField {
	field := &stringField{
		baseField: &baseField{
			name: string(tagName),
		},
	}
	for _, tagOption := range tagOptions {
		if string(tagOption) == tagOmitEmpty {
			field.omitEmpty = true
		}
	}
	return field
}

// Float32 field
type float32Field struct {
	*baseField
}

func (float32Field *float32Field) formatFnc(value reflect.Value, result resultFunc) {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			if !float32Field.omitEmpty {
				result(float32Field.name, "")
			}
			return
		}
		value = value.Elem()
	}
	f := value.Float()
	if f == 0 && float32Field.omitEmpty {
		return
	}
	result(float32Field.name, strconv.FormatFloat(f, 'f', -1, 32))
}

func newFloat32Field(tagName []byte, tagOptions [][]byte) *float32Field {
	field := &float32Field{
		baseField: &baseField{
			name: string(tagName),
		},
	}
	for _, tagOption := range tagOptions {
		if string(tagOption) == tagOmitEmpty {
			field.omitEmpty = true
		}
	}
	return field
}

// Float64 field
type float64Field struct {
	*baseField
}

func (float64Field *float64Field) formatFnc(v reflect.Value, result resultFunc) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !float64Field.omitEmpty {
				result(float64Field.name, "")
			}
			return
		}
		v = v.Elem()
	}
	f := v.Float()
	if f == 0 && float64Field.omitEmpty {
		return
	}
	result(float64Field.name, strconv.FormatFloat(f, 'f', -1, 64))
}

func newFloat64Field(tagName []byte, tagOptions [][]byte) *float64Field {
	field := &float64Field{
		baseField: &baseField{
			name: string(tagName),
		},
	}
	for _, tagOption := range tagOptions {
		if string(tagOption) == tagOmitEmpty {
			field.omitEmpty = true
		}
	}
	return field
}

// Complex64 field
type complex64Field struct {
	*baseField
}

func (complex64Field *complex64Field) formatFnc(v reflect.Value, result resultFunc) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !complex64Field.omitEmpty {
				result(complex64Field.name, "")
			}
			return
		}
		v = v.Elem()
	}
	c := v.Complex()
	if c == 0 && complex64Field.omitEmpty {
		return
	}
	result(complex64Field.name, strconv.FormatComplex(c, 'f', -1, 64))
}

func newComplex64Field(tagName []byte, tagOptions [][]byte) *complex64Field {
	field := &complex64Field{
		baseField: &baseField{
			name: string(tagName),
		},
	}
	for _, tagOption := range tagOptions {
		if string(tagOption) == tagOmitEmpty {
			field.omitEmpty = true
		}
	}
	return field
}

// Complex64 field
type complex128Field struct {
	*baseField
}

func (complex128Field *complex128Field) formatFnc(v reflect.Value, result resultFunc) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !complex128Field.omitEmpty {
				result(complex128Field.name, "")
			}
			return
		}
		v = v.Elem()
	}
	c := v.Complex()
	if c == 0 && complex128Field.omitEmpty {
		return
	}
	result(complex128Field.name, strconv.FormatComplex(c, 'f', -1, 128))
}

func newComplex128Field(tagName []byte, tagOptions [][]byte) *complex128Field {
	field := &complex128Field{
		baseField: &baseField{
			name: string(tagName),
		},
	}
	for _, tagOption := range tagOptions {
		if string(tagOption) == tagOmitEmpty {
			field.omitEmpty = true
		}
	}
	return field
}

// Time field
type timeField struct {
	*baseField
	timeFormat timeFormat
}

func (timeField *timeField) formatFnc(v reflect.Value, result resultFunc) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !timeField.omitEmpty {
				result(timeField.name, "")
			}
			return
		}
		v = v.Elem()
	}
	t := v.Interface().(time.Time)
	if t.IsZero() && timeField.omitEmpty {
		return
	}
	switch timeField.timeFormat {
	case timeFormatSecond:
		result(timeField.name, strconv.FormatInt(t.Unix(), 10))
	case timeFormatMillis:
		result(timeField.name, strconv.FormatInt(t.UnixNano()/1000000, 10))
	default:
		result(timeField.name, t.Format(time.RFC3339))
	}
}

func newTimeField(tagName []byte, tagOptions [][]byte) *timeField {
	field := &timeField{
		baseField: &baseField{
			name: string(tagName),
		},
	}
	for _, tagOption := range tagOptions {
		switch string(tagOption) {
		case tagOmitEmpty:
			field.omitEmpty = true
		case "second":
			field.timeFormat = timeFormatSecond
		case "millis":
			field.timeFormat = timeFormatMillis
		}
	}
	return field
}

type customField struct {
	*baseField
	tagOptions []string
	formatter  func(val interface{}, opts []string, result func(v string))
}

func (customField *customField) formatFnc(v reflect.Value, result resultFunc) {
	customField.formatter(v.Interface(), customField.tagOptions, func(v string) {
		result(customField.name, v)
	})
}

func newCustomField(tagName []byte, tagOptions [][]byte, formatter func(val interface{}, opts []string, result func(v string))) *customField {
	opts := make([]string, 0, len(tagOptions))
	for _, tagOption := range tagOptions {
		opts = append(opts, string(tagOption))
	}
	return &customField{
		baseField: &baseField{
			name:      string(tagName),
			omitEmpty: false,
		},
		tagOptions: opts,
		formatter:  formatter,
	}
}
