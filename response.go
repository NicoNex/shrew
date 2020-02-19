package main

import "encoding/json"

type Status struct {
	Ok    bool   `json:"ok,omitempty"`
	Error string `json:"error,omitempty"`
}

type Item struct {
	Name    string `json:"name"`
	Archive string `json:"archive"`
	Status
}

type Archive struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
	Status
}

type UploadResponse struct {
	Status
	Items []Item `json:"items"`
}

type HomeResponse struct {
	Status
	Archives []Archive `json:"archives"`
}

func NewStatus(err error) Status {
	var ok = (err == nil)
	var errmsg string

	if !ok {
		errmsg = err.Error()
	}
	return Status{ok, errmsg}
}

func NewArchive(name string, files []string) Archive {
	return Archive{
		name,
		files,
		NewStatus(nil),
	}
}

func NewArchiveErr(name string, err error) Archive {
	return Archive{
		name,
		[]string{},
		NewStatus(err),
	}
}

func NewItem(name, archive string, ok bool) Item {
	return Item{
		name,
		archive,
		NewStatus(nil),
	}
}

func NewUploadResponse(items []Item) UploadResponse {
	return UploadResponse{
		NewStatus(nil),
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

// This function returns a json string containing an error message.
func GetErrResponse(err error) string {
	res := NewStatus(err)
	str, err := getResponseStr(res)
	if err != nil {
		return GetErrResponse(err)
	}
	return str
}

// This function returns a json string containing all the useful info for the home page.
func GetUploadResponse(items []Item) string {
	res := NewUploadResponse(items)
	str, err := getResponseStr(res)
	if err != nil {
		return GetErrResponse(err)
	}
	return str
}

func GetHomeResponse(a []Archive) string {
	res := HomeResponse{
		NewStatus(nil),
		a,
	}
	str, err := getResponseStr(res)
	if err != nil {
		return GetErrResponse(err)
	}
	return str
}
