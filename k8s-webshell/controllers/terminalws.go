package controllers

import (
	"github.com/astaxie/beego"
	"golang.org/x/net/websocket"
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"strconv"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
)


type terminalsize struct {
	C chan *remotecommand.TerminalSize
}

func (self *terminalsize) Next() *remotecommand.TerminalSize {
	return <-self.C
}

func buildConfigFromContextFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func Handler(r io.Reader, w io.Writer, context string, namespace string, podname string, container string, rows string, cols string,cmd string) error {
	config, err := buildConfigFromContextFlags(context, beego.AppConfig.String("kubeconfig"))
	if err != nil {
		return err
	}
	groupversion := schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}
	config.GroupVersion = &groupversion
	config.APIPath = "/api"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	restclient, err := rest.RESTClientFor(config)
	if err != nil {
		return err
	}

	req := restclient.Post().
		Resource("pods").
		Name(podname).
		Namespace(namespace).
		SubResource("exec").
		Param("container", container).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("command", cmd).Param("tty", "true")
	c := make(chan *remotecommand.TerminalSize)
	t := &terminalsize{c}
	req.VersionedParams(
		&v1.PodExecOptions{
			Container: container,
			Command:   []string{},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		},
		scheme.ParameterCodec,
	)
	executor, err := remotecommand.NewSPDYExecutor(
		config, http.MethodPost, req.URL(),
	)
	if err != nil {
		return err
	}

	ro, err := strconv.Atoi(rows)
	if err != nil {
		return err
	}
	wo, err := strconv.Atoi(cols)
	if err != nil {
		return err
	}
	go func() {
		c <- &remotecommand.TerminalSize{
			Width:  uint16(ro),
			Height: uint16(wo),
		}
	}()
	err = executor.Stream(remotecommand.StreamOptions{
		//SupportedProtocols: remotecommandconsts.SupportedStreamingProtocols,
		Stdin:              r,
		Stdout:             w,
		Stderr:             w,
		Tty:                true,
		TerminalSizeQueue:  t,
	})
	return err
}

func EchoHandler(ws *websocket.Conn) {
	defer ws.Close()
	r := ws.Request()
	context := r.FormValue("context")
	namespace := r.FormValue("namespace")
	pod := r.FormValue("pod")
	container := r.FormValue("container")
	rows := r.FormValue("rows")
	cols := r.FormValue("cols")
	beego.Info(context, namespace, pod, container, rows, cols)
	if err:=Handler(ws, ws, context, namespace, pod, container, cols, rows,"/bin/bash"); err!=nil {
		beego.Error(err)
		log.Println(Handler(ws, ws, context, namespace,  pod, container, cols, rows,"/bin/sh"))
	}
}
