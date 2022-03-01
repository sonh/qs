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

func (embedField *embedField) formatFnc(v reflect.Value, result resultFunc) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !embedField.omitEmpty {
				result(embedField.name, "")
			}
			return nil
		}
		v = v.Elem()
	}
	for i, cachedField := range embedField.cachedFields {
		if cachedField == nil {
			continue
		}
		err := cachedField.formatFnc(v.Field(i), result)
		if err != nil {
			return err
		}
	}
	return nil
}

// Present for field with slice/array data type
type listField struct {
	*baseField
	cachedField cachedField
	arrayFormat listFormat
}

func (listField *listField) formatFnc(field reflect.Value, result resultFunc) error {
	switch listField.arrayFormat {
	case arrayFormatComma:
		var str strings.Builder
		for i := 0; i < field.Len(); i++ {
			elemVal := field.Index(i)
			if _, ok := listField.cachedField.(*customField); ok {
				elem := elemVal
				for elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				if !elem.IsValid() {
					continue
				}
			} else {
				for elemVal.Kind() == reflect.Ptr {
					elemVal = elemVal.Elem()
				}
				if !elemVal.IsValid() {
					continue
				}
			}
			err := listField.cachedField.formatFnc(elemVal, func(name string, val string) {
				if i > 0 {
					str.WriteByte(',')
				}
				str.WriteString(val)
			})
			if err != nil {
				return err
			}
		}
		returnStr := str.String()
		if returnStr[0] == ',' {
			returnStr = returnStr[1:]
		}
		result(listField.name, returnStr)
	case arrayFormatRepeat, arrayFormatBracket:
		for i := 0; i < field.Len(); i++ {
			elemVal := field.Index(i)
			if _, ok := listField.cachedField.(*customField); ok {
				elem := elemVal
				for elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				if !elem.IsValid() {
					continue
				}
			} else {
				for elemVal.Kind() == reflect.Ptr {
					elemVal = elemVal.Elem()
				}
				if !elemVal.IsValid() {
					continue
				}
			}
			err := listField.cachedField.formatFnc(elemVal, func(name string, val string) {
				result(listField.name, val)
			})
			if err != nil {
				return err
			}
		}
	case arrayFormatIndex:
		count := 0
		for i := 0; i < field.Len(); i++ {
			elemVal := field.Index(i)
			if _, ok := listField.cachedField.(*customField); ok {
				elem := elemVal
				for elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				if !elem.IsValid() {
					continue
				}
			} else {
				for elemVal.Kind() == reflect.Ptr {
					elemVal = elemVal.Elem()
				}
				if !elemVal.IsValid() {
					continue
				}
			}
			if v, ok := listField.cachedField.(*embedField); ok {
				err := v.formatFnc(elemVal, func(name string, val string) {
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
				if err != nil {
					return err
				}
				continue
			}
			err := listField.cachedField.formatFnc(elemVal, func(name string, val string) {
				var key strings.Builder
				key.WriteString(listField.name)
				key.WriteString(strconv.FormatInt(int64(count), 10))
				key.WriteString("]")
				result(key.String(), val)
				count++
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
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
		cachedField: newCacheFieldByType(elemTyp, nil, tagOptions),
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

	if field, ok := listField.cachedField.(*embedField); ok {
		e.structCaching(&field.cachedFields, reflect.Zero(elemTyp), nil)
	}

	return listField
}

type mapField struct {
	*baseField
	cachedKeyField   cachedField
	cachedValueField cachedField
}

func (mapField *mapField) formatFnc(field reflect.Value, result resultFunc) error {
	mapRange := field.MapRange()
	fieldName := make([]byte, 0, 36)
	fieldName = append(fieldName, mapField.name...)

	for mapRange.Next() {
		fieldName = fieldName[:len(mapField.name)]
		err := mapField.cachedKeyField.formatFnc(mapRange.Key(), func(_ string, val string) {
			fieldName = append(fieldName, '[')
			fieldName = append(fieldName, val...)
			fieldName = append(fieldName, ']')
		})
		if err != nil {
			return err
		}
		err = mapField.cachedValueField.formatFnc(mapRange.Value(), func(_ string, val string) {
			result(string(fieldName), val)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func newMapField(keyType reflect.Type, valueType reflect.Type, tagName []byte, tagOptions [][]byte) *mapField {
	removeIdx := -1
	for i, tagOption := range tagOptions {
		if string(tagOption) == tagOmitEmpty {
			removeIdx = i
		}
	}
	if removeIdx > -1 {
		tagOptions = append(tagOptions[:removeIdx], tagOptions[removeIdx+1:]...)
	}

	if !keyType.Implements(encoderType) {
		for keyType.Kind() == reflect.Ptr {
			keyType = keyType.Elem()
		}
	}

	if !valueType.Implements(encoderType) {
		for valueType.Kind() == reflect.Ptr {
			valueType = valueType.Elem()
		}
	}

	field := &mapField{
		baseField: &baseField{
			name: string(tagName),
		},
		cachedKeyField:   newCacheFieldByType(keyType, nil, nil),
		cachedValueField: newCacheFieldByType(valueType, nil, nil),
	}
	return field
}

//
type boolField struct {
	*baseField
	useInt bool
}

func (boolField *boolField) formatFnc(v reflect.Value, result resultFunc) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !boolField.omitEmpty {
				result(boolField.name, "")
			}
			return nil
		}
		v = v.Elem()
	}
	b := v.Bool()
	if !b && boolField.omitEmpty {
		return nil
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
	return nil
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

func (intField *intField) formatFnc(value reflect.Value, result resultFunc) error {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			if !intField.omitEmpty {
				result(intField.name, "")
			}
			return nil
		}
		value = value.Elem()
	}
	i := value.Int()
	if i == 0 && intField.omitEmpty {
		return nil
	}
	result(intField.name, strconv.FormatInt(i, 10))
	return nil
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

func (uintField *uintField) formatFnc(value reflect.Value, result resultFunc) error {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			if !uintField.omitEmpty {
				result(uintField.name, "")
			}
			return nil
		}
		value = value.Elem()
	}
	i := value.Uint()
	if i == 0 && uintField.omitEmpty {
		return nil
	}
	result(uintField.name, strconv.FormatUint(i, 10))
	return nil
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

func (stringField *stringField) formatFnc(value reflect.Value, result resultFunc) error {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			if !stringField.omitEmpty {
				result(stringField.name, "")
			}
			return nil
		}
		value = value.Elem()
	}
	str := value.String()
	if str == "" && stringField.omitEmpty {
		return nil
	}
	result(stringField.name, str)
	return nil
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

func (float32Field *float32Field) formatFnc(value reflect.Value, result resultFunc) error {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			if !float32Field.omitEmpty {
				result(float32Field.name, "")
			}
			return nil
		}
		value = value.Elem()
	}
	f := value.Float()
	if f == 0 && float32Field.omitEmpty {
		return nil
	}
	result(float32Field.name, strconv.FormatFloat(f, 'f', -1, 32))
	return nil
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

func (float64Field *float64Field) formatFnc(v reflect.Value, result resultFunc) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !float64Field.omitEmpty {
				result(float64Field.name, "")
			}
			return nil
		}
		v = v.Elem()
	}
	f := v.Float()
	if f == 0 && float64Field.omitEmpty {
		return nil
	}
	result(float64Field.name, strconv.FormatFloat(f, 'f', -1, 64))
	return nil
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

func (complex64Field *complex64Field) formatFnc(v reflect.Value, result resultFunc) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !complex64Field.omitEmpty {
				result(complex64Field.name, "")
			}
			return nil
		}
		v = v.Elem()
	}
	c := v.Complex()
	if c == 0 && complex64Field.omitEmpty {
		return nil
	}
	result(complex64Field.name, strconv.FormatComplex(c, 'f', -1, 64))
	return nil
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

func (complex128Field *complex128Field) formatFnc(v reflect.Value, result resultFunc) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !complex128Field.omitEmpty {
				result(complex128Field.name, "")
			}
			return nil
		}
		v = v.Elem()
	}
	c := v.Complex()
	if c == 0 && complex128Field.omitEmpty {
		return nil
	}
	result(complex128Field.name, strconv.FormatComplex(c, 'f', -1, 128))
	return nil
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

