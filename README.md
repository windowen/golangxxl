# Job

定时任务和消息队列消费项目


git remote set-url origin git@github.com:wind/goxxl.git


root@de233:/home/hroot/jenkins/run-config/workspace/go_server_dev/bin# cat /lib/systemd/system/coder-job233.service
[Unit]
Description=Tran Coder Service
After=network.target

[Service]
Type=simple
WorkingDirectory=/home/hroot/jenkins/run-config/workspace/go_server_dev/bin
ExecStart=/home/hroot/jenkins/run-config/workspace/go_server_dev/bin/coder_job_html -config config.yaml
Restart=always
RestartSec=5


[Install]


oot@de233:/home/hroot/jenkins/run-config/workspace/go_server_dev/bin# vim /lib/systemd/system/coder-job-index233.service
root@de233:/home/hroot/jenkins/run-config/workspace/go_server_dev/bin# systemctl daemon-reload
root@de233:/home/hroot/jenkins/run-config/workspace/go_server_dev/bin# systemctl enable coder-job-index233
Created symlink /etc/systemd/system/multi-user.target.wants/coder-job-index233.service -> /lib/systemd/system/coder-job-index233.service.
root@de233:/home/hroot/jenkins/run-config/workspace/go_server_dev/bin# systemctl start coder-job-index233
root@de233:/home/hroot/jenkins/run-config/workspace/go_server_dev/bin# ps aux | grep 'jobGenerateIndex'
root       62111  0.0  0.0   3328  1484 pts/0    S+   13:58   0:00 grep jobGenerateIndex
root@de233:/home/hroot/jenkins/run-config/workspace/go_server_dev/bin# cat /lib/systemd/system/coder-job-index233.service
[Unit]
Description=Tran Coder Service
After=network.target

[Service]
Type=simple
WorkingDirectory=/home/hroot/jenkins/run-config/workspace/go_server_dev/bin
ExecStart=/home/hroot/jenkins/run-config/workspace/go_server_dev/bin/coder_job -config config.yaml
Restart=always
RestartSec=5


[Install]
WantedBy=multi-user.target
root@de233:/home/hroot/jenkins/run-config/workspace/go_server_dev/bin# 



jenkins config 

#!/bin/bash
# 提权 远程复制文件
set -ex
DATE=$(date +%Y%m%d%H%M)
source /etc/profile
go version


# 打印环境变量验证
echo "HOME: $HOME"
echo "GOPATH: $GOPATH"
echo "GOCACHE: $GOCACHE"
echo "GOMODCACHE: $GOMODCACHE"
echo "PATH: $PATH"

# 检查目录是否存在
ls -ld $HOME
ls -ld $GOCACHE
ls -ld $GOMODCACHE



# 设置基本环境变量  /home/hroot/jenkins/run-config/workspace/go_server_dev
export HOME=/home/hroot/jenkins/run-config/workspace/go_server_dev  # 或其他合适的路径 /var/lib/jenkins
export GOPATH=$HOME/go
export GOCACHE=$HOME/.cache/go-build
export GOMODCACHE=$GOPATH/pkg/mod
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# 打印环境变量验证
echo "HOME: $HOME"
echo "GOPATH: $GOPATH"
echo "GOCACHE: $GOCACHE"
echo "GOMODCACHE: $GOMODCACHE"
echo "PATH: $PATH"

# 确保缓存目录存在
mkdir -p $GOCACHE
mkdir -p $GOMODCACHE

# 检查目录是否存在
ls -ld $HOME
ls -ld $GOCACHE
ls -ld $GOMODCACHE

go version


echo "[start] 开始 $DATE ..."
cd goapijob-main
go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/coder_job cmd/job/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/coder_queue cmd/queue/main.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/coder_translate cmd/tgbot/telegram.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/coder_job_html cmd/tgsinfo/telegram_info.go
cp -r config/config.yaml template bin/
cd bin/
chmod +x coder_*
exit 0

#nohup ./coder_job -config config.yaml > /dev/null 2>&1 &
#nohup ./coder_queue -config config.yaml > /dev/null 2>&1 &
#nohup ./coder_translate -config config.yaml > /dev/null 2>&1 &
#nohup ./coder_job_html -config config.yaml > /dev/null 2>&1 &

