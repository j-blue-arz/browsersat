GOOS=js GOARCH=wasm go build -o solver.wasm
cp solver.wasm ../client/public/
