### 开发环境
```
go version:go1.17.1
mysql version:5.7.36
```

### 初始化数据库
```
// 创建数据库logviewer设置字符集为utf8mb4
create database logviewer default character set utf8mb4;

// 未开启远程登录
// 创建用户并授权用户所有权限给logproject
CREATE USER 'logviewer'@'%' IDENTIFIED BY 'pass';
GRANT all ON logviewer.* TO 'logviewer'@'%';

// 导入初始化sql
use logviewer
source sql/init.sql
```

### 更新Server端配置文件
```
# logViewerServer端程序侦听端口
port = 9000

# 运行模式，true为生产环境，false为测试环境
release = false

# 日志压缩包上传路径，要保证服务启动的用户对此目录有读写权限
path = storage

# 日志服务端登陆地址
download = http://127.0.0.1:9000/logviewer

# 数据库配置信息
[mysql]
user = logviewer
password = password
host = 127.0.0.1
port = 3306
db = logviewer
```

### 启动server端服务
```
nohup ./logViewerServer conf/config.ini > logViewerServer.log &
```


