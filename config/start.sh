#!/bin/bash

IMAGE_PREFIX="wechat"
ARCH_SUFFIX="amd64"

# 获取本地可用的最大版本号（假设版本号格式是 x.y.z）
LATEST_TAG=$(docker images --format "{{.Repository}}:{{.Tag}}" \
  | grep "^${IMAGE_PREFIX}:" \
  | grep "${ARCH_SUFFIX}" \
  | sed -E "s/^${IMAGE_PREFIX}:([0-9]+\.[0-9]+\.[0-9]+)-${ARCH_SUFFIX}$/\1/" \
  | sort -Vr \
  | head -n 1)

if [[ -z "$LATEST_TAG" ]]; then
  echo "❌ 未找到符合格式的镜像版本（wechat:x.y.z-amd64）"
  exit 1
fi

FULL_IMAGE="${IMAGE_PREFIX}:${LATEST_TAG}-${ARCH_SUFFIX}"

echo "✅ 即将运行镜像: ${FULL_IMAGE}"

# 停止并删除旧容器（如存在）
docker rm -f wechat-app 2>/dev/null

# 启动容器
docker run -d \
  --name wechat-app \
  --restart=always \
  -p 8081:8081 \
  -v /data/wechat/config/app.yml:/app/config/app.yml \
  --network common-app-net \
  "$FULL_IMAGE"