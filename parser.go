package easyconfmgr

// Parser encode raw data to map `map[string]interface{}`
type Parser interface {
	// Parse unmarshal raw data into a config that type of `map[string]interface{}`
	// must ensure the sub value is type of `map[string]interface{}`
	Parse(rawData []byte) error
	// ConfigMap return config that type of `map[string]interface{}`
	ConfigMap() map[string]interface{}
	// Support returns true if the Parser supports this contentType, return false if the Parser not supports this contentType
	Support(contentType string) bool
}
