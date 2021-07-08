package easyconfmgrvaluer

import (
	"fmt"
	"time"

	"github.com/derekparker/trie"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"

	"github.com/soyacen/easyconfmgr"
)

type TrieValuer struct {
	config map[string]interface{}
	tree   *trie.Trie
}

// =================== Valuer ===================

func (v *TrieValuer) AddConfig(configs ...map[string]interface{}) {
	for _, config := range configs {
		mergeMaps(config, v.config)
	}
	MapToTrie(v.config, v.tree)
}

func (v *TrieValuer) AllConfigs() map[string]interface{} {
	return v.config
}

func (v *TrieValuer) UnmarshalKey(key string, rawVal interface{}) error {
	input, err := v.Get(key)
	if err != nil {
		return err
	}
	return decode(input, defaultDecoderConfig(rawVal))
}

func (v *TrieValuer) Unmarshal(rawVal interface{}) error {
	return v.UnmarshalKey("", rawVal)
}

func (v *TrieValuer) Get(key string) (interface{}, error) {
	node, ok := v.tree.Find(key)
	if !ok {
		return nil, fmt.Errorf("not found %s value", key)
	}
	return node.Meta(), nil
}

// =================== FloatValuer ===================

func (v *TrieValuer) GetFloat32(key string) (float32, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToFloat32E(str)
}

func (v *TrieValuer) GetFloat64(key string) (float64, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToFloat64E(str)
}

// =================== IntValuer ===================

func (v *TrieValuer) GetInt(key string) (int, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToIntE(str)
}

func (v *TrieValuer) GetInt8(key string) (int8, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToInt8E(str)
}

func (v *TrieValuer) GetInt16(key string) (int16, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToInt16E(str)
}

func (v *TrieValuer) GetInt32(key string) (int32, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToInt32E(str)
}

func (v *TrieValuer) GetInt64(key string) (int64, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToInt64E(str)
}

func (v *TrieValuer) GetUint(key string) (uint, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToUintE(str)
}

func (v *TrieValuer) GetUint8(key string) (uint8, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToUint8E(str)
}

func (v *TrieValuer) GetUint16(key string) (uint16, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToUint16E(str)
}

func (v *TrieValuer) GetUint32(key string) (uint32, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToUint32E(str)
}

func (v *TrieValuer) GetUint64(key string) (uint64, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToUint64E(str)
}

// =================== BoolValuer ===================

func (v *TrieValuer) GetBool(key string) (bool, error) {
	str, err := v.Get(key)
	if err != nil {
		return false, err
	}
	return cast.ToBoolE(str)
}

// =================== StringValuer ===================

func (v *TrieValuer) GetString(key string) (string, error) {
	str, err := v.Get(key)
	if err != nil {
		return "", err
	}
	return cast.ToStringE(str)
}

// =================== TimeValuer ===================

func (v *TrieValuer) GetTime(key string) (time.Time, error) {
	str, err := v.Get(key)
	if err != nil {
		return time.Time{}, err
	}
	return cast.ToTimeE(str)
}

func (v *TrieValuer) GetDuration(key string) (time.Duration, error) {
	str, err := v.Get(key)
	if err != nil {
		return 0, err
	}
	return cast.ToDurationE(str)
}

// =================== SliceValuer ===================

func (v *TrieValuer) GetSlice(key string) ([]interface{}, error) {
	str, err := v.Get(key)
	if err != nil {
		return nil, err
	}
	return cast.ToSliceE(str)
}

func (v *TrieValuer) GetIntSlice(key string) ([]int, error) {
	str, err := v.Get(key)
	if err != nil {
		return nil, err
	}
	return cast.ToIntSliceE(str)
}

func (v *TrieValuer) GetBoolSlice(key string) ([]bool, error) {
	str, err := v.Get(key)
	if err != nil {
		return nil, err
	}
	return cast.ToBoolSliceE(str)
}

func (v *TrieValuer) GetStringSlice(key string) ([]string, error) {
	str, err := v.Get(key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringSliceE(str)
}

func (v *TrieValuer) GetDurationSlice(key string) ([]time.Duration, error) {
	str, err := v.Get(key)
	if err != nil {
		return nil, err
	}
	return cast.ToDurationSliceE(str)
}

// =================== SliceValuer ===================

func (v *TrieValuer) GetStringMap(key string) (map[string]interface{}, error) {
	str, err := v.Get(key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapE(str)
}

func (v *TrieValuer) GetStringMapInt64(key string) (map[string]int64, error) {
	str, err := v.Get(key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapInt64E(str)
}

func (v *TrieValuer) GetStringMapInt(key string) (map[string]int, error) {
	str, err := v.Get(key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapIntE(str)
}

func (v *TrieValuer) GetStringMapBool(key string) (map[string]bool, error) {
	str, err := v.Get(key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapBoolE(str)
}

func (v *TrieValuer) GetStringMapString(key string) (map[string]string, error) {
	str, err := v.Get(key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapStringE(str)
}

func (v *TrieValuer) GetStringMapStringSlice(key string) (map[string][]string, error) {
	str, err := v.Get(key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapStringSliceE(str)
}

// defaultDecoderConfig returns default mapsstructure.DecoderConfig with suppot
// of time.Duration values & string slices
func defaultDecoderConfig(output interface{}) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToIPHookFunc(),
			mapstructure.StringToIPNetHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.TextUnmarshallerHookFunc(),
		),
	}
	return c
}

// A wrapper around mapstructure.Decode that mimics the WeakDecode functionality
func decode(input interface{}, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func NewTrieTreeValuer() easyconfmgr.Valuer {
	return &TrieValuer{
		config: make(map[string]interface{}),
		tree:   trie.New(),
	}
}
