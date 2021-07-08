package easyconfmgrparser

type NopParser struct {
}

func (n *NopParser) Parse(rawData []byte) error {
	return nil
}

func (n *NopParser) ConfigMap() map[string]interface{} {
	return nil
}

func (p *NopParser) Support(contentType string) bool {
	return false
}
