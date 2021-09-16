#!/bin/sh

go build
pkill navani
nohup ./navani >> log &
