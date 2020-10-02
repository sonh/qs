package qs

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

type basicVal struct {
	String  string    `qs:"string"`
	Bool    bool      `qs:"bool"`
	Int     int       `qs:"int"`
	Int8    int8      `qs:"int8"`
	Int16   int16     `qs:"int16"`
	Int32   int32     `qs:"int32"`
	Int64   int64     `qs:"int64"`
	Uint    uint      `qs:"uint"`
	Uint8   uint8     `qs:"uint8"`
	Uint16  uint16    `qs:"uint16"`
	Uint32  uint32    `qs:"uint32"`
	Uint64  uint64    `qs:"uint64"`
	Uintptr	uintptr	  `qs:"uintptr"`
	Float32 float32   `qs:"float32"`
	Float64 float64   `qs:"float64"`
	Time    time.Time `qs:"time"`
}

type basicValWithOmit struct {
	String  string    `qs:"string,omitempty"`
	Bool    bool      `qs:"bool,omitempty"`
	Int     int       `qs:"int,omitempty"`
	Int8    int8      `qs:"int8,omitempty"`
	Int16   int16     `qs:"int16,omitempty"`
	Int32   int32     `qs:"int32,omitempty"`
	Int64   int64     `qs:"int64,omitempty"`
	Uint    uint      `qs:"uint,omitempty"`
	Uint8   uint8     `qs:"uint8,omitempty"`
	Uint16  uint16    `qs:"uint16,omitempty"`
	Uint32  uint32    `qs:"uint32,omitempty"`
	Uint64  uint64    `qs:"uint64,omitempty"`
	Float32 float32   `qs:"float32,omitempty"`
	Float64 float64   `qs:"float64,omitempty"`
	Time    time.Time `qs:"time,omitempty"`
}

type basicPtr struct {
	String  *string    `qs:"string"`
	Bool    *bool      `qs:"bool"`
	Int     *int       `qs:"int"`
	Int8    *int8      `qs:"int8"`
	Int16   *int16     `qs:"int16"`
	Int32   *int32     `qs:"int32"`
	Int64   *int64     `qs:"int64"`
	Uint    *uint      `qs:"uint"`
	Uint8   *uint8     `qs:"uint8"`
	Uint16  *uint16    `qs:"uint16"`
	Uint32  *uint32    `qs:"uint32"`
	Uint64  *uint64    `qs:"uint64"`
	UinPtr  *uintptr   `qs:"uintptr"`
	Float32 *float32   `qs:"float32"`
	Float64 *float64   `qs:"float64"`
	Time    *time.Time `qs:"time"`
}

type basicPtrWithOmit struct {
	String  *string    `qs:"string,omitempty"`
	Bool    *bool      `qs:"bool,omitempty"`
	Int     *int       `qs:"int,omitempty"`
	Int8    *int8      `qs:"int8,omitempty"`
	Int16   *int16     `qs:"int16,omitempty"`
	Int32   *int32     `qs:"int32,omitempty"`
	Int64   *int64     `qs:"int64,omitempty"`
	Uint    *uint      `qs:"uint,omitempty"`
	Uint8   *uint8     `qs:"uint8,omitempty"`
	Uint16  *uint16    `qs:"uint16,omitempty"`
	Uint32  *uint32    `qs:"uint32,omitempty"`
	Uint64  *uint64    `qs:"uint64,omitempty"`
	Float32 *float32   `qs:"float32,omitempty"`
	Float64 *float64   `qs:"float64,omitempty"`
	Time    *time.Time `qs:"time,omitempty"`
}

func TestIgnore(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	v := struct {
		Test string `qs:"-"`
	}{}

	values, err := encoder.Values(v)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	assert.Equal(t, url.Values{}, values)
}

func TestEncodeBasicVal(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	tm := time.Unix(600, 0).UTC()

	s := basicVal{
		String:  "abc",
		Bool:    true,
		Int:     12,
		Int8:    int8(8),
		Int16:   int16(16),
		Int32:   int32(32),
		Int64:   int64(64),
		Uint:    24,
		Uint8:   uint8(8),
		Uint16:  uint16(16),
		Uint32:  uint32(32),
		Uint64:  uint64(64),
		Uintptr: uintptr(72),
		Float32: float32(0.1234),
		Float64: 1.2345,
		Time:    tm,
	}
	values, err := encoder.Values(s)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"string":  []string{"abc"},
		"bool":    []string{"true"},
		"int":     []string{"12"},
		"int8":    []string{"8"},
		"int16":   []string{"16"},
		"int32":   []string{"32"},
		"int64":   []string{"64"},
		"uint":    []string{"24"},
		"uint8":   []string{"8"},
		"uint16":  []string{"16"},
		"uint32":  []string{"32"},
		"uint64":  []string{"64"},
		"uintptr": []string{"72"},
		"float32": []string{"0.1234"},
		"float64": []string{"1.2345"},
		"time":    []string{tm.Format(time.RFC3339)},
	}
	assert.Equal(t, expected, values)
}

