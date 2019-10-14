# Themis

[![Build Status](https://travis-ci.com/xmidt-org/themis.svg?branch=master)](https://travis-ci.com/xmidt-org/themis)
[![codecov.io](http://codecov.io/github/xmidt-org/themis/coverage.svg?branch=master)](http://codecov.io/github/xmidt-org/themis?branch=master)
[![Code Climate](https://codeclimate.com/github/xmidt-org/themis/badges/gpa.svg)](https://codeclimate.com/github/xmidt-org/themis)
[![Issue Count](https://codeclimate.com/github/xmidt-org/themis/badges/issue_count.svg)](https://codeclimate.com/github/xmidt-org/themis)
[![Go Report Card](https://goreportcard.com/badge/github.com/xmidt-org/themis)](https://goreportcard.com/report/github.com/xmidt-org/themis)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/xmidt-org/themis/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/xmidt-org/themis?include_prereleases)](CHANGELOG.md)

## Summary

A JWT issuer for the CPE devices that connect to the XMiDT cluster.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How to Install](#how-to-install)
- [Usage](#usage)
- [Contributing](#contributing)

## Code of Conduct

This project and everyone participating in it are governed by the [XMiDT Code Of Conduct](https://xmidt.io/code_of_conduct/). 
By participating, you agree to this Code.

## How to Install

### Installation
- [Docker](https://www.docker.com/) (duh)
  - `brew install docker`

</br>

### Running
#### Build the docker image
```bash
docker build -t themis:local .
```
This `build.sh` script will build the binary and docker image

## Usage
Once everything is up and running you can start sending requests. Below are a few examples.
TODO: Add examples

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md).
