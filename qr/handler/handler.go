package handler

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	ierr "qr/error"
	"time"
)

type ResponseBody struct {
	ErrorCode       int         `json:"errorCode"`
	ErrorMessage    string      `json:"errorMessage,omitempty"`
	ErrorHttpStatus int         `json:"errorHttpStatus,omitempty"`
	ErrorTrace      interface{} `json:"errorTrace,omitempty"`
	Data            interface{} `json:"data,omitempty"`
	Time            string      `json:"time"`
	RequestID       string      `json:"request_id"`
}

func Post(r *mux.Router, path string, cusHandler func(resp http.ResponseWriter, req *http.Request) (payload interface{}, err error)) *mux.Route {
	return initRoute(r, path, http.MethodPost, cusHandler)
}
func Get(r *mux.Router, path string, cusHandler func(resp http.ResponseWriter, req *http.Request) (payload interface{}, err error)) *mux.Route {
	return initRoute(r, path, http.MethodGet, cusHandler)
}
func Put(r *mux.Router, path string, cusHandler func(resp http.ResponseWriter, req *http.Request) (payload interface{}, err error)) *mux.Route {
	return initRoute(r, path, http.MethodPut, cusHandler)
}
func Delete(r *mux.Router, path string, cusHandler func(resp http.ResponseWriter, req *http.Request) (payload interface{}, err error)) *mux.Route {
	return initRoute(r, path, http.MethodDelete, cusHandler)
}

func initRoute(r *mux.Router, path string, method string, cusHandler func(resp http.ResponseWriter, req *http.Request) (payload interface{}, err error)) *mux.Route {
	return r.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		result, err := cusHandler(writer, request)
		response(request, writer, result, err)
	}).Methods(method)
}

func response(req *http.Request, resp http.ResponseWriter, payload interface{}, err error) {
	if err != nil {
		//Check internal error
		if ie, ok := ierr.InternalErrors[err]; ok {
			responseError(req, resp, ie.HttpStatus, ie.Message, ie.Code)
		} else {
			responseError(req, resp, http.StatusInternalServerError, err.Error(), -1)
		}
		return
	} else {
		responseOk(req, resp, payload)
		return
	}
}

func responseOk(req *http.Request, resp http.ResponseWriter, payload interface{}) {
	body := ResponseBody{
		Time:      calculateTimeRequest(req),
		Data:      payload,
		RequestID: getRequestID(req),
	}

	outputJSON(req, resp, http.StatusOK, body)
}

func responseError(req *http.Request, resp http.ResponseWriter, status int, errMessage string, errCode int) {
	body := ResponseBody{
		ErrorCode:       errCode,
		ErrorMessage:    errMessage,
		ErrorHttpStatus: status,
		Time:            calculateTimeRequest(req),
		Data:            nil,
		RequestID:       getRequestID(req),
	}

	outputJSON(req, resp, status, body)
}

func outputJSON(req *http.Request, resp http.ResponseWriter, status int, payload interface{}) {
	output, err := json.Marshal(payload)
	if err != nil {
		log.Err(err).Send()
		body := ResponseBody{
			ErrorCode:       100000,
			ErrorMessage:    err.Error(),
			ErrorHttpStatus: http.StatusInternalServerError,
			Time:            calculateTimeRequest(req),
			Data:            nil,
		}

		output, _ = json.Marshal(body)
		_, _ = resp.Write(output)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(status)
	_, _ = resp.Write(output)
	return
}

func calculateTimeRequest(req *http.Request) string {
	timeReq, ok := req.Context().Value("time_request").(time.Time)
	if !ok {
		panic("error no context time_request")
	}

	diff := time.Since(timeReq).String()
	return diff
}

func getRequestID(req *http.Request) string {
	reId, ok := req.Context().Value("request_id").(string)
	if !ok {
		log.Err(errors.New("error no context request_id")).Send()
		reId = ""
	}

	return reId
}
