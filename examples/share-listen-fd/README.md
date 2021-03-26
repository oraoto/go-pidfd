## Share listen socket

Build:

```
go build
```

Start a http server:

```
./main
Listen fd = 3 , pid = 989509
```

In anthoer terminal, "steal"" the listen fd:

```
sudo ./main -fd 3 -pid 989509
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
