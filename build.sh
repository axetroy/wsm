#!/bin/bash

# Reference:
# https://github.com/golang/go/blob/master/src/go/build/syslist.go
os_archs=(
    darwin/amd64
    linux/amd64
    windows/amd64
)

releases=()
fails=()

for os_arch in "${os_archs[@]}"
do
    goos=${os_arch%/*}
    goarch=${os_arch#*/}

    filename=terminal

    if [[ ${goos} == "windows" ]];
    then
        filename+=.exe
    fi

    echo building ${os_arch}

    CGO_ENABLED=0 GOOS=${goos} GOARCH=${goarch} go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-s -w" -o ./bin/${filename} ./cmd/user/main.go >/dev/null 2>&1

    # if build success
    if [[ $? == 0 ]];then
        releases+=(${os_arch})
        cd ./bin

        tar -czf terminal_${goos}_${goarch}.tar.gz ${filename}

        rm -rf ./${filename}

        cd ../
    else
        fails+=(${os_arch})
    fi
done

echo

if [[ -n "$fails" ]]; then
    echo "fails:"

    for os_arch in "${fails[@]}"
    do
        printf "\t%s\n" "${os_arch}"
    done
fi


if [[ -n "releases" ]]; then
    echo "release:"

    for os_arch in "${releases[@]}"
    do
        printf "\t%s\n" "${os_arch}"
    done
else
    echo "there's no build success"
    exit 1
fi

echo