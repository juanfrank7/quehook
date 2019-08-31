package handlers

// Request provides a generalization of CloudWatch and API Gateway events
type Request struct {
	Body       string            `json:"body"`
	HTTPMethod string            `json:"httpMethod"`
	Headers    map[string]string `json:"headers"`
	Source     string            `json:"source"`
	Year       int               `json:"year"`
	Month      int               `json:"month"`
	Day        int               `json:"day"`
	Hour       int               `json:"hour"`
}
