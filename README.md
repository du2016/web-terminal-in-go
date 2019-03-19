## web-terminal-in-go

### k8s-webshell

通过web界面连接k8s容器，需要修改beego config配置，指定kubeconfig位置，设置为空则为incluster连接方式，需要传递以下参数：
- context
- namespace
- podname
- containername


#### 功能

- web终端实现
- 多集群支持
- 根据浏览器窗口调整tty大小

***

### container-webshell


通过web界面连接容器，需要开放docker api,需要传递以下参数：

- host ip
- docker port
- container id

***

### 演示

![demo](./demo.gif)