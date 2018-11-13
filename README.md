# mssh
这是一个通过SSH在多机器执行相同命令的小工具。

## 安装
`go get -u github.com/zxfishhack/mssh`

## 配置文件格式
```yaml
- ip: 192.168.2.2 # 机器的IP地址
  port: 22        # SSH端口
  username: root  # 登录使用的用户名
  password: pass  # 登录使用的密码
  path: /root     # 执行命令的工作目录
```

## 使用
`mssh config.yaml egrep -h foo bar*.log`

## TODO
1. 增加公钥登录支持
1. 增加只往配置文件的特定机器执行命令功能
