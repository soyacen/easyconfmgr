package easyconfmgr

// Loader load data from some medium.
type Loader interface {
	// ContentType get config content type
	ContentType() string
	// Load load config raw data from some where, like local file or remote server
	Load() error
	// RawData get raw data
	RawData() []byte
}
