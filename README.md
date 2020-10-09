# qs #
[![Build](https://github.com/sonh/qs/workflows/build/badge.svg?branch=main)](https://github.com/sonh/qs/actions)
[![Codecov](https://codecov.io/gh/sonh/qs/branch/main/graph/badge.svg)](https://codecov.io/gh/sonh/qs)
[![GoReportCard](https://goreportcard.com/badge/github.com/sonh/qs)](https://goreportcard.com/report/github.com/sonh/qs)
[![Release](https://img.shields.io/github/release/sonh/qs.svg?color=brightgreen)](https://github.com/sonh/qs/releases/)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/sonh/qs)](https://pkg.go.dev/github.com/sonh/qs)
[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/sonh/qs/blob/main/LICENSE)

Package sonh/qs encodes structs into url.Values.

## Installation
```shell
go get github.com/sonh/qs
```

## Usage
```go
import (
    "github.com/sonh/qs"
)
```
Package qs exports `NewEncoder()` function to create an encoder. 

Encoder caches struct info to speed up encoding process, use a single instance is highly recommended. 

Use `WithTagAlias()` func to register custom tag alias (default is `qs`)
```go
encoder = qs.NewEncoder(
    qs.WithTagAlias("myTag"),
)
```

Encoder has `Values()` and `Encode()` functions to encode structs into `url.Values`.

### Supported data types:
- all basic types (`bool`, `uint`, `string`, `float64`,...)
- `struct`
- `slice`, `array`
- `pointer`
- `time.Time`   
- custom type

### Example
```go
type Query struct {
    Tags   []string  `qs:"tags"`
    Limit  int       `qs:"limit"`
    From   time.Time `qs:"from"`
    Active bool      `qs:"active,omitempty"`  //omit empty value
    Ignore float64   `qs:"-"`                 //ignore
}

query := &Query{
    Tags:   []string{"docker", "golang", "reactjs"},
    Limit:  24,
    From:   time.Unix(1580601600, 0).UTC(),
    Ignore: 0,
}

encoder := qs.NewEncoder()
values, err := encoder.Values(query)
if err != nil {
    // Handle error
}
fmt.Println(values.Encode()) //(unescaped) output: "from=2020-02-02T00:00:00Z&limit=24&tags=docker&tags=golang&tags=reactjs"
```
### Bool format
Use `int` option to encode bool to integer
```go
type Query struct {
    DefaultFmt bool `qs:"default_fmt"`
    IntFmt     bool `qs:"int_fmt,int"`
}

query := &Query{
    DefaultFmt: true, 
    IntFmt:     true,
}
values, _ := encoder.Values(query)
fmt.Println(values.Encode()) // (unescaped) output: "default_fmt=true&int_fmt=1"
```
### Time format
By default, package encodes time.Time values as RFC3339 format. 

Including the `"second"` or `"millis"` option to signal that the field should be encoded as second or millisecond.
```go
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
}

encoder := qs.NewEncoder()
values, _ := encoder.Values(query)
fmt.Println(values.Encode()) // (unescaped) output: "default_fmt=2020-02-02T00:00:00Z&millis_fmt=1580601600000&second_fmt=1580601600"
```

### Slice/Array Format
Slice and Array default to encoding into multiple URL values of the same value name.
```go
type Query struct {
    Tags []string `qs:"tags"`
}

values, _ := encoder.Values(&Query{Tags: []string{"foo","bar"}})
fmt.Println(values.Encode()) //(unescaped) output: "tags=foo&tags=bar"
```

Including the `comma` option to signal that the field should be encoded as a single comma-delimited value.
```go
type Query struct {
    Tags []string `qs:"tags,comma"`
}

values, _ := encoder.Values(&Query{Tags: []string{"foo","bar"}})
fmt.Println(values.Encode()) //(unescaped) output: "tags=foo,bar"
```

Including the `bracket` option to signal that the multiple URL values should have "[]" appended to the value name.
```go
type Query struct {
    Tags []string `qs:"tags,bracket"`
}

values, _ := encoder.Values(&Query{Tags: []string{"foo","bar"}})
fmt.Println(values.Encode()) //(unescaped) output: "tags[]=foo&tags[]=bar"
```

The `index` option will append an index number with brackets to value name.
```go
type Query struct {
    Tags []string `qs:"tags,index"`
}

values, _ := encoder.Values(&Query{Tags: []string{"foo","bar"}})
fmt.Println(values.Encode()) //(unescaped) output: "tags[0]=foo&tags[1]=bar"
```

### Nested structs
All nested structs are encoded including the parent value name with brackets for scoping.
```go
type User struct {
    Verified bool      `qs:"verified"`
    From     time.Time `qs:"from,millis"`
}

type Query struct {
    User User `qs:"user"`
}

query := Query{
    User: User{
        Verified: true,
        From: time.Now(),
    },
}
values, _ := encoder.Values(query)
fmt.Println(values.Encode()) //(unescaped) output: "user[from]=1601623397728&user[verified]=true"
```

### Limitation
- `map` is not supported yet
- `struct`, `slice/array` multi-level nesting are limited
- no decoder yet

_Will improve in future versions_ 

## License
Distributed under MIT License, please see license file in code for more details.
