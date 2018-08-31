package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func PostRegist(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"fail", "description":"Invalid method"}`))
		return
	}

	if r.ContentLength == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"fail", "description":"Invalid Content-Length"}`))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"fail", "description":"Cannot read body"}`))
		return
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"fail", "description":"Cannot parse to json"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v\n", jsonBody)
	fmt.Printf("%v\n", jsonBody)

	_, _ = CreateSession(jsonBody["Name"].(string))

	s := GetSession()
	fmt.Println(s)
}
