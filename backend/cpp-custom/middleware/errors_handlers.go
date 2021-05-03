package middleware

import (
	"cpp-custom/logger"
	"encoding/json"
	"net/http"
)

func err_handling(err error, msg string, w http.ResponseWriter) error {
	if err != nil {
		logger.Error.Println(err, ":", msg)
		err_res := error_response{
			Err_type: err.Error(),
			Message: "error: " + msg,
		}
		json.NewEncoder(w).Encode(err_res)
	}
	return err
}
