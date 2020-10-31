package contract

// Request defines the request structure for metrics
type Request struct {
	Date string
	Path string
}

// Response defines the response structure for metrics
type Response struct {
	Date string
	Path string
	Code int
}
