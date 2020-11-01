package config

// MetricsConfig defines config for metrics
type MetricsConfig struct {
	minForSync           uint
	periodicSyncInterval uint
}

// GetMinForSync returns the minimum number of reqs / responses which can be pending before data is syncronized with DB
func (m MetricsConfig) GetMinForSync() uint {
	return m.minForSync
}

// GetPeriodicSyncInterval returns periodic sync interval (in seconds)
func (m MetricsConfig) GetPeriodicSyncInterval() uint {
	return m.periodicSyncInterval
}
