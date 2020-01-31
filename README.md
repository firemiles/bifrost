# Bifrost

[![Build Status](https://semaphoreci.com/api/v1/firemiles/bifrost-3/branches/master/shields_badge.svg)](https://semaphoreci.com/firemiles/bifrost-3)
[![Build Status](https://travis-ci.com/firemiles/bifrost.svg?branch=master)](https://travis-ci.com/firemiles/bifrost)
[![Coverage Status](https://coveralls.io/repos/github/firemiles/bifrost/badge.svg?branch=master)](https://coveralls.io/github/firemiles/bifrost?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/firemiles/bifrost)](https://goreportcard.com/report/github.com/firemiles/bifrost)

In Norse mythology, Bifröst is a burning rainbow bridge that reaches between Midgard (Earth) and Asgard, the realm of the gods

![](https://upload.wikimedia.org/wikipedia/commons/thumb/e/e3/Heimdall_an_der_Himmelsbr%C3%BCcke.jpg/440px-Heimdall_an_der_Himmelsbr%C3%BCcke.jpg)

Project `bifrost` is a CNI plugin inspired by [aws/containers-roadmap](https://github.com/aws/containers-roadmap/issues/398), aims to provide a centralized `IPAM` on any `Cloud Provider`

## Feature

1. Provide a `IPAM` Plugin, which provide two ip allocation methods: `IPBlock` and `Endpoint`.
2. Support add `Cloud Provide` network card to Pod , such as `AWS ENI` 
3. Provide `framework` to add custom network type. 