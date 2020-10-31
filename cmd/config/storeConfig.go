package config

// StoreConfig defines the config for store
type StoreConfig struct {
	fileName  string
	filePerms int
}

// GetFileName returns the file name to be used for storage
func (s *StoreConfig) GetFileName() string {
	return s.fileName
}

// GetFilePerms returns the permissions to be used for file
func (s *StoreConfig) GetFilePerms() int {
	return s.filePerms
}
