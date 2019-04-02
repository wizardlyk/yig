package meta

import (
	"database/sql"
	"github.com/journeymidnight/yig/helper"
	"github.com/journeymidnight/yig/redis"
	"github.com/journeymidnight/yig/yigtracer"
	"time"
)

type CacheType int

const (
	NoCache CacheType = iota
	EnableCache
	SimpleCache
)

var cacheNames = [...]string{"NOCACHE", "EnableCache", "SimpleCache"}

type MetaCache interface {
	Get(table redis.RedisDatabase, key string,
		onCacheMiss func() (interface{}, error),
		unmarshaller func([]byte) (interface{}, error), willNeed bool) (value interface{}, err error)
	Remove(table redis.RedisDatabase, key string)
	GetCacheHitRatio() float64
}

type disabledMetaCache struct{}

type entry struct {
	table redis.RedisDatabase
	key   string
	value interface{}
}

func newMetaCache(myType CacheType) (m MetaCache) {

	helper.Logger.Printf(10, "Setting Up Metadata Cache: %s\n", cacheNames[int(myType)])
	if myType == SimpleCache {
		m := new(enabledSimpleMetaCache)
		m.Hit = 0
		m.Miss = 0
		return m
	}
	return &disabledMetaCache{}
}

func (m *disabledMetaCache) Get(table redis.RedisDatabase, key string,
	onCacheMiss func() (interface{}, error),
	unmarshaller func([]byte) (interface{}, error), willNeed bool) (value interface{}, err error) {

	return onCacheMiss()
}

func (m *disabledMetaCache) Remove(table redis.RedisDatabase, key string) {
	return
}

func (m *disabledMetaCache) GetCacheHitRatio() float64 {
	return -1
}

type enabledSimpleMetaCache struct {
	Hit  int64
	Miss int64
}

func (m *enabledSimpleMetaCache) Get(table redis.RedisDatabase, key string,
	onCacheMiss func() (interface{}, error),
	unmarshaller func([]byte) (interface{}, error), willNeed bool) (value interface{}, err error) {
	var tracer = yigtracer.New()
	var tracerLogger = helper.TracerLogger
	var startTime int64
	var finishTime int64
	var consumeTime int64

	helper.Logger.Println(10, "enabledSimpleMetaCache Get. table:", table, "key:", key)

	//redis span 开始
	spanRedis := tracer.StartSpan("redis")
	time.Sleep(time.Millisecond * 9)
	value, err = redis.Get(table, key, unmarshaller)
	spanRedis.Finish()
	//redis span 结束

	if err != nil {
		helper.Logger.Println(5, "enabledSimpleMetaCache Get err:", err, "table:", table, "key:", key)
	}
	if err == nil && value != nil {
		m.Hit = m.Hit + 1

		spans := tracer.FinishedSpans()
		redisSpan := spans[0]
		startTime = redisSpan.StartTime.UnixNano() / 1e6
		finishTime = redisSpan.FinishTime.UnixNano() / 1e6
		consumeTime = finishTime - startTime
		tracerLogger.Println(6, "redis耗时：", consumeTime, "ms")
		tracerLogger.Println(6, "cache---TracerID：", redisSpan.SpanContext.TraceID)

		return value, nil
	}

	//if redis doesn't have the entry
	if onCacheMiss != nil {
		//tidb span 开始
		spanTidb := tracer.StartSpan("tidb")
		time.Sleep(time.Millisecond * 13)

		value, err = onCacheMiss()

		spanTidb.Finish()
		//tidb span 结束
		spans := tracer.FinishedSpans()
		redisSpan := spans[0]
		tidbSpan := spans[1]

		startTime = redisSpan.StartTime.UnixNano() / 1e6
		finishTime = redisSpan.FinishTime.UnixNano() / 1e6
		consumeTime = finishTime - startTime
		tracerLogger.Println(6, "redis耗时：", consumeTime, "ms")

		startTime = tidbSpan.StartTime.UnixNano() / 1e6
		finishTime = tidbSpan.FinishTime.UnixNano() / 1e6
		consumeTime = finishTime - startTime
		tracerLogger.Println(6, "tidb耗时：", consumeTime, "ms")
		tracerLogger.Println(6, "cache---TracerID：", tidbSpan.SpanContext.TraceID)

		if err != nil {
			if err != sql.ErrNoRows {
				helper.ErrorIf(err, "exec onCacheMiss() err.")
			}
			return
		}

		if willNeed == true {
			err = redis.Set(table, key, value)
			if err != nil {
				helper.Logger.Println(5, "WARNING: redis is down!")
				//do nothing, even if redis is down.
			}
		}
		m.Miss = m.Miss + 1
		return value, nil
	}
	return nil, nil
}

func (m *enabledSimpleMetaCache) Remove(table redis.RedisDatabase, key string) {
	redis.Remove(table, key)
}

func (m *enabledSimpleMetaCache) GetCacheHitRatio() float64 {
	return float64(m.Hit) / float64(m.Hit+m.Miss)
}