func (timeField *timeField) formatFnc(v reflect.Value, result resultFunc) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !timeField.omitEmpty {
				result(timeField.name, "")
			}
			return nil
		}
		v = v.Elem()
	}
	t := v.Interface().(time.Time)
	if t.IsZero() && timeField.omitEmpty {
		return nil
	}
	switch timeField.timeFormat {
	case timeFormatSecond:
		result(timeField.name, strconv.FormatInt(t.Unix(), 10))
	case timeFormatMillis:
		result(timeField.name, strconv.FormatInt(t.UnixNano()/1000000, 10))
	default:
		result(timeField.name, t.Format(time.RFC3339))
	}
	return nil
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

// Zeroer represents an object has zero value
// IsZeroer is used to check whether an object is zero to
// determine whether it should be omitted when encoding
type Zeroer interface {
	IsZero() bool
}

// QueryParamEncoder is an interface implemented by any type to encode itself into query param
type QueryParamEncoder interface {
	EncodeParam() (string, error)
}

type customField struct {
	*baseField
	isZeroer bool
}

func (customField *customField) formatFnc(v reflect.Value, result resultFunc) error {
	elem := v
	for elem.Kind() == reflect.Ptr {
		elem = v.Elem()
	}
	if !elem.IsValid() {
		if !customField.omitEmpty {
			result(customField.name, "")
		}
		return nil
	}
	valueInterface := v.Interface()
	if customField.isZeroer && valueInterface.(Zeroer).IsZero() {
		if !customField.omitEmpty {
			result(customField.name, "")
		}
		return nil
	}
	str, err := valueInterface.(QueryParamEncoder).EncodeParam()
	if err != nil {
		return err
	}
	result(customField.name, str)
	return nil
}

