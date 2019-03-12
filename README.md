<div align="center"><img src=https://user-images.githubusercontent.com/6259384/52863883-42ec5100-317c-11e9-89f4-640f7bd26938.png "proskenion"></div>

[![CircleCI](https://circleci.com/gh/proskenion/proskenion.svg?style=svg)](https://circleci.com/gh/proskenion/proskenion)

It is a BlockChain platform that bases between the content provider and the spectator. So, "Proskenion" named from the etymology of "Prosium Arch" (Greek).

- Having a high expression power by a combination of primitive instruction sets.
- Easy to customize incentive / consensus algorithm.
- Can change the incentive / consensus algorithm without hard fork.

A decentralized and creator-based management system can be realized with Proskenion.

Proskenion's incentive/consensus algorithm is designed by [Prosl](https://github.com/proskenion/proskenion/tree/master/prosl).

Proskenion documentations is [here](https://proskenion.github.io).

[![demo movie](http://img.youtube.com/vi/UyQdHKpAxJ0/1.jpg)](https://www.youtube.com/embed/UyQdHKpAxJ0)

## Develop Environment
- glide version 0.13.1
- go version go1.11
- libprotoc 3.6.0

## Boot proskenion
```
$ git clone https://github.com/proskenion/proskenion.git
$ cd proskenion
$ docker build -t proskenion .
$ docker run proskenion
```

## Example (for Mac)
[Example: https://github.com/proskenion/proskenion/tree/master/example](https://github.com/proskenion/proskenion/tree/master/example)
### terminal1
```
$ make example
$ make ipset
$ docker run -p $LOCAL_HOST_IP:50052:50052 proskenion -c example/configRoot.yaml
```
### terminal2
```
$ make ipset
$ docker run -p $LOCAL_HOST_IP:50053:50053 proskenion:latest -c example/config1.yaml
```
### terminal3
```
$ make ipset
$ docker run -p $LOCAL_HOST_IP:50054:50054 proskenion:latest -c example/config2.yaml
```
### terminal4
```
$ make ipset
$ docker run-p $LOCAL_HOST_IP:50055:50055 proskenion:latest -c example/config3.yaml
```
### terminal5
```
$ ./bin/example
```

