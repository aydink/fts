package main

const (
	INDEXED      = true
	NOT_INDEXED  = false
	STORED       = true
	NOT_STORED   = false
	ANALYZED     = true
	NOT_ANALYZED = false

	FIELD_TYPE_TEXT    uint32 = 1
	FIELD_TYPE_INT32   uint32 = 2
	FIELD_TYPE_INT64   uint32 = 3
	FIELD_TYPE_FLOAT32 uint32 = 4
	FIELD_TYPE_FLOAT64 uint32 = 5
)

type DocValue interface {
	Type() uint32
}

type Int32 int32
type Int64 int64
type Float32 float32
type Float64 float64

func (i Int32) Type() uint32 {
	return FIELD_TYPE_INT32
}

func (i Int64) Type() uint32 {
	return FIELD_TYPE_INT64
}

func (i Float32) Type() uint32 {
	return FIELD_TYPE_FLOAT32
}

func (i Float64) Type() uint32 {
	return FIELD_TYPE_FLOAT64
}

// type DocValue2 interface {
// 	Name() string
// 	Type() int
// 	Value() interface{}
// }

type TextField struct {
	fid       uint32
	name      string
	value     string
	indexed   bool
	stored    bool
	analyzed  bool
	boost     float32
	fieldType uint32
	tokens    []Token
}

func NewTextField(fieldName, filedValue string, indexed, stored, analyzed bool) TextField {
	f := TextField{}
	f.fieldType = FIELD_TYPE_TEXT
	f.name = fieldName
	f.value = filedValue
	f.indexed = indexed
	f.stored = stored
	f.analyzed = analyzed
	f.boost = 1.0
	f.tokens = make([]Token, 0)
	return f
}

func (f *TextField) Type() uint32 {
	return f.fieldType
}

func (f *TextField) SetBoost(b float32) {
	f.boost = b
}

func (f *TextField) Boost() float32 {
	return f.boost
}

func (f *TextField) Tokens() []Token {
	return f.tokens
}

func (f *TextField) SetTokens(tokens []Token) {
	f.tokens = tokens
}

func (f *TextField) Indexed() bool {
	return f.indexed
}

func (f *TextField) Stored() bool {
	return f.stored
}

func (f *TextField) Analyzed() bool {
	return f.analyzed
}

func (f *TextField) Name() string {
	return f.name
}

func (f *TextField) StringValue() string {
	return f.value
}

// type Int64Value struct {
// 	name      string
// 	value     int64
// 	fieldType int
// 	stored    bool
// }

// func NewInt64Value(fieldName string, filedValue int64, stored bool) Int64Value {
// 	f := Int64Value{}
// 	f.fieldType = FIELD_TYPE_INT64
// 	f.name = fieldName
// 	f.value = filedValue
// 	f.stored = stored
// 	return f
// }

// func (f *Int64Value) Name() string {
// 	return f.name
// }

// func (f *Int64Value) Type() int {
// 	return f.fieldType
// }

// func (f *Int64Value) Stored() bool {
// 	return f.stored
// }

// func (f *Int64Value) Value() int64 {
// 	return f.value
// }

// //-----------------------------------------------

// type Int32Value struct {
// 	name      string
// 	value     int32
// 	fieldType int
// 	stored    bool
// }

// func NewInt32Value(fieldName string, filedValue int32, stored bool) Int32Value {
// 	f := Int32Value{}
// 	f.fieldType = FIELD_TYPE_INT64
// 	f.name = fieldName
// 	f.value = filedValue
// 	f.stored = stored
// 	return f
// }

// func (f *Int32Value) Name() string {
// 	return f.name
// }

// func (f *Int32Value) Type() int {
// 	return f.fieldType
// }

// func (f *Int32Value) Stored() bool {
// 	return f.stored
// }

// func (f *Int32Value) Value() int32 {
// 	return f.value
// }

// //-----------------------------------------------

// type Float32Value struct {
// 	name      string
// 	value     float32
// 	fieldType int
// 	stored    bool
// }

// func NewFloat32Value(fieldName string, filedValue float32, stored bool) Float32Value {
// 	f := Float32Value{}
// 	f.fieldType = FIELD_TYPE_FLOAT32
// 	f.name = fieldName
// 	f.value = filedValue
// 	f.stored = stored
// 	return f
// }

// func (f *Float32Value) Name() string {
// 	return f.name
// }

// func (f *Float32Value) Type() int {
// 	return f.fieldType
// }

// func (f *Float32Value) Stored() bool {
// 	return f.stored
// }

// func (f *Float32Value) Value() float32 {
// 	return f.value
// }

// //-----------------------------------------------

// type Float64Value struct {
// 	name      string
// 	value     float64
// 	fieldType int
// 	stored    bool
// }

// func NewFloat64Value(fieldName string, filedValue float64, stored bool) Float64Value {
// 	f := Float64Value{}
// 	f.fieldType = FIELD_TYPE_FLOAT64
// 	f.name = fieldName
// 	f.value = filedValue
// 	f.stored = stored
// 	return f
// }

// func (f *Float64Value) Name() string {
// 	return f.name
// }

// func (f *Float64Value) Type() int {
// 	return f.fieldType
// }

// func (f *Float64Value) Stored() bool {
// 	return f.stored
// }

// func (f *Float64Value) Value() float64 {
// 	return f.value
// }