func TestEncodeBasicPtr(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	tm := time.Unix(600, 0).UTC()

	s := basicPtr{
		String:  withStr("abc"),
		Bool:    withBool(true),
		Int:     withInt(12),
		Int8:    withInt8(int8(8)),
		Int16:   withInt16(int16(16)),
		Int32:   withInt32(int32(32)),
		Int64:   withInt64(int64(64)),
		Uint:    withUint(uint(24)),
		Uint8:   withUint8(uint8(8)),
		Uint16:  withUint16(uint16(16)),
		Uint32:  withUint32(uint32(32)),
		Uint64:  withUint64(uint64(64)),
		UinPtr:  withUintPtr(uintptr(72)),
		Float32: withFloat32(float32(0.1234)),
		Float64: withFloat64(1.2345),
		Time:    withTime(tm),
	}
	values, err := encoder.Values(s)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"string":  []string{"abc"},
		"bool":    []string{"true"},
		"int":     []string{"12"},
		"int8":    []string{"8"},
		"int16":   []string{"16"},
		"int32":   []string{"32"},
		"int64":   []string{"64"},
		"uint":    []string{"24"},
		"uint8":   []string{"8"},
		"uint16":  []string{"16"},
		"uint32":  []string{"32"},
		"uint64":  []string{"64"},
		"uintptr": []string{"72"},
		"float32": []string{"0.1234"},
		"float64": []string{"1.2345"},
		"time":    []string{tm.Format(time.RFC3339)},
	}
	assert.Equal(t, expected, values)
}

func TestZeroVal(t *testing.T)  {
	test := assert.New(t)
	encoder := NewEncoder()

	values, err := encoder.Values(basicVal{})
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"string":  []string{""},
		"bool":    []string{"false"},
		"int":     []string{"0"},
		"int8":    []string{"0"},
		"int16":   []string{"0"},
		"int32":   []string{"0"},
		"int64":   []string{"0"},
		"uint":    []string{"0"},
		"uint8":   []string{"0"},
		"uint16":  []string{"0"},
		"uint32":  []string{"0"},
		"uint64":  []string{"0"},
		"uintptr": []string{"0"},
		"float32": []string{"0"},
		"float64": []string{"0"},
		"time":    []string{time.Time{}.Format(time.RFC3339)},
	}
	assert.Equal(t, expected, values)
}

func TestZeroPtr(t *testing.T)  {
	test := assert.New(t)
	encoder := NewEncoder()

	values, err := encoder.Values(basicPtr{})
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"string":  []string{""},
		"bool":    []string{""},
		"int":     []string{""},
		"int8":    []string{""},
		"int16":   []string{""},
		"int32":   []string{""},
		"int64":   []string{""},
		"uint":    []string{""},
		"uint8":   []string{""},
		"uint16":  []string{""},
		"uint32":  []string{""},
		"uint64":  []string{""},
		"uintptr": []string{""},
		"float32": []string{""},
		"float64": []string{""},
		"time":    []string{""},
	}
	assert.Equal(t, expected, values)
}

func TestOmitZeroVal(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()
	values, err := encoder.Values(basicValWithOmit{})
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	assert.Equal(t, url.Values{}, values)
}

func TestOmitZeroPtr(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	values, err := encoder.Values(basicPtrWithOmit{})
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	assert.Equal(t, url.Values{}, values)
}

func TestTimeFormat(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	tm := time.Unix(600, 0).UTC()

	times := struct {
		Rfc3339    time.Time  `qs:"default_fmt"`
		Second     time.Time  `qs:"default_second,second"`
		Millis     time.Time  `qs:"default_millis,millis"`
		Rfc3339Ptr *time.Time `qs:"default_fmt_ptr"`
		SecondPtr  *time.Time `qs:"default_second_ptr,second"`
		MillisPtr  *time.Time `qs:"default_millis_ptr,millis"`
	}{
		tm,
		tm,
		tm,
		&tm,
		&tm,
		&tm,
	}
	values, err := encoder.Values(times)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"default_fmt":    []string{"1970-01-01T00:10:00Z"},
		"default_second": []string{"600"},
		"default_millis": []string{"600000"},
		"default_fmt_ptr":    []string{"1970-01-01T00:10:00Z"},
		"default_second_ptr": []string{"600"},
		"default_millis_ptr": []string{"600000"},
	}
	assert.Equal(t, expected, values)
}

