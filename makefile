export tag=v1.0.0

root:
	export ROOT=github.com/LiuJianwen/cncamp

build:
	echo "build httpserver binary"
	mkdir -p bin/amd64
	cd HTTPServer && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../bin/amd64

release: build
	echo "build httpserver container"
	sudo docker build -t wxwd14388/httpserver:${tag} .

push: release
	echo "push wxwd14388/httpserver"
	sudo docker push wxwd14388/httpserver:${tag}
