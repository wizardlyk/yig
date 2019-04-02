package yigtracer

import (
	"github.com/opentracing/opentracing-go"
	"sync"
)

type YigTracer struct {
	sync.RWMutex
	finishedSpans []*YigSpan
	injectors     map[interface{}]Injector
	extractors    map[interface{}]Extractor
}

func New() *YigTracer {
	t := &YigTracer{
		finishedSpans: []*YigSpan{},
		injectors:     make(map[interface{}]Injector),
		extractors:    make(map[interface{}]Extractor),
	}

	// register default injectors/extractors
	textPropagator := new(TextMapPropagator)
	t.RegisterInjector(opentracing.TextMap, textPropagator)
	t.RegisterExtractor(opentracing.TextMap, textPropagator)

	httpPropagator := &TextMapPropagator{HTTPHeaders: true}
	t.RegisterInjector(opentracing.HTTPHeaders, httpPropagator)
	t.RegisterExtractor(opentracing.HTTPHeaders, httpPropagator)

	return t
}

func (t *YigTracer) FinishedSpans() []*YigSpan {
	t.RLock()
	defer t.RUnlock()
	spans := make([]*YigSpan, len(t.finishedSpans))
	copy(spans, t.finishedSpans)
	return spans
}

func (t *YigTracer) recordSpan(span *YigSpan) {
	t.Lock()
	defer t.Unlock()
	t.finishedSpans = append(t.finishedSpans, span)
}

func (TextMapPropagator) Inject(ctx YigSpanContext, carrier interface{}) error {
	panic("implement me")
}

func (TextMapPropagator) Extract(carrier interface{}) (YigSpanContext, error) {
	panic("implement me")
}

func (t *YigTracer) RegisterInjector(format interface{}, injector Injector) {
	t.injectors[format] = injector
}

func (t *YigTracer) RegisterExtractor(format interface{}, extractor Extractor) {
	t.extractors[format] = extractor
}

func (t *YigTracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	sso := opentracing.StartSpanOptions{}
	for _, o := range opts {
		o.Apply(&sso)
	}
	return newYigSpan(t, operationName, sso)
}

func (t *YigTracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	panic("implement me")
}

func (t *YigTracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	panic("implement me")
}

type Injector interface {
	Inject(ctx YigSpanContext, carrier interface{}) error
}

type Extractor interface {
	Extract(carrier interface{}) (YigSpanContext, error)
}

type TextMapPropagator struct {
	HTTPHeaders bool
}
