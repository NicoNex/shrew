package main

import (
	"encoding/json"
)

type Response struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type Item struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Path    string `json:"path"`
}

type UploadResponse struct {
	Response
	Item Item `json:"item"`
}

type HomeResponse struct {
	Response
	Items []Item `json:"items"`
}

func newResponse(err error) Response {
	var ok = (err == nil)
	var errmsg string
	
	if !ok {
		errmsg = err.Error()
	}
	return Response{ok, errmsg}
}

func newItem(name, version, path string) Item {
	return Item{
		name,
		version,
		path,
	}
}

func newUploadResponse(name, version, path string, err error) UploadResponse {
	return UploadResponse{
		newResponse(err),
		newItem(name, version, path),
	}
}

func newHomeResponse(items []Item) HomeResponse {
	return HomeResponse{
		newResponse(nil),
		items,
	}
}

// This function marshals any struct and returns the string containing the resulting json.
func getResponseStr(r interface{}) (string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// This function returns the string containing the response json resulting after a file upload attempt.
func GetUploadResponse(name, version, path string, err error) string {
	res := newUploadResponse(name, version, path, err)
	str, err := getResponseStr(res)
	if err != nil {
		return GetErrResponse(err)
	}
	return str
}

// This function returns a json string containing an error message.
func GetErrResponse(err error) string {
	res := newResponse(err)
	str, err := getResponseStr(res)
	if err != nil {
		return GetErrResponse(err)
	}
	return str
}

// This function returns a json string containing all the useful info for the home page.
func GetHomeResponse(items []Item) string {
	res := newHomeResponse(items)
	str, err := getResponseStr(res)
	if err != nil {
		return GetErrResponse(err)
	}
	return str
}
