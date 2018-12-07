package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Contect struct {
	Id string
}

func Getexecid(host string, port string, containerid string) string {
	log.SetFlags(log.Llongfile)
	client := &http.Client{}
	request, err := http.NewRequest("POST",
		fmt.Sprintf("http://%s:%s/containers/%s/exec", host, port, containerid),
		strings.NewReader("{\"Tty\": true, \"Cmd\": [\"/bin/sh\"], \"AttachStdin\": true, \"AttachStderr\": true, \"Privileged\": true, \"AttachStdout\": true}"),
	)
	if err != nil {
		log.Println(err)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
	v := &Contect{}
	err = json.Unmarshal(content, v)
	if err != nil {
		log.Println(err)
	}
	return v.Id
}

func Resizecontainer(host string, port string, execid string, width string, height string) {
	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%s:%s/exec/%s/resize?h=%s&w=%s", host, port, execid, width, height),
		nil,
	)
	if err != nil {
		log.Println(err)
	}
	request.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		log.Println(err)
	}
}
