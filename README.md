# About

This is the test example  "Social tournament service"


# Install

## Get project

- `git clone https://github.com/vxdiv/sts.git sts/src/sts`

-  set environment vars $GOROOT  and $GOPATH

## Install dependency

-  [Docker Compose](https://docs.docker.com/compose/install/)

-  [Glide (Vendor Package Management)](https://github.com/Masterminds/glide#install)

## Prepare environment

-  Start DB instance: `docker-compose up -d`

-  Build application `make all`

-  Add connection to DB in config file `config.yaml`  (`default: localhost:27017`)


# Usage

- Run application: `/sts`
