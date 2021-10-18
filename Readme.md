## 1009模块三作业

1. 创建Dockerfile文件，编写，两端构建
2. docker build -t httpserver:v1.0.0 .
3. docker run -d --name ht -P httpserver:v1.0.0
4. docker inspect 查看容器PID
5. 使用nsenter 进入但是没有任何指令（top、ps、hostname等），可能是因为使用了scrach镜像作为运行的基镜像，后来修改为go的镜像
6. sudo nsenter -a -t $PID hostname 查看容器所在namespace的hostname和容器ID一致
7. docker tag httpserver:v1.0.0 wxwd14388/httpserver:v1.0.0 (sudo make release)
8. docker login
9. docker push wxwd14388/httpserver:v1.0.0 (sudo make push)