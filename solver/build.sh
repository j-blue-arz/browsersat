#/bin/bash

build_wasm_golang() {
    GOOS=js GOARCH=wasm go build -o solver.wasm wasm/golang/main.go
    mv solver.wasm ../client/public/solver_golang.wasm
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ../client/public/wasm_exec_golang.js
}

build_wasm_tinygo() {
    tinygo build -o solver.wasm -target wasm wasm/tinygo/main.go
    mv solver.wasm ../client/public/solver.wasm
    cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js ../client/public/wasm_exec.js
}

if [ "$1" = "golang" ]; then
    build_wasm_golang
elif [ "$1" = "tinygo" ]; then
    build_wasm_tinygo
else
    echo "Usage: ./build.sh golang or ./build.sh tinygo"
fi


