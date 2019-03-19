#!/bin/bash

protoc -I proto/ proto/main.proto --go_out=plugins=grpc:sathelperapp
