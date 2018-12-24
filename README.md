# WebPW

Web password generator

## Build

```bash
go build
```

Prepared docker image 

Build docker image:

```bash
# prepare custom golang image (with git)
cd docker/golang
docker build -t golang:webpw .

# build binary file for alpine
cd ..
./build

# build image
cd ..
docker build -t webpw -f docker/Dockerfile .
```

## License

This source code is governed by a MIT license that can be found in the [LICENSE](https://github.com/z0rr0/webpw/blob/master/LICENSE) file.
