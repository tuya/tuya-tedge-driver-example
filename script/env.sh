#!/bin/bash
echo "machine registry.code.tuya-inc.top login tuyacoderobot password tuya@2014" > ~/.netrc

git config --global url."https://registry.code.tuya-inc.top".insteadOf "https://gitlab.com"

go env -w GOPROXY=https://goproxy.cn,direct

go env -w GOSUMDB=off