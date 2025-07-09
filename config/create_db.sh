#!/bin/bash

source /data/docker/.env

echo "等待 MySQL 启动完成..."

# 等待 mysql 容器的健康检查通过
until docker inspect --format "{{json .State.Health.Status }}" mysql | grep -q '"healthy"'; do
  sleep 2
done

echo "MySQL 已就绪，开始检查并创建数据库..."

# 执行数据库初始化（注意此处是直接通过 docker exec 调用）
export MYSQL_PWD=$MYSQL_ROOT_PASSWORD
docker exec -e MYSQL_PWD=$MYSQL_PWD -i mysql mysql -uroot <<EOF
CREATE DATABASE IF NOT EXISTS wechat DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
EOF
unset MYSQL_PWD


echo "数据库创建完成"