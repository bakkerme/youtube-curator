package api

//go:generate oapi-codegen -generate types -package api -o types.gen.go api.v1.yaml
//go:generate oapi-codegen -generate server -package api -o server.gen.go api.v1.yaml
