#/bin/bash

if [ "$1" = "golang" ]; then
    GOOS=js GOARCH=wasm go build -o solver.wasm wasm/golang/main.go
    mv solver.wasm ../client/public/
elif [ "$1" = "tinygo" ]; then
    tinygo build -o solver.wasm -target wasm wasm/tinygo/main.go
    mv solver.wasm ../client/public/
else
    echo "Usage: ./build.sh golang or ./build.sh tinygo"
fi


