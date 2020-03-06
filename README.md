# Kubernetes DNS

[![Build Status](https://travis-ci.org/kubernetes/dns.svg?branch=master)](https://travis-ci.org/kubernetes/dns)
[![Coverage Status](https://coveralls.io/repos/github/kubernetes/dns/badge.svg?branch=master)](https://coveralls.io/github/kubernetes/dns?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubernetes/dns)](https://goreportcard.com/report/github.com/kubernetes/dns)

This is the repository for [Kubernetes DNS](http://kubernetes.io/docs/admin/dns/).

## Images

* [kube-dns](http://kubernetes.io/docs/admin/dns/)
* [sidecar](docs/sidecar/README.md)
* [dnsmasq](images/dnsmasq)
* [node-cache](http://kubernetes.io/docs/tasks/administer-cluster/nodelocaldns/)

## Building

`make` targets:

| target | description |
| ---- | ---- |
|all, build   | build all binaries |
|test         | run unit tests |
|containers   | build the containers |
|images-clean | clear image build artifacts from workdir |
|push         | push containers to the registry |
|help         | this help message |
|version      | show package version |
|{build,containers,push}-ARCH | do action for specific ARCH |
|all-{build,containers,push}  | do action for all ARCH |
|only-push-BINARY             | push just BINARY |

* Setting `VERBOSE=1` will show additional build logging.
* Setting `VERSION` will override the container version tag.

[![Analytics](https://kubernetes-site.appspot.com/UA-36037335-10/GitHub/dns/README.md?pixel)]()

## Release process

1. Build and test (`make images-clean`; `make build`; `make containers`; `make test`)
1. Update [go dependencies](docs/go-dependencies.md) if needed.
1. Update the release tag. We use [semantic versioning](http://semver.org) to
   name releases.
1. Push the containers (`make push`)
1. Submit a PR for the kubernetes/kubernetes repository to switch to the new
   version of the containers.
1. Build and push for all architectures (`make all-push`)

## Upgrading CoreDNS version in node-cache
The coreDNS version used in node-cache image will be atmost 2 minor versions behind the latest CoreDNS image.
Any vulnerability fixes will be picked up soon after a CoreDNS release is published.
