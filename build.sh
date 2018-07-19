#!/bin/bash

sudo go build knitctl.go

sudo cp knitctl /usr/local/bin/

#sudo rm -f knitctl

echo "build success. /usr/local/bin/knitctl"
