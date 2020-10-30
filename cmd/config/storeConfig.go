package config

type StoreConfig struct {
	fileName  string
	filePerms int
}

func (s *StoreConfig) GetFileName() string {
	return s.fileName
}

func (s *StoreConfig) GetFilePerms() int {
	return s.filePerms
}
