package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	app := App{Name: r.URL.Query().Get(":name")}
	app.Get()

	if app.Id == 0 {
		http.NotFound(w, r)
	}
}

func AppInfo(w http.ResponseWriter, r *http.Request) {
	app := App{Name: r.URL.Query().Get(":name")}
	app.Get()

	if app.Id == 0 {
		http.NotFound(w, r)
	} else {
		b, err := json.Marshal(app)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, bytes.NewBuffer(b).String())
	}
}

func CreateAppHandler(w http.ResponseWriter, r *http.Request) {
	var app App

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &app)
	if err != nil {
		panic(err)
	}

	err = app.Create()
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, "success")
}
