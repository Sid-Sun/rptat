package config

// MetricsConfig defines config for metrics
type MetricsConfig struct {
	minForSync int
}

// GetMinForSync returns the minimum number of reqs / responses which can be pending before data is syncronized with DB
func (m MetricsConfig) GetMinForSync() int {
	return m.minForSync
}
