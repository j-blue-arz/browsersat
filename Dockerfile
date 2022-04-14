FROM golang:1.17-alpine as solver-base

WORKDIR /go_app

COPY solver/go.mod ./
COPY solver/go.sum ./

RUN go mod download

COPY solver/main.go ./
COPY solver/sat ./sat

FROM solver-base as solver-test

RUN go test github.com/j-blue-arz/browsersat/solver/sat

FROM solver-base as solver-build

RUN GOOS=js GOARCH=wasm go build -o solver.wasm
RUN cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .

FROM node:16 as react-build

WORKDIR /react_app

COPY client ./
COPY --from=solver-build /go_app/solver.wasm public/
COPY --from=solver-build /go_app/wasm_exec.js public/
RUN npm ci
RUN npm run build

FROM nginx:latest as release

COPY --from=react-build /react_app/dist /usr/share/nginx/html