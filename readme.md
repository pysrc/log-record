## 日志收集系统

### 服务端(server)

### 客户端(client)

example:

`tail -f test.log | ./LogClient -s svc-demo -i svc-id -h http://127.0.0.1:9587 &`

只需要将待收集的日志输出管道输出到收集客户端即可，客户端会同步到服务端

### 相同服务日志集中

`http://127.0.0.1:9587/svc?svc=svc-demo`

在server上即可集中svc-demo服务的日志svc-demo.log

