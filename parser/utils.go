package parser

import (
	"github.com/spf13/cast"

	"github.com/soyacen/easyconfmgr"
)

func Parsers() []easyconfmgr.Parser {
	return []easyconfmgr.Parser{
		NewYamlParser(),
		NewJsonParser(),
		NewTomlParser(),
	}
}

func Parser(contentType string) easyconfmgr.Parser {
	switch contentType {
	case YAML, YML:
		return NewYamlParser()
	case JSON:
		return NewJsonParser()
	case TOML:
		return NewTomlParser()
	default:
		return &NopParser{}
	}
}

// standardizedMap recursively standardized map to map[string]interface{}
func standardizedMap(m map[string]interface{}) {
	for key, val := range m {
		switch val.(type) {
		case map[interface{}]interface{}:
			val = cast.ToStringMap(val)
			standardizedMap(val.(map[string]interface{}))
		case map[string]interface{}:
			standardizedMap(val.(map[string]interface{}))
		case []interface{}:
			standardizedSlice(val.([]interface{}))
		}
		m[key] = val
	}
}

// standardizedSlice recursively standardized map to map[string]interface{}
func standardizedSlice(v []interface{}) {
	for i, item := range v {
		switch item.(type) {
		case map[interface{}]interface{}:
			item = cast.ToStringMap(item)
			standardizedMap(item.(map[string]interface{}))
		case map[string]interface{}:
			standardizedMap(item.(map[string]interface{}))
		}
		v[i] = item
	}
}
