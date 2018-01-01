// Package main provides ...
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

var (
	interval  time.Duration // 间隔时间
	long      int           // 保留天数
	backupDB  string        // 备份数据库
	bucket    string        // 七牛云空间
	domain    string        // 七牛云域名
	accessKey string        // 七牛云 AccessKey
	secretKey string        // 七牛云 SecretKey

	err   error
	inter string
)

func init() {
	flag.IntVar(&long, "l", 30, "how many days")
	flag.StringVar(&inter, "i", "7d", "how long interval backup, like 7s/7m/7h/7d")
	flag.StringVar(&backupDB, "db", "eiblog", "back up which database")
	flag.StringVar(&bucket, "b", "", "qiniu bucket, must")
	flag.StringVar(&domain, "d", "", "bucket's domain, must")
	flag.StringVar(&accessKey, "ak", "", "qiniu AccessKey, must")
	flag.StringVar(&secretKey, "sk", "", "qiniu SecretKey, must")
}

func main() {
	flag.Parse()

	// days
	l := os.Getenv("BACKUP_LONG")
	if l != "" {
		days, err := strconv.Atoi(l)
		if err != nil {
			panic(err)
		}
		long = days
	}
	// interval
	i := os.Getenv("BACKUP_INTERVAL")
	if i != "" {
		inter = i
	}
	interval, err = parseDuration(inter)
	if err != nil {
		log.Panic(err)
	}
	// db
	db := os.Getenv("BACKUP_DB")
	if db != "" {
		backupDB = db
	}
	// bucket
	b := os.Getenv("QINIU_BUCKET")
	if b != "" {
		bucket = b
	}
	if bucket == "" {
		log.Panic("which bucket do you want to save?")
	}
	// domain
	d := os.Getenv("QINIU_DOMAIN")
	if d != "" {
		domain = d
	}
	if domain == "" {
		log.Panic("we need the bucket's domain.")
	}
	// key
	ak := os.Getenv("ACCESS_KEY")
	sk := os.Getenv("SECRECT_KEY")
	if ak != "" || sk != "" {
		accessKey = ak
		secretKey = sk
	}
	if accessKey == "" || secretKey == "" {
		log.Panic("we need the accessKey and secretKey")
	}

	// mongodump
	ch := make(chan string)
	go mongoDump(ch)
	go qiniuUpload(ch)

	select {}
}

func mongoDump(ch chan string) {
	ips, err := net.LookupIP("mongodb")
	if err != nil {
		log.Panic(err)
	}
	log.Println(ips)
	if len(ips) == 0 {
		log.Panic("not found host: mongodb")
	}

	t := time.NewTicker(interval)
	for {
		now := <-t.C

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*20)
		defer cancel()
		cmd := exec.CommandContext(ctx, "mongodump", "-h", ips[0].String()+":27017", "-d", backupDB, "-o", "/tmp")
		err = cmd.Run()
		if err != nil {
			log.Println("mongodump", err)
			continue
		}

		name := fmt.Sprintf("%s-%s.tar.gz", backupDB, now.Format("2006-01-02"))
		cmd = exec.CommandContext(ctx, "tar", "czf", name, "/tmp/"+backupDB)
		err = cmd.Run()
		if err != nil {
			log.Println("tar", err)
			continue
		}
		log.Println("created " + name + " success, uploading to qiniu.")
		ch <- name
	}
}

func qiniuUpload(ch chan string) {
	for key := range ch {
		mac := qbox.NewMac(accessKey, secretKey)
		putPolicy := &storage.PutPolicy{
			Scope:      bucket,
			Expires:    3600,
			InsertOnly: 1,
		}
		upToken := putPolicy.UploadToken(mac)

		cfg := &storage.Config{
			Zone:     &storage.ZoneHuadong,
			UseHTTPS: true,
		}

		// uploader
		uploader := storage.NewFormUploader(cfg)
		ret := new(storage.PutRet)
		putExtra := &storage.PutExtra{}

		err := uploader.PutFile(nil, ret, upToken, key, key, putExtra)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("uploaded %s success!\n", key)

		// delete file
		bucketManager := storage.NewBucketManager(mac, cfg)
		err = bucketManager.DeleteAfterDays(bucket, key, long)
		if err != nil {
			log.Println(err)
		}
	}
}

// parse duration
func parseDuration(d string) (time.Duration, error) {
	if len(d) == 0 {
		return 0, errors.New("incorrect duration input.")
	}

	length := len(d)
	switch d[length-1] {
	case 's', 'm', 'h':
		return time.ParseDuration(d)
	case 'd':
		di, err := strconv.Atoi(d[:length-1])
		if err != nil {
			return 0, err
		}
		return time.Duration(di) * time.Hour * 24, nil
	}

	return 0, errors.New("unsupported duration.")
}
