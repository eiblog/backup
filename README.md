# backup

用于备份 Eiblog 数据库到七牛云

你可以使用如下命令查看帮助：
```
$ registry.cn-hangzhou.aliyuncs.com/deepzz/backup -h
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
环境变量将会覆盖命令行参数。
