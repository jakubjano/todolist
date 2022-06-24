##             MonitorAPIs microservice Makefile
##
##  Simple Makefile containing implementation of targets for generating protobuf file
##
##  To generate swagger, angular models and go messages&grpc use this command:
##    $ make gen
SHELL=/bin/bash

PROTO_INCLUDES=-Iproto -I3rdparty -I=.
PROTO_GRPC_ARGS=paths=source_relative
##
##  \e[1mTargets\e[0m
##   \e[34mhelp\e[0m
##       Shows this help
help:
	@echo -e "$$(sed -n 's/^##//p' Makefile)"

##   \e[34mgen\e[0m
##       Shortcut for generate
gen:
	make generate SERVICE_NAME=task
	make generate/go SERVICE_NAME=task
	make generate/gateway SERVICE_NAME=task

	make generate SERVICE_NAME=user
	make generate/go SERVICE_NAME=user
	make generate/gateway SERVICE_NAME=user



##   \e[34mgenerate\e[0m
##       Generates Go and Angular models
generate: generate/go

##   \e[34mgenerate/go\e[0m
##       Generates go grpc files and messages from proto file
generate/go:
	protoc $(PROTO_INCLUDES) \
		  proto/v1/${SERVICE_NAME}.proto \
		   --go_out=$(PROTO_GRPC_ARGS):go-sdk/${SERVICE_NAME} \
		   --go-grpc_out=$(PROTO_GRPC_ARGS):go-sdk/${SERVICE_NAME} \


##   \e[34mgenerate/gateway\e[0m
##       Generates grpc-gateway files from proto file
generate/gateway:
	protoc $(PROTO_INCLUDES) \
	  proto/v1/${SERVICE_NAME}.proto \
	   --grpc-gateway_out=logtostderr=true,paths=source_relative:go-sdk/${SERVICE_NAME} \
	   --openapiv2_out=logtostderr=true:.