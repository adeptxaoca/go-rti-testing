FROM golang:1.15.3-alpine AS build
WORKDIR $GOPATH/src/go-rti-testing
COPY . $GOPATH/src/go-rti-testing
RUN CGO_ENABLED=0 go build -installsuffix cgo -ldflags="-w -s" -o /bin/server .

FROM scratch
COPY --from=build /bin/server /bin/server
ENTRYPOINT ["/bin/server"]
EXPOSE 8080