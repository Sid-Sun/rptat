package store

import (
	"github.com/sid-sun/rptat/cmd/config"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

// Store defines the store interface
type Store interface {
	Write(data []byte) error
	Read() ([]byte, error)
}

type jsonStore struct {
	lgr       *zap.Logger
	fileName  string
	filePerms int
}

// Write takes raw bytes and writes it to file
func (j *jsonStore) Write(data []byte) error {
	err := ioutil.WriteFile(j.fileName, data, os.FileMode(j.filePerms))
	if err != nil {
		j.lgr.Sugar().Errorf("[Store] [Write] [WriteFile] %v", err)
		return err
	}
	return nil
}

// Read reads raw bytes from file and returns the raw bytes and / or an error
func (j *jsonStore) Read() ([]byte, error) {
	file, err := os.Open(j.fileName)
	if err != nil {
		j.lgr.Sugar().Errorf("[Store] [Read] [Open] %v", err)
		return nil, err
	}

	defer func() {
		err = file.Close()
		if err != nil {
			j.lgr.Sugar().Errorf("[Store] [Read] [Close] %v", err)
		}
	}()

	return ioutil.ReadAll(file)
}

// NewStore returns a new store implementation
func NewStore(s *config.StoreConfig, lgr *zap.Logger) Store {
	return &jsonStore{
		lgr:       lgr,
		fileName:  s.GetFileName(),
		filePerms: s.GetFilePerms(),
	}
}
