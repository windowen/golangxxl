root@de233:~# ls -al /lib/systemd/system |grep coder
-rw-r--r--  1 root root   389 Jul 23 11:17 coder-jenkins39017.service
-rw-r--r--  1 root root   324 Aug  2 13:57 coder-job-index233.service
-rw-r--r--  1 root root   329 Aug  1 17:30 coder-job233.service
-rw-r--r--  1 root root   302 Aug  1 15:24 coder-tran233.service
-rw-r--r--  1 root root   274 Aug  2 10:53 coder-xxl-job39270.service

root@de233:~# cat /lib/systemd/system/coder-jenkins39017.service
[Unit]
Description=Jenkins Daemon
After=syslog.target network-online.target
Wants=network-online.target

[Service]
Type=notify
ExecStart=/usr/bin/java -jar -DJENKINS_HOME=/home/hroot/jenkins/run-config /home/hroot/jenkins/jenkins.war --httpPort=39017
RuntimeDirectory=/home/hroot/jenkins/run-config
ReadWriteDirectories=/home/hroot/jenkins/run-config

[Install]
WantedBy=multi-user.target
root@de233:~# 

job 单个生成 单个的html 机器人
root@de233:~# cat /lib/systemd/system/coder-job-index233.service
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

# 生成首页的 xxljob 调度
root@de233:~# cat /lib/systemd/system/coder-job233.service
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
WantedBy=multi-user.target

# 翻译机器人
root@de233:~# cat /lib/systemd/system/coder-tran233.service
[Unit]
Description=Tran Coder Service
After=network.target

[Service]
Type=simple
WorkingDirectory=/home/hroot/jenkins/run-config/workspace/go_server_dev
ExecStart=/home/hroot/jenkins/run-config/workspace/go_server_dev/coder_translate
Restart=always
RestartSec=5

# xxljob的基本配置
[Install]
WantedBy=multi-user.target

root@de233:~# cat /lib/systemd/system/coder-xxl-job39270.service
[Unit]
Description=XXL-Job Admin Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/bin/java -jar /home/hroot/libs/xxl-job-3.1.1/xxl-job-admin-3.1.1.jar
Restart=always
RestartSec=5
StandardOutput=null
StandardError=null

[Install]
WantedBy=multi-user.target


### 服务查找
2. 检查进程
root@de233:~# ps aux | grep coder_job_html
root        2743  0.0  0.1 1242084 22688 ?       Ssl  Aug04   0:05 /home/hroot/jenkins/run-config/workspace/go_server_dev/bin/coder_job_html -config config.yaml
root       65219  0.0  0.0   3328  1388 pts/0    S+   10:37   0:00 grep coder_job_html

如何查找关联的 service？
1. 检查 systemd 服务
   运行以下命令查找是否有对应的 systemd 服务：
2. systemctl list-units --type=service | grep -i coder_job_html

或检查所有 service 文件：
root@de233:~# grep -l "coder_job_html" /etc/systemd/system/*.service /lib/systemd/system/coder* 2>/dev/null
/lib/systemd/system/coder-job233.service