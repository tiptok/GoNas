package global

import (
	"time"

	"github.com/tiptok/gotransfer/comm"
)

//ICache  缓存类
type ICache interface {
	GetCache(key string) interface{}
	AddCache(key string, value interface{}) error
	RemoveCache(key string) error
}

//CacheBase 缓存基类
type CacheBase struct {
	//Timer comm.TimerWork
	CacheValue comm.DataContext
	TimerTask  *comm.Task
}

func (cache *CacheBase) NewCache(cacheId string, span int, workFucn func(p interface{})) *CacheBase {
	cache.TimerTask = &comm.Task{
		Interval: time.Duration(span),
		TaskId:   cacheId,
		Work:     workFucn,
		Param:    cache.CacheValue,
	}
	return cache
}

func (cache *CacheBase) GetCache(key string) interface{} {
	return cache.CacheValue.Get(key)
}

func (cache *CacheBase) AddCache(key string, value interface{}) error {
	cache.CacheValue.Set(key, value)
	return nil
}

func (cache *CacheBase) RemoveCache(key string) error {
	cache.CacheValue.Delete(key)
	return nil
}

func (cache *CacheBase) ContaninKey(key string) bool {
	_, result := cache.CacheValue.GetOk(key)
	return result
}
