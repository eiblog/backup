# backup

Only for eiblog v1, so archived.

用于备份 Eiblog 数据库到七牛云，程序使用 `mongodump` 命令进行数据导出，并进行 gzip 压缩，最后上传到七牛云。默认文件名为 `eiblog-$date.tag.gz`。程序将连接 mongodb（host 地址）。如未使用 eiblog 下的 docker-compose.yml ，请自行寻找方法。

使用 [docker-compose.yml](https://github.com/eiblog/eiblog) 运行的朋友请修改该文件的相关配置。

> 注意，最好再新建一个私有的 bucket 进行数据库的备份。

你可以使用如下命令查看帮助：
```
$ docker run --rm registry.cn-hangzhou.aliyuncs.com/deepzz/backup -h
Usage of ./backup:
  -ak string
    	qiniu AccessKey, must
  -b string
    	qiniu bucket, must
  -d string
    	bucket's domain, must
  -db string
    	back up which database (default "eiblog")
  -i string
    	how long interval backup, like 7s/7m/7h/7d (default "7d")
  -l int
    	how many days (default 60)
  -sk string
    	qiniu SecretKey, must
```

你可以通过命令行指定参数：
```
$ docker run --rm registry.cn-hangzhou.aliyuncs.com/deepzz/backup \
  -l 60 \
  -i 7d \
  -db eiblog \
  -b eiblog \
  -d xx.example.com \
  -ak xxxxxxxxxxxx \
  -sk xxxxxxxxxxxx
```

也可以通过环境变量指定参数：
```
$ docker run --rm \
  -e BACKUP_LONG=60 \
  -e BACKUP_INTERVAL=7d \
  -e BACKUP_DB=eiblog \
  -e QINIU_BUCKET=eiblog \
  -e QINIU_DOMAIN=xx.example.com \
  -e ACCESS_KEY=xxxxxxxxxx \
  -e SECRECT_KEY=xxxxxxxxxx \
  registry.cn-hangzhou.aliyuncs.com/deepzz/backup
```
这里的配置可以解读为：数据库 eiblog 将会每 7 天备份一次，每次备份数据保留 60 天，备份在七牛云的 eiblog。环境变量将会覆盖命令行参数。
