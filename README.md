Wechat 项目，使用go语言编写，使用gorm作为数据库驱动，使用gin作为web框架，使用gormgen作为数据库模型生成工具，使用gormgen作为数据库模型生成工具，使用gormgen作为数据库模型生成工具。

config目录下存放配置文件，develop.yml为应用本地开发配置，app.yml为正式应用配置文件，create_db.sh为创建数据库脚本文件，start.sh为更新镜像最新版本脚本文件。

在启动项目运行时之前，先执行create_db.sh脚本文件，创建数据库。

创建数据库后，执行start.sh脚本文件，更新镜像最新版本。

# 项目目录
- [1.正式项目目录运行时结构](#1正式项目目录运行时结构)
- [2.快速使用](#2快速使用)
    - [1. 服务器部署](#1-服务器部署)
    - [2. 上传app.yml到config目录下：](#2-上传appyml到config目录下)
    - [3. 创建create_db.sh文件：](#3-创建create_dbsh文件)
    - [4. 编写对应自动跟新镜像版本脚本start.sh文件：](#4-编写对应自动跟新镜像版本脚本startsh文件)
- [3.打包对应的架构镜像以及运行镜像](#3打包对应的架构镜像以及运行镜像)
    - [3.1 在本地打包适合x86_64架构的镜像](#31-在本地打包适合x86_64架构的镜像)
    - [3.2 保存镜像到本地](#32-保存镜像到本地)
    - [3.3 上传镜像到服务器](#33-上传镜像到服务器)
    - [3.4 登录到服务器加载镜像](#34-登录到服务器加载镜像)
    - [3.5 更新最新版本镜像并运行](#35-更新最新版本镜像并运行)


## 1.正式项目目录运行时结构

```
wechat/
├── config/
│   └── app.yml
├── create_db.sh  # 创建数据库脚本文件
└── start.sh  # 更新镜像最新版本脚本文件
```

## 2.快速使用
### 1. 服务器部署
    - 创建相关目录
    - mkdir -p /data/wechat/config
### 2. 上传app.yml到config目录下：
    - app.yml内容
    ```
        app: wechat
        port: 8081
        debug: true
        db:
        wechat:
            dialect: mysql
            dsn: root:123456@tcp(mysql:3306)/wechat?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai
        log:
        outputs:
            - stdout
            - ./logs/api.log
    ```
### 3. 创建create_db.sh文件：
    ```
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
    ```
### 4. 编写对应自动跟新镜像版本脚本start.sh文件：
    ```
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
    ```


## 3.打包对应的架构镜像以及运行镜像
### 3.1 在本地打包适合x86_64架构的镜像
首先，进入对应的项目目录下，执行以下命令：

```
docker buildx build \
  --platform linux/amd64 \
  -t wechat:1.0.0-amd64 \
  --output type=docker \
  .
```

### 3.2 保存镜像到本地

```bash
docker save -o /docker-images/wechat-1.0.0-amd64.tar wechat:1.0.0-amd64
```

### 3.3 上传镜像到服务器

```
scp /docker-images/wechat-1.0.0-amd64.tar root@192.168.10.31:/data/wechat
```
### 3.4 登录到服务器加载镜像
```
docker load -i /data/wechat/wechat-1.0.0-amd64.tar
```
### 3.5 更新最新版本镜像并运行
``` 
./start.sh
```
