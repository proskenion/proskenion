# proskenion

[![CircleCI](https://circleci.com/gh/proskenion/proskenion.svg?style=svg)](https://circleci.com/gh/proskenion/proskenion)

Blockchain project for creators.

## Develop Environment
- glide version 0.13.1
- go version go1.11
- libprotoc 3.6.0

## Boot proskenion
```
docker build -t proskenion .
docker run proskenion
```

### Example
```$xslt
make example
make ipset
docker run -p $LOCAL_HOST_IP:50052:50052 proskenion -c example/configRoot.yaml
```
```
make ipset
docker run -p $LOCAL_HOST_IP:50053:50053 proskenion:latest -c example/config1.yaml
```
```
make ipset
docker run -p $LOCAL_HOST_IP:50054:50054 proskenion:latest -c example/config2.yaml
```
```
make ipset
docker run-p $LOCAL_HOST_IP:50055:50055 proskenion:latest -c example/config3.yaml
```
```$xslt
./bin/example
```