func TestIgnoreEmptySlice(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	s := struct {
		A []string	`qs:"a"`
		B []string	`qs:"b"`
		C *[]string	`qs:"c"`
	}{
		A: nil,
		B: []string{},
		C: nil,
	}

	values, err := encoder.Values(s)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	test.Equal(url.Values{}, values)
}

func TestSliceValWithBasicVal(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	s := struct {
		StringList []string `qs:"str_list"`
		BoolList   []bool   `qs:"bool_list"`
		IntList    []int    `qs:"int_list"`
	}{
		StringList: []string{"", "a", "b", "c"},
		BoolList:   []bool{true, false},
		IntList:    []int{0, 1, 2, 3},
	}
	values, err := encoder.Values(s)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"str_list":  []string{"", "a", "b", "c"},
		"bool_list": []string{"true", "false"},
		"int_list":  []string{"0", "1", "2", "3"},
	}
	assert.Equal(t, expected, values)
}

func TestSliceValWithBasicPtr(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	s := struct {
		StringList []*string `qs:"str_list"`
		BoolList   []*bool   `qs:"bool_list"`
		IntList    []*int    `qs:"int_list"`
	}{
		StringList: []*string{withStr(""), withStr("a"), withStr("b"), withStr("c")},
		BoolList:   []*bool{withBool(true), withBool(false)},
		IntList:    []*int{withInt(0), withInt(1), withInt(2), withInt(3)},
	}
	values, err := encoder.Values(s)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"str_list":  []string{"", "a", "b", "c"},
		"bool_list": []string{"true", "false"},
		"int_list":  []string{"0", "1", "2", "3"},
	}
	assert.Equal(t, expected, values)
}

func TestSlicePtrWithBasicVal(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	strList := []string{"", "a", "b", "c"}
	boolList := []bool{true, false}
	intList := []int{0, 1, 2, 3}

	s := struct {
		StringList *[]string `qs:"str_list"`
		BoolList   *[]bool   `qs:"bool_list"`
		IntList    *[]int    `qs:"int_list"`
	}{
		StringList: &strList,
		BoolList:   &boolList,
		IntList:    &intList,
	}
	values, err := encoder.Values(s)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"str_list":  []string{"", "a", "b", "c"},
		"bool_list": []string{"true", "false"},
		"int_list":  []string{"0", "1", "2", "3"},
	}
	assert.Equal(t, expected, values)
}

func TestSlicePtrWithBasicPtr(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	strList := []*string{withStr(""), withStr("a"), withStr("b"), withStr("c")}
	boolList := []*bool{withBool(true), withBool(false)}
	intList := []*int{withInt(0), withInt(1), withInt(2), withInt(3)}

	s := struct {
		StringList *[]*string `qs:"str_list"`
		BoolList   *[]*bool   `qs:"bool_list"`
		IntList    *[]*int    `qs:"int_list"`
	}{
		StringList: &strList,
		BoolList:   &boolList,
		IntList:    &intList,
	}
	values, err := encoder.Values(s)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"str_list":  []string{"", "a", "b", "c"},
		"bool_list": []string{"true", "false"},
		"int_list":  []string{"0", "1", "2", "3"},
	}
	assert.Equal(t, expected, values)
}

func TestArrayFormat_Comma(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	s := struct {
		StringList []string `qs:"str_list,comma"`
	}{
		StringList: []string{"a", "b", "c"},
	}
	values, err := encoder.Values(s)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"str_list":  []string{"a,b,c"},
	}
	assert.Equal(t, expected, values)
}

func TestArrayFormat_Repeat(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	s := struct {
		StringList []string `qs:"str_list"`
	}{
		StringList: []string{"a", "b", "c"},
	}
	values, err := encoder.Values(s)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"str_list":  []string{"a", "b", "c"},
	}
	assert.Equal(t, expected, values)
}

func TestArrayFormat_Bracket(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	s := struct {
		StringList []string `qs:"str_list,bracket"`
	}{
		StringList: []string{"a", "b", "c"},
	}
	values, err := encoder.Values(s)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"str_list[]":  []string{"a", "b", "c"},
	}
	assert.Equal(t, expected, values)
}

func TestArrayFormat_Index(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	s := struct {
		StringList []string `qs:"str_list,index"`
	}{
		StringList: []string{"a", "b", "c"},
	}
	values, err := encoder.Values(s)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"str_list[0]":  []string{"a"},
		"str_list[1]":  []string{"b"},
		"str_list[2]":  []string{"c"},
	}
	assert.Equal(t, expected, values)
}

