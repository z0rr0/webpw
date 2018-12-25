# WebPW

Web password generator.

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

Based on [github.com/z0rr0/gopwgen/pwgen](https://github.com/z0rr0/gopwgen) package.

## Build

```bash
go build
```

Build docker [image](https://cloud.docker.com/repository/docker/z0rr0/webpw):

```bash
bash docker/build.sh
```

## License

This source code is governed by a MIT license that can be found in the [LICENSE](https://github.com/z0rr0/webpw/blob/master/LICENSE) file.
