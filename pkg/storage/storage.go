package storage

import (
	"path"
	"strings"
)

var (
	memory map[string]interface{} = make(map[string]interface{})
)

func Get(key string) (interface{}, bool) {
	content, isPresent := memory[key]
	return content, isPresent
}

func GetCollection(key string) map[string]interface{} {

	var itemsInCollection = make(map[string]interface{})

	for location, item := range memory {
		if strings.HasPrefix(location, key) {
			var subLocation = strings.TrimPrefix(location, key+"/")
			if path.Base(subLocation) == subLocation {
				itemsInCollection[location] = item
			}
		}
	}

	return itemsInCollection
}

func Store(key string, content interface{}) {
	memory[key] = content
}

func Remove(key string) bool {

	_, prs := memory[key]
	if prs {
		delete(memory, key)
	}
	return prs
}
