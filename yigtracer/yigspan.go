package yigtracer

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"sync"
	"sync/atomic"
	"time"
)

type YigSpanContext struct {
	TraceID int
	SpanID  int
	Sampled bool
	Baggage map[string]string
}

type YigSpan struct {
	sync.RWMutex
	ParentID int

	OperationName string
	StartTime     time.Time
	FinishTime    time.Time

	SpanContext YigSpanContext
	tags        map[string]interface{}
	logs        []YigLogRecord
	tracer      *YigTracer
}

func (c YigSpanContext) ForeachBaggageItem(handler func(k, v string) bool) () {
	for k, v := range c.Baggage {
		if !handler(k, v) {
			break
		}
	}
}

func (s *YigSpan) Finish() {
	s.Lock()
	s.FinishTime = time.Now()
	s.Unlock()
	s.tracer.recordSpan(s)
}

func (s *YigSpan) FinishWithOptions(opt opentracing.FinishOptions) {
	s.Lock()
	s.FinishTime = time.Now()
	s.Unlock()

	time := opt.LogRecords[0].Timestamp
	field := opt.LogRecords[0].Fields
	fmt.Println(time, field)
	//......

}

func (s *YigSpan) Context() opentracing.SpanContext {
	s.Lock()
	defer s.Unlock()
	return &s.SpanContext
}

func (s *YigSpan) SetOperationName(operationName string) opentracing.Span {
	s.Lock()
	defer s.Unlock()
	s.OperationName = operationName
	return s
}

func (s *YigSpan) SetTag(key string, value interface{}) opentracing.Span {
	s.Lock()
	defer s.Unlock()
	if key == string(ext.SamplingPriority) {
		if v, ok := value.(uint16); ok {
			s.SpanContext.Sampled = v > 0
			return s
		}
		if v, ok := value.(int); ok {
			s.SpanContext.Sampled = v > 0
			return s
		}
	}
	s.tags[key] = value
	return s
}

func (s *YigSpan) LogFields(fields ...log.Field) {
	fmt.Println("LogFields")
	return
}

func (s *YigSpan) LogKV(keyValues ...interface{}) {
	fmt.Println("LogKV")
	return
}

func (s *YigSpan) SetBaggageItem(key, value string) opentracing.Span {
	fmt.Println("SetBaggageItem")
	//s.Lock()
	//defer s.Unlock()
	//s.SpanContext = s.SpanContext.WithBaggageItem(key, val)
	return s
}

func (s *YigSpan) BaggageItem(key string) string {
	fmt.Println("BaggageItem")
	return ""
}

func (s *YigSpan) Tracer() opentracing.Tracer {
	fmt.Println("Tracer")
	return s.tracer
}

func (s *YigSpan) LogEvent(event string) {
	fmt.Println("LogEvent")
}

func (s *YigSpan) LogEventWithPayload(event string, payload interface{}) {
	fmt.Println("LogEventWithPayload")
}

func (s *YigSpan) Log(data opentracing.LogData) {

}

func newYigSpan(t *YigTracer, name string, opts opentracing.StartSpanOptions) *YigSpan {
	tags := opts.Tags
	if tags == nil {
		tags = map[string]interface{}{}
	}
	traceID := nextYigID()
	parentID := int(0)
	var baggage map[string]string
	sampled := true
	if len(opts.References) > 0 {
		traceID = opts.References[0].ReferencedContext.(YigSpanContext).TraceID
		parentID = opts.References[0].ReferencedContext.(YigSpanContext).SpanID
		sampled = opts.References[0].ReferencedContext.(YigSpanContext).Sampled
		baggage = opts.References[0].ReferencedContext.(YigSpanContext).Baggage
	}
	spanContext := YigSpanContext{traceID, nextYigID(), sampled, baggage}
	startTime := opts.StartTime
	if startTime.IsZero() {
		startTime = time.Now()
	}
	return &YigSpan{
		ParentID:      parentID,
		OperationName: name,
		StartTime:     startTime,
		tags:          tags,
		logs:          []YigLogRecord{},
		SpanContext:   spanContext,

		tracer: t,
	}
}

var yigIDSource = uint32(42)

func nextYigID() int {
	return int(atomic.AddUint32(&yigIDSource, 1))
}

func (s *YigSpan) Tags() map[string]interface{} {
	s.RLock()
	defer s.RUnlock()
	tags := make(map[string]interface{})
	for k, v := range s.tags {
		tags[k] = v
	}
	return tags
}