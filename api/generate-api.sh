#!/bin/sh
oapi-codegen -generate types api.v1.yaml > types.gen.go
oapi-codegen -generate server api.v1.yaml > server.gen.go
