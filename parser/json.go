package parser

import (
	"encoding/json"
	"fmt"

	"github.com/soyacen/easyconfmgr"
)

type JsonParser struct {
	config map[string]interface{}
}

func (p *JsonParser) Parse(rawData []byte) error {
	config := make(map[string]interface{})
	err := json.Unmarshal(rawData, &config)
	if err != nil {
		return fmt.Errorf("failed parse json config, %w", err)
	}
	standardizedMap(config)
	p.config = config
	return nil
}

func (p *JsonParser) ConfigMap() map[string]interface{} {
	return p.config
}

func (p *JsonParser) Support(contentType string) bool {
	return JSON == contentType
}

func NewJsonParser() easyconfmgr.Parser {
	return &JsonParser{}
}
