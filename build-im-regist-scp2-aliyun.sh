#!/bin/bash

set -e

VERSION=$1
CUR_VERSION=""
IMAGE_NAME=registerserver
if test -z $VERSION; then
   echo "请输入版本号!!!"
   exit 1
fi

if [ "$(echo $VERSION | grep "beta-")" != "" ]; then
  CUR_VERSION="develop"
fi

if [ "$(echo $VERSION | grep "v-")" != "" ]; then
  CUR_VERSION="release"
fi

if [ "$CUR_VERSION" = "" ]; then
  echo "版本号错误,无法推送!!!"
  exit 1
fi

echo "开始编译..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o registerserver cmd/open_im_register_api.go
echo "编译完成..."

# 写入版本号
echo "写入版本号"
echo "$1" > ./config/VERSION

echo "生成Dockerfile文件"
cat>Dockerfile<<EOF
FROM alpine

WORKDIR /root/

COPY registerserver /root/registerserver
COPY config /root/config

CMD ["/root/registerserver"]
EOF

echo "生成docker-compose.yml文件"
cat>docker-compose.yml<<EOF
version: '3'
services:
  registerserver:
    build: .
    container_name: 'registerserver'
    ports:
      - "42233:42233"
    restart: always
    networks:
      - traefik
    labels:
      - "traefik.enable=true"
      - "traefik.docker.network=traefik"
    logging:
      driver: "json-file"
      options:
        max-size: "1m"
    volumes:
      - ./config:/root/config
      - ./logs:/root/logs
      - /etc/timezone:/etc/timezone
      - /etc/localtime:/etc/localtime

networks:
  traefik:
    external: true
EOF

echo "压缩文件"
tar -zcvf .build/registerserver.tgz registerserver Dockerfile docker-compose.yml config/VERSION config/config.yaml
rm registerserver
echo "scp 文件到服务器"
scp .build/registerserver.tgz aliyun-stone:/root/registerserver

# run.sh
echo "生成 run.sh"
cat>.build/run.sh<<EOF
tar zxvf registerserver.tgz -C ./
docker-compose down && docker-compose up -d --build \
&& rm registerserver && rm registerserver.tgz
EOF

#echo "更新最新版本库"
#docker tag $(docker images | grep ${IMAGE_NAME} | head -1 | awk '{print $3}') ${IMAGE_NAME}:$CUR_VERSION
#
#echo "清理本地无用镜像..."
#docker rmi -f $(docker images | grep ${IMAGE_NAME} | awk '{print $3}')

echo "完事"

chmod +x .build/run.sh
echo "scp run.sh 文件到服务器"
scp .build/run.sh aliyun-stone:/root/registerserver
