package yigtracer

import (
	"fmt"
	"reflect"
	"time"

	"github.com/opentracing/opentracing-go/log"
)

// MockLogRecord represents data logged to a Span via Span.LogFields or
// Span.LogKV.
type YigLogRecord struct {
	Timestamp time.Time
	Fields    []YigKeyValue
}

// MockKeyValue represents a single key:value pair.
type YigKeyValue struct {
	Key string

	// All MockLogRecord values are coerced to strings via fmt.Sprint(), though
	// we retain their type separately.
	ValueKind   reflect.Kind
	ValueString string
}

// EmitString belongs to the log.Encoder interface
func (m *YigKeyValue) EmitString(key, value string) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// EmitBool belongs to the log.Encoder interface
func (m *YigKeyValue) EmitBool(key string, value bool) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// EmitInt belongs to the log.Encoder interface
func (m *YigKeyValue) EmitInt(key string, value int) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// EmitInt32 belongs to the log.Encoder interface
func (m *YigKeyValue) EmitInt32(key string, value int32) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// EmitInt64 belongs to the log.Encoder interface
func (m *YigKeyValue) EmitInt64(key string, value int64) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// EmitUint32 belongs to the log.Encoder interface
func (m *YigKeyValue) EmitUint32(key string, value uint32) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// EmitUint64 belongs to the log.Encoder interface
func (m *YigKeyValue) EmitUint64(key string, value uint64) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// EmitFloat32 belongs to the log.Encoder interface
func (m *YigKeyValue) EmitFloat32(key string, value float32) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// EmitFloat64 belongs to the log.Encoder interface
func (m *YigKeyValue) EmitFloat64(key string, value float64) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// EmitObject belongs to the log.Encoder interface
func (m *YigKeyValue) EmitObject(key string, value interface{}) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// EmitLazyLogger belongs to the log.Encoder interface
func (m *YigKeyValue) EmitLazyLogger(value log.LazyLogger) {
	var meta YigKeyValue
	value(&meta)
	m.Key = meta.Key
	m.ValueKind = meta.ValueKind
	m.ValueString = meta.ValueString
}
