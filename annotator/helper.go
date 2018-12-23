package annotator

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func logRequest(r *http.Request) {
	if r.Body == nil {
		return
	}

	body, _ := ioutil.ReadAll(r.Body)

	obj := make(map[string]interface{})
	if err := json.Unmarshal(body, &obj); err == nil {
		log.WithData(obj).Debugf("%s: %s", r.Method, r.RequestURI)
		return
	}

	arr := make([]interface{}, 0)
	if err := json.Unmarshal(body, &arr); err == nil {
		log.WithData(arr).Debugf("%s: %s", r.Method, r.RequestURI)
		return
	}

	log.WithData(string(body)).Debugf("%s: %s", r.Method, r.RequestURI)
}
