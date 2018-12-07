package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"golang.org/x/net/websocket"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubernetes/pkg/util/interrupt"
	"github.com/docker/docker/pkg/term"
	"net/http"
)

func (self terminalsize) Read(p []byte) (int, error) {
	var reply string
	var msg map[string]uint16
	if err := websocket.Message.Receive(self.conn, &reply); err != nil {
		return 0, err
	}
	if err := json.Unmarshal([]byte(reply), &msg); err != nil {
		return copy(p, reply), nil
	} else {
		self.sizeChan <- &remotecommand.TerminalSize{
			msg["cols"],
			msg["rows"],
		}
		return 0, nil
	}
}

type terminalsize struct {
	conn     *websocket.Conn
	sizeChan chan *remotecommand.TerminalSize
}

func (self *terminalsize) Next() *remotecommand.TerminalSize {
	size := <-self.sizeChan
	beego.Debug("terminal size to width: %s height: %s", size.Width,size.Height)
	return size
}

func buildConfigFromContextFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func Handler(ws *websocket.Conn, context string, namespace string, podname string, container string, cmd string) error {
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
	fn := func() error {
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
		t := &terminalsize{ws, c}
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
		return executor.Stream(remotecommand.StreamOptions{
			Stdin:             t,
			Stdout:            ws,
			Stderr:            ws,
			Tty:               true,
			TerminalSizeQueue: t,
		})
	}
	inFd, isTerminal := term.GetFdInfo(ws)
	beego.Info(isTerminal)
	state, err := term.SaveState(inFd)
	return interrupt.Chain(nil, func() {
		term.RestoreTerminal(inFd, state)
	}).Run(fn)
}

func EchoHandler(ws *websocket.Conn) {
	defer ws.Close()
	r := ws.Request()
	context := r.FormValue("context")
	namespace := r.FormValue("namespace")
	pod := r.FormValue("pod")
	container := r.FormValue("container")
	beego.Debug("connect context: %s namespace: %s pod: %s container: %s ", context, namespace, pod, container)
	if err := Handler(ws, context, namespace, pod, container, "/bin/bash"); err != nil {
		beego.Error(err)
		beego.Error(Handler(ws, context, namespace, pod, container, "/bin/sh"))
	}
	Handler(ws, context, namespace, pod, container, "exit")
}