func newCustomField(typ reflect.Type, tagName []byte, tagOptions [][]byte) *customField {
	field := &customField{
		baseField: &baseField{
			name:      string(tagName),
			omitEmpty: false,
		},
	}

	if typ.Implements(zeroerType) {
		field.isZeroer = true
	}

	for _, tagOption := range tagOptions {
		if string(tagOption) == tagOmitEmpty {
			field.omitEmpty = true
		}
	}
	return field
}

type interfaceField struct {
	*baseField
	tagName    []byte
	tagOptions [][]byte
	fieldMap   map[reflect.Type]cachedField
}

func (interfaceField *interfaceField) formatFnc(v reflect.Value, result resultFunc) error {

	v = v.Elem()

	if v.IsValid() && v.Type().Implements(encoderType) {
		elem := v
		for elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		if !elem.IsValid() {
			return nil
		}
	} else {
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if !v.IsValid() {
			if !interfaceField.omitEmpty {
				result(interfaceField.name, "")
			}
			return nil
		}
	}

	if field := interfaceField.fieldMap[v.Type()]; field == nil {
		interfaceField.fieldMap[v.Type()] = newCacheFieldByType(v.Type(), interfaceField.tagName, interfaceField.tagOptions)
	}
	if field := interfaceField.fieldMap[v.Type()]; field != nil {
		err := field.formatFnc(v, result)
		if err != nil {
			return err
		}
	}
	return nil
}

func newInterfaceField(tagName []byte, tagOptions [][]byte) *interfaceField {
	copiedTagName := make([]byte, len(tagName))
	copy(copiedTagName, tagName)
	copiedTagOptions := make([][]byte, len(tagOptions))
	copy(copiedTagOptions, tagOptions)

	field := &interfaceField{
		baseField: &baseField{
			name: string(copiedTagName),
		},
		tagName:    copiedTagName,
		tagOptions: copiedTagOptions,
		fieldMap:   make(map[reflect.Type]cachedField, 5),
	}
	for _, tagOption := range tagOptions {
		if string(tagOption) == tagOmitEmpty {
			field.omitEmpty = true
		}
	}
	return field
}