func TestNestedStruct(t *testing.T) {
	test := assert.New(t)
	encoder := NewEncoder()

	tm := time.Unix(600, 0)

	type Nested struct{
		Time time.Time	`qs:"time,second"`
	}

	s := struct {
		Nested Nested `qs:"nested"`
	}{
		Nested: Nested{
			tm,
		},
	}

	values, err := encoder.Values(s)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	expected := url.Values{
		"nested[time]":  []string{"600"},
	}
	assert.Equal(t, expected, values)
}

//------------------------------------------------

func withStr(v string) *string {
	return &v
}

func withBool(v bool) *bool {
	return &v
}

func withInt(v int) *int {
	return &v
}

func withInt8(v int8) *int8 {
	return &v
}

func withInt16(v int16) *int16 {
	return &v
}

func withInt32(v int32) *int32 {
	return &v
}

func withInt64(v int64) *int64 {
	return &v
}

func withUint(v uint) *uint {
	return &v
}

func withUint8(v uint8) *uint8 {
	return &v
}

func withUint16(v uint16) *uint16 {
	return &v
}

func withUint32(v uint32) *uint32 {
	return &v
}

func withUint64(v uint64) *uint64 {
	return &v
}

func withUintPtr(v uintptr) *uintptr {
	return &v
}

func withFloat32(v float32) *float32 {
	return &v
}

func withFloat64(v float64) *float64 {
	return &v
}

func withTime(v time.Time) *time.Time {
	return &v
}

/*
func TestGetTagNameAndOpts(t *testing.T) {
	test := assert.New(t)

	table := []struct{
		expectedName string
		expectedOpts []string
	} {
		{expectedName: "Age", expectedOpts: []string{}},
		{expectedName: "-", expectedOpts: []string{}},
		{expectedName: "msg", expectedOpts: []string{}},
		{expectedName: "time", expectedOpts: []string{"unixSecond", "omitempty"}},
	}

	st := reflect.TypeOf(struct {
		Age      int
		Verified bool      `qs:"-"`
		Message  string    `qs:"msg"`
		Time     time.Time `qs:"time,unixSecond,omitempty"`
	}{})
	for i := 0; i < st.NumField(); i++ {
		name, opts := getTagNameAndOpts(st.Field(i))
		test.Equal(table[i].expectedName, name)
		test.Equal(table[i].expectedOpts, opts)
	}
}

func TestReflectSlice(t *testing.T) {
	test := assert.New(t)

	values := make(url.Values)
	intSlice := []int{ 1, 2, 3, 4, 5 }
	reflectSliceAndArray(values, reflect.ValueOf(intSlice), "list", nil)

	test.Len(values, 1)

	expectedValues := url.Values{
		"list" : []string{ "1", "2", "3", "4", "5"},
	}

	test.Equal(expectedValues, values)
}

func TestFormatSliceWithBracket(t *testing.T) {
	test := assert.New(t)

	values := make(url.Values)
	intSlice := []int{ 1, 2, 3, 4, 5 }
	reflectSliceAndArray(values, reflect.ValueOf(intSlice), "list", []string{"bracket"})

	test.Len(values, 1)

	expectedValues := url.Values{
		"list[]" : []string{ "1", "2", "3", "4", "5"},
	}
	test.Equal(expectedValues, values)
}

func TestReflectSliceWithIndex(t *testing.T) {
	test := assert.New(t)

	values := make(url.Values)
	intSlice := []int{ 1, 2, 3, 4, 5 }
	reflectSliceAndArray(values, reflect.ValueOf(intSlice), "list", []string{"index"})

	test.Len(values, 5)

	expectedValues := url.Values{
		"list[0]" : []string{"1"},
		"list[1]" : []string{"2"},
		"list[2]" : []string{"3"},
		"list[3]" : []string{"4"},
		"list[4]" : []string{"5"},
	}
	test.Equal(expectedValues, values)
}*/

/*func BenchmarkValues(b *testing.B) {
	param := struct {
		//Price       decimal.Decimal `qs:"price,omitempty"`
		EpochSecond time.Time `qs:"start,unixSecond"`
		EpochMillis time.Time `qs:"end,unixMillis"`
		Bool        bool      `qs:"bool"`
		Int         int       `qs:"int"`
	}{
		//Price:       decimal.NewFromFloat(float64(5)),
		EpochSecond: time.Now(),
		EpochMillis: time.Now(),
		Bool:        true,
		Int:         12000,
	}
	b.ResetTimer()
	for i := 1; i <= b.N; i++ {
		_, _ = Values(param)
	}
}*/
