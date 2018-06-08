package utils

import (
	"gorage/src/config"
	"gorage/src/data"
)

// GetListWithStartAndEnd get the list with start and end
func GetListWithStartAndEnd(start int, end int) []data.KeyMap {
	var UUIDArray []data.KeyMap
	if start > len(config.KeyCacheArray)-1 || start >= end || start < 0 {
		return UUIDArray
	}
	if end > len(config.KeyCacheArray)-1 {
		end = len(config.KeyCacheArray) - 1
	}
	for i := start; i < end; i++ {
		UUIDArray = append(UUIDArray, config.KeyCacheArray[i])
	}
	return UUIDArray
}
