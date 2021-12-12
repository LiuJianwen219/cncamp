FROM golang:alpine AS build
COPY . /httpserver
WORKDIR /httpserver
RUN cd HTTPServer && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../bin/amd64

FROM golang:alpine
EXPOSE 80
ENV VERSION=v1.0.0
COPY --from=build /tmp/ /tmp/
COPY --from=build /httpserver/bin/amd64/HTTPServer /bin/httpserver
ENTRYPOINT ["/bin/httpserver", "-log_dir=/tmp/", "-logtostderr=true"]
