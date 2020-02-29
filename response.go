package main

import "encoding/json"

type Status struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type Item struct {
	Name    string `json:"name"`
	Archive string `json:"archive"`
	Path    string `json:"path,omitempty"`
	Hash	string `json:"sha256sum,omitempty"`
	Status
}

type Archive struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
	Status
}

type ItemsResponse struct {
	Status
	Items []Item `json:"items"`
}

type RootResponse struct {
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

func NewItem(name string, archive string, path string, err error) Item {
	return Item{
		Name: name,
		Archive: archive,
		Path: path,
		Status: NewStatus(err),
	}
}

func NewItemsResponse(items []Item) ItemsResponse {
	return ItemsResponse{
		NewStatus(nil),
		items,
	}
}

// This function marshals any struct and returns the string containing
// the resulting json.
func marshalResponse(r interface{}) (string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// This function returns a json string containing all the useful info for
// the home page.
func GetItemsResponse(items []Item) string {
	res := NewItemsResponse(items)
	str, err := marshalResponse(res)
	if err != nil {
		return GetStatusResponse(err)
	}
	return str
}

func GetRootResponse(a []Archive) string {
	res := RootResponse{
		NewStatus(nil),
		a,
	}
	str, err := marshalResponse(res)
	if err != nil {
		return GetStatusResponse(err)
	}
	return str
}

// This function returns a json string containing only the status of a request.
func GetStatusResponse(e error) string {
	s := NewStatus(e)
	str, err := marshalResponse(s)
	if err != nil {
		return GetStatusResponse(err)
	}
	return str
}
