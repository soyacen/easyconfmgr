package easyconfmgrparser

import (
	"fmt"

	"github.com/BurntSushi/toml"

	"github.com/soyacen/easyconfmgr"
)

type TomlParser struct {
	config map[string]interface{}
}

func (p *TomlParser) Parse(rawData []byte) error {
	config := make(map[string]interface{})
	err := toml.Unmarshal(rawData, &config)
	if err != nil {
		return fmt.Errorf("failed parse toml config, %w", err)
	}
	standardizedMap(config)
	p.config = config
	return nil
}

func (p *TomlParser) ConfigMap() map[string]interface{} {
	return p.config
}

func (p *TomlParser) Support(contentType string) bool {
	return TOML == contentType
}

func NewTomlParser() easyconfmgr.Parser {
	return &TomlParser{}
}
