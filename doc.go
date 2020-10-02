// Copyright 2020 Son Huynh. All rights reserved.

/*
Package qs encodes structs into url.Values.

Package exports `NewEncoder()` function to create an encoder.
Use `WithTagAlias()` func to register custom tag alias (default is `qs`)

	encoder = qs.NewEncoder(
		qs.WithTagAlias("myTag"),
	)

Encoder has `.Values()` and `Encode()` functions to encode structs into url.Values.

Supported data types:
	- all basic types (`bool`, `uint`, `string`, `float64`,...)
	- struct
	- slice/array
	- pointer
	- time.Time
	- custom type

Example

	type Query struct {
		Tags   []string  `qs:"tags"`
		Limit  int       `qs:"limit"`
		From   time.Time `qs:"from"`
		Active bool      `qs:"active,omitempty"`
		Ignore float64   `qs:"-"` //ignore
	}

	query := &Query{
		Tags:   []string{"docker", "golang", "reactjs"},
		Limit:  24,
		From:   time.Unix(1580601600, 0).UTC(),
		Ignore: 0,
	}

	encoder = qs.NewEncoder()

	values, err := encoder.Values(query)
	if err != nil {
		// Handle error
	}
	fmt.Println(values.Encode()) //(unescaped) output: "from=2020-02-02T00:00:00Z&limit=24&tags=docker&tags=golang&tags=reactjs"

Ignoring Fields

	type Struct struct {
        Field string `form:"-"` //using `-` to to tell qs to ignore fields
    }

Omitempty

	type Struct struct {
		Field1 string `form:",omitempty"` 		//using `omitempty` to to tell qs to omit empty field
		Field2 *int `form:"field2,omitempty"`
	}

By default, package encodes time.Time values as RFC3339 format.

Including the `"second"` or `"millis"` option to signal that the field should be encoded as second or millisecond.

	type Query struct {
		Default time.Time   `qs:"default_fmt"`
		Second  time.Time   `qs:"second_fmt,second"` //use `second` option
		Millis  time.Time   `qs:"millis_fmt,millis"` //use `millis` option
	}

	t := time.Unix(1580601600, 0).UTC()
	query := &Query{
		Default: t,
		Second:  t,
		Millis:  t,
		Decimal: decimal.NewFromFloat(0.012147483648),
	}

	encoder = qs.NewEncoder()
	values, _ := encoder.Values(query)
	fmt.Println(values.Encode()) // (unescaped) output: "default_fmt=2020-02-02T00:00:00Z&millis_fmt=1580601600000&second_fmt=1580601600"

Slice and Array default to encoding into multiple URL values of the same value name.

	type Query struct {
		Tags []string `qs:"tags"`
	}

	values, _ := encoder.Values(&Query{Tags: []string{"foo","bar"}})
	fmt.Println(values.Encode()) //(unescaped) output: "tags=foo&tags=bar"

Including the `comma` option to signal that the field should be encoded as a single comma-delimited value.

	type Query struct {
		Tags []string `qs:"tags,comma"`
	}

	values, _ := encoder.Values(&Query{Tags: []string{"foo","bar"}})
	fmt.Println(values.Encode()) //(unescaped) output: "tags=foo,bar"

Including the `bracket` option to signal that the multiple URL values should have "[]" appended to the value name.

	type Query struct {
		Tags []string `qs:"tags,bracket"`
	}

	values, _ := encoder.Values(&Query{Tags: []string{"foo","bar"}})
	fmt.Println(values.Encode()) //(unescaped) output: "tags[]=foo&tags[]=bar"


The `index` option will append an index number with brackets to value name

	type Query struct {
		Tags []string `qs:"tags,index"`
	}

	values, _ := encoder.Values(&Query{Tags: []string{"foo","bar"}})
	fmt.Println(values.Encode()) //(unescaped) output: "tags[0]=foo&tags[1]=bar"


All nested structs are encoded including the parent value name with brackets for scoping.

	type User struct {
		Verified bool      `qs:"verified"`
		From     time.Time `qs:"from,millis"`
	}

	type Query struct {
		User User `qs:"user"`
	}

	querys := Query{
		User: User{
			Verified: true,
			From: time.Now(),
		},
	}
	values, _ := encoder.Values(querys)
	fmt.Println(values.Encode()) //(unescaped) output: "user[from]=1601623397728&user[verified]=true"

Limitation
	- `interface`\, `[]interface`\, `map` are not supported yet
	- `struct`, `slice`/`array` multi-level nesting are limited
	- no decoder yet
_Will improve in future versions_
*/
package qs
