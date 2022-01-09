package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	//初始化oss配置参数
	conf     = NewConf()
	//CLI参数
	filename string
	help bool
)

//声明CLI参数
func init () {
	flag.StringVar(&filename,"f","","请输入需要上传文件的路径")
	flag.BoolVar(&help,"h",false,"打印本工具的使用说明")
}

//默认帮助信息
func usage(){
	fmt.Fprintf(os.Stderr,`cloud-station version: 0.0.1
	Usage: cloud-station [-h] -f <uplaod_file_path>
	Options:
	`)
		flag.PrintDefaults()
}

//执行CLI逻辑
func LoadArgsFromCLI(){
	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}
}

//OSS配置构造函数,参数配置从环境变量获取
func NewConf() *Conf {
	return &Conf{
		endpint:    os.Getenv("ALI_OSS_ENDPOINT"),
		accessKey:  os.Getenv("ALI_AK"),
		sercetKey:  os.Getenv("ALI_SK"),
		bucketName: os.Getenv("ALI_OSS_BUCKETNAME"),
	}
}

type Conf struct {
	endpint    string
	accessKey  string
	sercetKey  string
	bucketName string
}

//上传文件
func (c Conf) upload(filePath string) (downloadURL string, err error) {
	client, err := oss.New(c.endpint, c.accessKey, c.sercetKey)
	if err != nil {
		err = fmt.Errorf("new client error,%s", err)
		return
	}
	bucket, err := client.Bucket(c.bucketName)
	if err != nil {
		err = fmt.Errorf("get bucket error,%s", err)
		return
	}
	err = bucket.PutObjectFromFile(filePath, filePath)
	if err != nil {
		err = fmt.Errorf("upload bucket %s error", err)
		return
	}

	return bucket.SignURL(filePath, oss.HTTPGet, 60*60*24*3)
}

//配置字段校验
func (c Conf) validate() error {
	if c.endpint == "" {
		return fmt.Errorf("endpint missed")
	} else if c.accessKey == "" || c.sercetKey == "" {
		return fmt.Errorf("access key or secret key missed")
	} else if c.bucketName == "" {
		return fmt.Errorf("bucket name missed")
	}
	return nil
}

func main() {
	if err := conf.validate(); err != nil {
		fmt.Printf("parameter is nil %s", err)
		os.Exit(1)
	}

	LoadArgsFromCLI()

	downloadURL, err := conf.upload(filename)
	if err != nil {
		fmt.Printf("upload file error %s", err)
		os.Exit(2)
	}

	fmt.Printf("文件:%s 上传成功", filename)

	fmt.Printf("下载连接: %s\n", downloadURL)
	fmt.Println("\n注意: 文件下载有效期为1天, 中转站保存时间为3天, 请及时下载")
	os.Exit(0)
}
