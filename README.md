# logging-go

A go logging lib, with a lot of logrus hooks


The base codes are from [hellofresh/logging-go](https://github.com/hellofresh/logging-go). Change all the hooks.



# usage

yaml config

```
level: info
format: text
formatSettings:
  ts: RFC3339Nano
writer: stderr
hooks:
- type: file
  settings: {name: myProject.log, keep: 7, path: logs}
- type: sentry
  settings: {dsn: mySentryDSN}
- type: redis
  settings: {host: 127.0.0.1, port: 6379, db: 0, key: apigateway, logformat: logstashv1, poolsize: 5}
```


# supported hooks

- file
- redis
- sentry



