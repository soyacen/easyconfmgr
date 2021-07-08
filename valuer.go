package easyconfmgr

import (
	"time"
)

// Valuer is Gets the value of a key, and Unmarshal value to a struct
type Valuer interface {
	// AddConfig add map[string]interface{} config
	AddConfig(configs ...map[string]interface{})
	// AllConfigs merges all config and returns them as a map[string]interface{}.
	AllConfigs() map[string]interface{}
	// UnmarshalKey takes a single key and unmarshal it into a Struct.
	UnmarshalKey(key string, rawVal interface{}) error
	// Unmarshal unmarshal the config into a Struct.
	Unmarshal(rawVal interface{}) error
	// Get can retrieve any value given the key to use.
	Get(key string) (interface{}, error)
	FloatValuer
	IntValuer
	BoolValuer
	StringValuer
	TimeValuer
	SliceValuer
	MapValuer
}

type FloatValuer interface {
	// GetFloat32 returns the value associated with the key as a float32.
	GetFloat32(key string) (float32, error)
	// GetFloat64 returns the value associated with the key as a float64.
	GetFloat64(key string) (float64, error)
}

type IntValuer interface {
	// GetInt returns the value associated with the key as an int.
	GetInt(key string) (int, error)
	// GetInt8 returns the value associated with the key as an int8.
	GetInt8(key string) (int8, error)
	// GetInt16 returns the value associated with the key as an int16.
	GetInt16(key string) (int16, error)
	// GetInt32 returns the value associated with the key as an int32.
	GetInt32(key string) (int32, error)
	// GetInt64 returns the value associated with the key as an int64.
	GetInt64(key string) (int64, error)

	// GetUint returns the value associated with the key as an uint.
	GetUint(key string) (uint, error)
	// GetUint8 returns the value associated with the key as an uint8.
	GetUint8(key string) (uint8, error)
	// GetUint16 returns the value associated with the key as an uint16.
	GetUint16(key string) (uint16, error)
	// GetUint32 returns the value associated with the key as an uint32.
	GetUint32(key string) (uint32, error)
	// GetUint64 returns the value associated with the key as an uint64.
	GetUint64(key string) (uint64, error)
}

type BoolValuer interface {
	// GetBool returns the value associated with the key as a bool.
	GetBool(key string) (bool, error)
}

type StringValuer interface {
	// GetString returns the value associated with the key as a string.
	GetString(key string) (string, error)
}

type TimeValuer interface {
	// GetTime returns the value associated with the key as time.
	GetTime(key string) (time.Time, error)
	// GetDuration returns the value associated with the key as a duration.
	GetDuration(key string) (time.Duration, error)
}

type SliceValuer interface {
	// GetSlice returns the value associated with the key as a slice of interface values.
	GetSlice(key string) ([]interface{}, error)
	// GetIntSlice returns the value associated with the key as a slice of int values.
	GetIntSlice(key string) ([]int, error)
	// GetBoolSlice returns the value associated with the key as a slice of bool values.
	GetBoolSlice(key string) ([]bool, error)
	// GetStringSlice returns the value associated with the key as a slice of strings.
	GetStringSlice(key string) ([]string, error)
	// GetDurationSlice returns the value associated with the key as a slice of duration.
	GetDurationSlice(key string) ([]time.Duration, error)
}

type MapValuer interface {
	// GetStringMap returns the value associated with the key as a map of interfaces.
	GetStringMap(key string) (map[string]interface{}, error)
	// GetStringMapInt64 returns the value associated with the key as a map of int64.
	GetStringMapInt64(key string) (map[string]int64, error)
	// GetStringMapInt returns the value associated with the key as a map of int.
	GetStringMapInt(key string) (map[string]int, error)
	// GetStringMapBool returns the value associated with the key as a map of bool.
	GetStringMapBool(key string) (map[string]bool, error)
	// GetStringMapString returns the value associated with the key as a map of strings.
	GetStringMapString(key string) (map[string]string, error)
	// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
	GetStringMapStringSlice(key string) (map[string][]string, error)
}
