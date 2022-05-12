GOOS=js GOARCH=wasm go build -o solver.wasm wasm/golang/main.go
mv solver.wasm ../client/public/
