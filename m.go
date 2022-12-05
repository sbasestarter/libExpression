package libexpression

import (
	"reflect"
	"sync"
)

var (
	opValueTypeMap = sync.Map{}

	tmLock sync.Mutex
	tm     map[string]reflect.Type
)

func UpdateOpValueTypeMap(key string, t reflect.Type) {
	opValueTypeMap.Store(key, t)

	tmLock.Lock()
	defer tmLock.Unlock()

	tm = nil
}

func getOpValueTypeM() map[string]reflect.Type {
	if tm != nil {
		return tm
	}

	tmLock.Lock()
	defer tmLock.Unlock()

	if tm != nil {
		return tm
	}

	newM := make(map[string]reflect.Type)

	opValueTypeMap.Range(func(key, value any) bool {
		if t, ok := value.(reflect.Type); ok {
			k, _ := key.(string)
			newM[k] = t
		}

		return true
	})

	tm = newM

	return tm
}
