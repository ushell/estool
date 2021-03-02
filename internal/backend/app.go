package backend

import (
	"encoding/json"
	"estool/internal/elasticsearch"
	"fmt"
)

type App struct {}

type Form struct {
	Url string `json:"url"`
	AuthUser string `json:"auth_user,omitempty"`
	AuthPassword string `json:"auth_password,omitempty"`
	Query elasticsearch.Query `json:"query"`
	Filename string `json:"filename,omitempty"`
}

type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

func NewApp() *App {
	return &App{}
}

func (*App) Search(obj map[string]interface{}) Response {
	f, err := parseForm(obj)
	if err != nil {
		return Response{Code:1, Message: err.Error()}
	}

	err = elasticsearch.Init(f.Url, f.AuthUser, f.AuthPassword)
	if err != nil {
		return Response{Code:-1, Message: err.Error()}
	}

	r, err := elasticsearch.Search(f.Query)
	if err != nil {
		return Response{Code:1, Message: err.Error()}
	}

	return Response{Code: 0, Message: "ok", Data: r}
}

func (*App) Download(obj map[string]interface{}) Response {
	f, err := parseForm(obj)
	if err != nil {
		return Response{Code:-1, Message: err.Error()}
	}

	err = elasticsearch.Init(f.Url, f.AuthUser, f.AuthPassword)
	if err != nil {
		return Response{Code:1, Message: err.Error()}
	}

	r, err := elasticsearch.Download(f.Query)
	if err != nil {
		return Response{Code:1, Message: err.Error()}
	}

	return Response{Code: 0, Message: "ok", Data: r}
}

func parseForm(obj map[string]interface{}) (Form, error) {
	var f Form

	objStr, err := json.Marshal(obj)
	if err != nil {
		return f, err
	}

	fmt.Println("[*] ", string(objStr))

	err = json.Unmarshal(objStr, &f)
	if err != nil {
		return f, err
	}

	return f, nil
}