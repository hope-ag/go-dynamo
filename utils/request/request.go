package request

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func newResponse(data any, status int, success bool) *Response {
	payload := map[string]any{}
	if success {
		payload["success"] = true
	} else {
		payload["success"] = false
	}
	payload["data"] = data
	return &Response{
		status,
		payload,
	}
}

func (res *Response) bytes() []byte {
	data, _ := json.Marshal(res.Data)
	return data
}

func (res *Response) string() string {
	return string(res.bytes())
}

func (res *Response) sendResponse(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(res.Status)
	_, _ = w.Write(res.bytes())
	log.Println(res.string())
}

func SendSuccessResponse(w http.ResponseWriter, r *http.Request, data any) {
	var status int
	if data == nil {
		status = http.StatusNoContent
	} else {
		status = http.StatusOK
	}
	newResponse(data, status, true).sendResponse(w, r)
}

func SendSuccessResponseWithStatus(w http.ResponseWriter, r *http.Request, data any, status int) {
	newResponse(data, status, true).sendResponse(w, r)
}
func SendNoContentSuccessResponse(w http.ResponseWriter, r *http.Request) {
	newResponse(nil, http.StatusNoContent, true).sendResponse(w, r)
}

func SendErrorResponse(w http.ResponseWriter, r *http.Request, err error, status int) {
	var s int
	data := map[string]any{"error": err.Error()}
	if status == 0 {
		s = http.StatusInternalServerError
	} else {
		s = status
	}
	newResponse(data, s, false).sendResponse(w, r)
}

func SendBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]any{"error": err.Error()}
	newResponse(data, http.StatusBadRequest, false).sendResponse(w, r)
}

func SendInternalError(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]any{"error": err.Error()}
	newResponse(data, http.StatusInternalServerError, false).sendResponse(w, r)
}
