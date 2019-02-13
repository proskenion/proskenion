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
docker run -p 127.0.0.1:50052:50052 proskenion -c example/configRoot.yaml
```
```
docker run proskenion:latest -c example/config1.yaml
```
```
docker run proskenion:latest -c example/config2.yaml
```
```
docker run proskenion:latest -c example/config3.yaml
```
```$xslt
./bin/example
```

