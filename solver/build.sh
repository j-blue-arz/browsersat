GOOS=js GOARCH=wasm go build -o solver.wasm
mv solver.wasm ../client/public/
