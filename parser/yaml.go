package parser

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/soyacen/easyconfmgr"
)

type YamlParser struct {
	config map[string]interface{}
}

func (p *YamlParser) Parse(rawData []byte) error {
	config := make(map[string]interface{})
	err := yaml.Unmarshal(rawData, &config)
	if err != nil {
		return fmt.Errorf("failed parse yaml config, %w", err)
	}
	standardizedMap(config)
	p.config = config
	return nil
}

func (p *YamlParser) ConfigMap() map[string]interface{} {
	return p.config
}

func (p *YamlParser) Support(contentType string) bool {
	return YAML == contentType || YML == contentType
}

func NewYamlParser() easyconfmgr.Parser {
	return &YamlParser{}
}
