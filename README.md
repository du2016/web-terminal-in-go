## web-terminal-in-go

### container-webshell


通过web界面连接容器，需要开放docker api,需要传递以下参数：

- 宿主机ip
- 端口
- 容器id


### k8s-webshell


通过web界面连接k8s容器，需要修改beego config配置，指定kubeconfig位置，设置为空则为incluster连接方式，需要传递以下参数：

- namespace
- podname
- containername
