package store

import (
	"github.com/sid-sun/rptat/cmd/config"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

type Store interface {
	Write(data []byte) error
	Read() ([]byte, error)
}

type jsonStore struct {
	lgr       *zap.Logger
	fileName  string
	filePerms int
}

func (j *jsonStore) Write(data []byte) error {
	err := ioutil.WriteFile(j.fileName, data, os.FileMode(j.filePerms))
	if err != nil {
		j.lgr.Sugar().Errorf("[Store] [Write] [WriteFile] %v", err)
		return err
	}
	return nil
}

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

func NewStore(s config.StoreConfig, lgr *zap.Logger) Store {
	return &jsonStore{
		lgr:       lgr,
		fileName:  s.GetFileName(),
		filePerms: s.GetFilePerms(),
	}
}
