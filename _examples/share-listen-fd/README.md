## Share listen socket

`PTRACE_MODE_ATTACH_REALCREDS` permission is required:

```
echo 0 | sudo tee /proc/sys/kernel/yama/ptrace_scope
```

Build:

```
go build
```

Start a http server:

```
./share-listen-fd
Listen fd = 3 , pid = 989509
```

In anthoer terminal, "steal"" the listen fd:

```
sudo ./share-listen-fd -fd 3 -pid 989509
```

Then, send some http request:

```
$ curl 127.0.0.1:8080
Hello from 997217
$ curl 127.0.0.1:8080
Hello from 997217
$ curl 127.0.0.1:8080
Hello from 997217
$ curl 127.0.0.1:8080
Hello from 989509      <-- the new server
$ curl 127.0.0.1:8080
Hello from 997217
$ curl 127.0.0.1:8080
Hello from 989509
```
