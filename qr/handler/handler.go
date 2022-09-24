package handler

import (
	"encoding/json"
	"net/http"
	ierr "qr/error"
	"time"
)

type ResponseBody struct {
	ErrorCode       int         `json:"errorCode"`
	ErrorMessage    string      `json:"errorMessage,omitempty"`
	ErrorHttpStatus int         `json:"errorHttpStatus,omitempty"`
	ErrorTrace      interface{} `json:"errorTrace,omitempty"`
	Payload         interface{} `json:"payload,omitempty"`
	Time            string      `json:"time"`
}

func Response(req *http.Request, resp http.ResponseWriter, payload interface{}, err error) {
	if err != nil {
		//Check internal error
		if ie, ok := ierr.InternalErrors[err]; ok {
			ResponseError(req, resp, ie.HttpStatus, ie.Message, ie.Code)
		} else {
			ResponseError(req, resp, http.StatusInternalServerError, err.Error(), -1)
		}
		return
	} else {
		ResponseOk(req, resp, payload)
		return
	}
}

func ResponseOk(req *http.Request, resp http.ResponseWriter, payload interface{}) {
	body := ResponseBody{
		Time:    calculateTimeRequest(req),
		Payload: payload,
	}

	outputJSON(req, resp, http.StatusOK, body)
}

func ResponseError(req *http.Request, resp http.ResponseWriter, status int, errMessage string, errCode int) {
	body := ResponseBody{
		ErrorCode:       errCode,
		ErrorMessage:    errMessage,
		ErrorHttpStatus: status,
		Time:            calculateTimeRequest(req),
		Payload:         nil,
	}

	outputJSON(req, resp, status, body)
}

func calculateTimeRequest(req *http.Request) string {
	timeReq, ok := req.Context().Value("time_request").(time.Time)
	if !ok {
		panic("error no context time_request")
	}
	return time.Since(timeReq).String()
}

func outputJSON(req *http.Request, resp http.ResponseWriter, status int, payload interface{}) {
	output, err := json.Marshal(payload)
	if err != nil {
		body := ResponseBody{
			ErrorCode:       100000,
			ErrorMessage:    err.Error(),
			ErrorHttpStatus: http.StatusInternalServerError,
			Time:            calculateTimeRequest(req),
			Payload:         nil,
		}

		output, _ = json.Marshal(body)
		_, _ = resp.Write(output)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = resp.Write(output)
	resp.WriteHeader(status)
	return
}
