package mackenzie

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
	"time"
)

type cacheItem[T any] struct {
	item    T
	created time.Time
}

type Cache[T any] struct {
	store             map[string]*cacheItem[T]
	outType           reflect.Type
	inTypes           []reflect.Type
	caller            reflect.Value
	willError         bool
	config            Config
	lock              sync.RWMutex
	stopCleanInterval chan bool
}

type Config struct {
	Lifetime      time.Duration
	CleanInterval time.Duration
}

func Create[T any](caller any, cfg ...Config) (*Cache[T], error) {
	tType := reflect.ValueOf(new(T)).Elem().Type()
	callerVal := reflect.ValueOf(caller)
	if callerVal.Type().Kind() != reflect.Func {
		return nil, ErrCallerMustBeFunction
	}
	if callerVal.Type().NumIn() == 0 {
		return nil, ErrCallerMustHaveAtLeastOneArgument
	}
	if callerVal.Type().NumOut() == 0 {
		return nil, ErrCallerMustHaveAtLeastOneReturnValue
	}
	// ensure that there's only two return values
	if callerVal.Type().NumOut() > 2 {
		return nil, ErrCallerMustHaveNoMoreThanTwoReturnValues
	}
	// check if the first return value is equal to T
	if callerVal.Type().Out(0) != tType {
		return nil, ErrCallerMustReturnTAsItsFirstMethod
	}
	returnsError := false
	if callerVal.Type().NumOut() == 2 {
		returnsError = true
	}
	// check if the last return value is error
	if returnsError {
		if callerVal.Type().Out(1) != reflect.TypeOf(new(error)).Elem() {
			return nil, ErrCallerMustReturnAnErrorAsItsLastMethod
		}
	}

	c := &Cache[T]{
		store:     make(map[string]*cacheItem[T]),
		outType:   tType,
		willError: returnsError,
		inTypes: func() []reflect.Type {
			var types []reflect.Type
			for i := 0; i < callerVal.Type().NumIn(); i++ {
				types = append(types, callerVal.Type().In(i))
			}
			return types
		}(),
		caller: callerVal,
		config: func() Config {
			if len(cfg) > 0 {
				return cfg[0]
			}
			return Config{
				Lifetime: 60,
			}
		}(),
		stopCleanInterval: make(chan bool),
	}

	if c.config.CleanInterval > 0 {
		go func() {
			for {
				select {
				case <-time.Tick(c.config.CleanInterval):
					c.ClearExpired()
				case <-c.stopCleanInterval:
					return
				}
			}
		}()
	}
	return c, nil
}

func (c *Cache[T]) Unload() {
	c.stopCleanInterval <- true
}

func (c *Cache[T]) Get(in ...any) (T, error) {
	if len(in) != len(c.inTypes) {
		return *new(T), ErrIncorrectNumberOfArguments
	}
	// ensure the types are correct
	for i, v := range in {
		if reflect.TypeOf(v) != c.inTypes[i] {
			return *new(T), ErrIncorrectTypeForArgument(c.inTypes[i].String(), reflect.TypeOf(v).String())
		}
	}
	key := c.getKey(in)
	c.lock.RLock()
	ci, ok := c.store[key]
	c.lock.RUnlock()
	if ok {
		if c.config.Lifetime > 0 && time.Since(ci.created) > c.config.Lifetime {
			slog.Debug("cache expired")
			c.lock.Lock()
			delete(c.store, key)
			c.lock.Unlock()
		} else {
			slog.Debug("returning cached item")
			return ci.item, nil
		}
	}

	slog.Debug("calling")
	item := c.caller.Call(func() []reflect.Value {
		var values []reflect.Value
		for _, v := range in {
			values = append(values, reflect.ValueOf(v))
		}
		return values
	}())

	v := item[0].Interface().(T)
	var err error
	if c.willError && !item[1].IsNil() {
		err = errors.New(fmt.Sprintf("%v", item[1].Interface()))
		return v, err
	}
	c.lock.Lock()
	c.store[key] = &cacheItem[T]{
		item:    v,
		created: time.Now(),
	}
	c.lock.Unlock()
	return v, err
}

func (c *Cache[T]) getKey(in []any) string {
	var key string
	for _, v := range in {
		key += fmt.Sprintf("%#v", v)
	}
	slog.Debug("cache key", "key", key)
	return key
}

func (c *Cache[T]) ClearExpired() {
	for k, v := range c.store {
		if c.config.Lifetime > 0 && time.Since(v.created) > c.config.Lifetime {
			c.lock.Lock()
			delete(c.store, k)
			c.lock.Unlock()
		}
	}
}

func (c *Cache[T]) InvalidateAll() {
	c.lock.Lock()
	c.store = make(map[string]*cacheItem[T])
	c.lock.Unlock()
}

func (c *Cache[T]) Invalidate(in ...any) {
	key := c.getKey(in)
	c.lock.Lock()
	delete(c.store, key)
	c.lock.Unlock()
}

func (c *Cache[T]) ForceGet(in ...any) (T, error) {
	c.Invalidate(in...)
	return c.Get(in...)
}

func (c *Cache[T]) ExpiresIn(in ...any) time.Duration {
	key := c.getKey(in)
	c.lock.RLock()
	defer c.lock.RUnlock()
	if item, ok := c.store[key]; ok {
		return c.config.Lifetime - time.Since(item.created)
	}
	return 0
}

func IsErrCallInvalid(err error) bool {
	return errors.Is(err, ErrCallInvalid)
}

func IsErrMackenzie(err error) bool {
	return errors.Is(err, ErrMackenzie)
}
