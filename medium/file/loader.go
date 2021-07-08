package easyconfmgrfile

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/soyacen/easyconfmgr"
)

type Loader struct {
	filename    string
	contentType string
	rawData     []byte
	log         easyconfmgr.Logger
}

func (l *Loader) ContentType() string {
	return l.contentType
}

func (l *Loader) Load() error {
	l.log.Info("reading file:", l.filename)
	rawData, err := ioutil.ReadFile(l.filename)
	if err != nil {
		return err
	}
	l.log.Debug("file content:", string(rawData))
	l.rawData = rawData
	return nil
}

func (l *Loader) RawData() []byte {
	return l.rawData
}

func NewLoader(filename string, log easyconfmgr.Logger) easyconfmgr.Loader {
	contentType := filepath.Ext(filename)
	if strings.HasPrefix(contentType, ".") {
		contentType = contentType[1:]
	}
	return &Loader{filename: filename, contentType: contentType, log: log}
}
