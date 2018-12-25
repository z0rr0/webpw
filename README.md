# WebPW

Web password generator

```bash
Usage of ./webpw:
  -host string
        host (default "0.0.0.0")
  -index string
        HTML template file (default "index.html")
  -port uint
        port (default 30080)
  -timeout uint
        handling timeout, seconds (default 30)
```

## Build

```bash
go build
```

Build docker image:

```bash
bash docker/build.sh
```

## License

This source code is governed by a MIT license that can be found in the [LICENSE](https://github.com/z0rr0/webpw/blob/master/LICENSE) file.
