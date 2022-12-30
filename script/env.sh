#!/bin/bash
git config --global --add url."git@github.com:".insteadOf "https://github.com/"

go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=off