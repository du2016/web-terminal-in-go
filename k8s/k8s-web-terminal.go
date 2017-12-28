package main

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	remotecommandconsts "k8s.io/apimachinery/pkg/util/remotecommand"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	//"strings"
	"io"
	//"fmt"
)

func Handler(r io.Reader,w io.Writer,containername string,podname string,namespace string) {
	log.SetFlags(log.Llongfile)
	config, err := clientcmd.BuildConfigFromFlags("", "./config")
	if err != nil {
		log.Fatalln(err)
	}
	groupversion := schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}
	config.GroupVersion = &groupversion
	config.APIPath = "/api"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: api.Codecs}
	restclient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatalln(err)
	}

	req := restclient.Post().
		Resource("pods").
		Name(podname).
		Namespace(namespace).
		SubResource("exec").
		Param("container", containername).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("command", "/bin/bash").Param("tty", "true")

	req.VersionedParams(
		&api.PodExecOptions{
			Container: containername,
			Command:   []string{"sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		},
		api.ParameterCodec,
	)
	log.Println(req.URL().String())
	executor, err := remotecommand.NewExecutor(
		config, http.MethodPost, req.URL(),
	)
	if err != nil {
		log.Println(err)
	}
	//strings.NewReader("touc /aa.txt")
	err = executor.Stream(remotecommand.StreamOptions{
		SupportedProtocols: remotecommandconsts.SupportedStreamingProtocols,
		Stdin:              r,
		Stdout:             w,
		Stderr:             os.Stderr,
		Tty:                true,
		TerminalSizeQueue:  nil,
	})

	if err != nil {
		log.Println(err)
	}
}

func echoHandler(ws *websocket.Conn) {
	Handler(ws,ws,"test","test","default")
}


func main() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}