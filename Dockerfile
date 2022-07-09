FROM tinygo/tinygo:0.24.0 as solver-base

WORKDIR /go_app

COPY solver/go.mod ./
COPY solver/go.sum ./

RUN go mod download

COPY solver/wasm ./wasm
COPY solver/maxsat ./maxsat

FROM solver-base as solver-test

RUN go test github.com/j-blue-arz/browsersat/solver/maxsat

FROM solver-base as solver-build

RUN tinygo build -o solver.wasm -target wasm wasm/main.go
RUN cp "$(tinygo env TINYGOROOT)/targets/wasm_exec.js" .

FROM node:16 as react-build

WORKDIR /react_app

COPY client ./
COPY --from=solver-build /go_app/solver.wasm public/
COPY --from=solver-build /go_app/wasm_exec.js public/
RUN npm ci
RUN npm run build

FROM nginx:latest as release

COPY --from=react-build /react_app/dist /usr/share/nginx/html