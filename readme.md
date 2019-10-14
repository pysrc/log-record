## 日志收集系统

### 服务端(server)

config.json配置数据库连接信息，第一次会生成log_record日志表

### 客户端(client)

example:

`tail -lf test.log | ./LogClient -s ClientSystem -h http://127.0.0.1:9587 &`

只需要将待收集的日志输出管道输出到收集客户端即可，客户端会同步到服务端

