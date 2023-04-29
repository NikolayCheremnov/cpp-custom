package middleware

type CheckerResponse struct {
	Message         string `json:"message"`
	LexinatorErrors string `json:"lexical_errors"`
	CheckerErrors   string `json:"checker_errors"`
}

type analyze_response struct {
	Message          string `json:"message"`
	LexinatorErrors  string `json:"lexical_errors"`
	ParsenatorErrors string `json:"syntax_errors"`
}

type basic_response struct {
	Message string `json:"message,omitempty"`
}

type error_response struct {
	Err_type string `json:"err_type,omitempty"`
	Message  string `json:"message,omitempty"`
}
