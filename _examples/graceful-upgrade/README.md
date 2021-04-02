## Graceful upgrade http server

`PTRACE_MODE_ATTACH_REALCREDS` permission is required:

```
echo 0 | sudo tee /proc/sys/kernel/yama/ptrace_scope
```

Build:

```
go build
```

Start http server:

```
./graceful-upgrade
```

Send some request and upgrade

```
$ curl 127.0.0.1:8001
Hello from 2893674      <- old server

$ kill -HUP 2893674     # upgrade

$ curl 127.0.0.1:8001
Hello from 2893969      <- new server
```
