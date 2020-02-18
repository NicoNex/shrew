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
	Archive string `json:"archive"`
	Ok 		bool   `json:"ok,omitempty"`
	Desc	string `json:"err_description,omitempty"`
}

type Archive struct {
	Name string `json:"name"`
	Files []string `json:"files"`
	Ok 		bool   `json:"ok,omitempty"`
	Desc	string `json:"err_description,omitempty"`
}

type UploadResponse struct {
	Response
	Item Item `json:"item"`
}

type ItemsResponse struct {
	Response
	Items []Item `json:"items"`
}

type HomeResponse struct {
	Response
	Archives []Archive `json:"archives"`
}

func newResponse(err error) Response {
	var ok = (err == nil)
	var errmsg string
	
	if !ok {
		errmsg = err.Error()
	}
	return Response{ok, errmsg}
}

func newArchive(name string, files []string) Archive {
	return Archive{
		Name: name,
		Files: files,
	}
}

func newArchiveErr(name string, err error) Archive {
	return Archive{
		name,
		[]string{},
		false,
		err.Error(),
	}
}

// func newItem(name, archive string, ok bool) Item {
// 	return Item{
// 		name,
// 		archive,
// 		ok,
// 	}
// }

// func newUploadResponse(name, path string, err error) UploadResponse {
// 	return UploadResponse{
// 		newResponse(err),
// 		newItem(name, path, err == nil),
// 	}
// }

func newItemsResponse(items []Item) ItemsResponse {
	return ItemsResponse{
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
// func GetUploadResponse(name, archive string, err error) string {
// 	res := newUploadResponse(name, archive, err)
// 	str, err := getResponseStr(res)
// 	if err != nil {
// 		return GetErrResponse(err)
// 	}
// 	return str
// }

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
func GetItemsResponse(items []Item) string {
	res := newItemsResponse(items)
	str, err := getResponseStr(res)
	if err != nil {
		return GetErrResponse(err)
	}
	return str
}

func GetHomeResponse(a []Archive) string {
	res := HomeResponse{
		newResponse(nil),
		a,
	}
	str, err := getResponseStr(res)
	if err != nil {
		return GetErrResponse(err)
	}
	return str
}
