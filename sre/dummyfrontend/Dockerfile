FROM golang:1.17-alpine AS build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go test -o dummyfrontend
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o dummyfrontend

FROM scratch
WORKDIR /
ENV ListenPort=
ENV BackendAddress=
COPY --from=build /build/dummyfrontend /
ENTRYPOINT ["/dummyfrontend"]
