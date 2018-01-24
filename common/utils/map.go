// AbstractMap
//
// This type provides convenient methods,
// you should use "MakeAbstractMap(interface{])" to intialize this object.
//
//	var sourceMap map[string]int8
//	abstractMap := MakeAbstractMap(sourceMap)
//
//	// Type conversion
//	abstractMap.ToType("string-type", int32(0))
//
//	// Batch process
//	abstractMap.BatchProcess(8, batchCallback, restCallback)
package utils

import (
	"fmt"
	"reflect"
)

// TODO 以下的部分, 考虑放到公共组件库
func KeysOfMap(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}

	return keys
}

// Wrapping type for any map
type AbstractMap struct {
	currentMap reflect.Value
	mapType    reflect.Type
}

// Initialized a map by any instance of map.
func MakeAbstractMap(sourceMap interface{}) *AbstractMap {
	valueOfMap := reflect.ValueOf(sourceMap)

	if valueOfMap.Kind() != reflect.Map {
		panic(fmt.Sprintf("The type of object is not map: [%T]", sourceMap))
	}

	return &AbstractMap{
		currentMap: valueOfMap,
		mapType:    valueOfMap.Type(),
	}
}

// Converts the type of key and value of a map to desired types.
//
// This method uses "ConvertToByReflect(reflect.Value, reflect.Type)" to convert
// the instances of key and value from source map to desired ones.
func (m *AbstractMap) ToType(keyType reflect.Type, elemType reflect.Type) interface{} {
	valueOfSourceMap := m.currentMap

	resultMap := reflect.MakeMap(reflect.MapOf(keyType, elemType))

	for _, key := range valueOfSourceMap.MapKeys() {
		sourceElem := valueOfSourceMap.MapIndex(key)

		targetKey := ConvertToByReflect(key, keyType)
		targetElem := ConvertToByReflect(sourceElem, elemType)

		resultMap.SetMapIndex(targetKey, targetElem)
	}

	return resultMap.Interface()
}

func (m *AbstractMap) ToTypeOfTarget(keyOfTarget interface{}, elemOfTarget interface{}) interface{} {
	return m.ToType(
		reflect.TypeOf(keyOfTarget),
		reflect.TypeOf(elemOfTarget),
	)
}

// Processes the map with desired size.
//
// The "batchProcessor" will accept the batch of map which fits the desired size.
// After that, the "restProcessor" will accept the batch of map which is LESS THAN the desired size.
func (m *AbstractMap) BatchProcess(
	batchSize int,
	batchProcessor func(batch interface{}), restProcessor func(rest interface{}),
) {
	resultMap := reflect.MakeMap(m.mapType)

	for _, key := range m.currentMap.MapKeys() {
		resultMap.SetMapIndex(key, m.currentMap.MapIndex(key))

		if resultMap.Len() == batchSize {
			batchProcessor(resultMap.Interface())
			resultMap = reflect.MakeMap(m.mapType)
		}
	}

	if resultMap.Len() > 0 {
		restProcessor(resultMap.Interface())
	}
}

// Processes the map with desired size.
//
// This method is simple version of "BatchProcess".
func (m *AbstractMap) SimpleBatchProcess(
	batchSize int,
	batchProcessor func(batch interface{}),
) {
	m.BatchProcess(batchSize, batchProcessor, batchProcessor)
}
