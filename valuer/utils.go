package valuer

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/derekparker/trie"
	"github.com/spf13/cast"
)

func mapToTrie(keyPrefix string, config interface{}, tree *trie.Trie) {
	tree.Add(keyPrefix, config)
	switch m := config.(type) {
	case map[interface{}]interface{}:
		for subKey, conf := range m {
			key := getKey(keyPrefix, cast.ToString(subKey))
			mapToTrie(key, conf, tree)
		}
	case map[string]interface{}:
		for subKey, conf := range m {
			key := getKey(keyPrefix, subKey)
			mapToTrie(key, conf, tree)
		}
	case []interface{}:
		for i, val := range m {
			key := getKey(keyPrefix, strconv.Itoa(i))
			mapToTrie(key, val, tree)
		}
	}
}

func MapToTrie(config interface{}, tree *trie.Trie) {
	mapToTrie("", config, tree)
}

func getKey(keyPrefix string, subKey string) string {
	var key string
	if keyPrefix != "" {
		key = fmt.Sprintf("%s.%s", keyPrefix, subKey)
	} else {
		key = subKey
	}
	return key
}

// mergeMaps merges two string maps.
func mergeMaps(src, target map[string]interface{}) {
	for key, srcVal := range src {
		// if target not contains srcKey, add srcVal and continue
		targetVal, ok := target[key]
		if !ok {
			target[key] = srcVal
			continue
		}

		// if targetVal and srcVal is not same type, ignore this val, and log it
		srcValType := reflect.TypeOf(srcVal)
		targetValType := reflect.TypeOf(targetVal)
		if targetValType != nil && srcValType != targetValType {
			fmt.Printf(
				"merge map, srcValType != targetValType; key=%s, srcValType=%v, targetValType=%v, srcVal=%v, targetVal=%v",
				key, srcValType, targetValType, srcVal, targetVal)
			continue
		}

		// if targetVal is `map[string]interface{}`, merge sub string map.
		// else add srcVal and continue
		switch subTargetVal := targetVal.(type) {
		case map[string]interface{}:
			subSrvVal := srcVal.(map[string]interface{})
			mergeMaps(subSrvVal, subTargetVal)
		default:
			target[key] = srcVal
			continue
		}
	}
}
