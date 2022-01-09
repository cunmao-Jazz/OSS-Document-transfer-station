package aliyun

import (
	"fmt"

	"gitee.com/wenqirui5600129/go-share-examples/cloud-station/pro/store"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/go-playground/validator/v10"
)

var (
	//字段校验器
	validate = validator.New()
)

//构造函数传入属性值，并调用Validate进行属性值的校验,校验通过返回接口
//(该结构体实现了Upload方法，等于实现了OSSUploader接口,构造函数返回接口更加规范标准,权限分配也更加细化，返回这个接口，他就只调用接口定义的upload方法)
func NewUploader(eq, ak, sk string) (store.OSSUploader, error) {
	upload := &aliyun{
		EQ: eq,
		Ak: ak,
		SK: sk,
		listner: NewOssProgressListener(),
	}
	if err := upload.Validate(); err != nil {
		return nil, fmt.Errorf("validate params error,%s", err)
	}
	return upload, nil
}

//定义阿里云的机构体
type aliyun struct {
	//将oos地址 accessKey secretKey定义为结构体属性,并打上tag,对属性值进行校验
	EQ string `validate:"required"`
	Ak string `validate:"required"`
	SK string `validate:"required"`
	listner oss.ProgressListener
}

//定义字段校验方法，对传入的属性值进行校验
func (a *aliyun) Validate() error {
	return validate.Struct(a)
}
//将阿里云的上传过程封装成结构体方法
func (a *aliyun) Upload(bucketName, objectKey, fileName string) (downloadUrl string, err error) {
	//创建oss对象,传入oss地址,和accesskey和secretkey
	client, err := oss.New(a.EQ, a.Ak, a.SK)
	if err != nil {
		err = fmt.Errorf("new client error,%s", err)
		return
	}
	//于bucket建立连接
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		err = fmt.Errorf("get bucket error,%s", err)
		return
	}
	//上传文件 指定上传名称，和文件名称
	err = bucket.PutObjectFromFile(objectKey, fileName,oss.Progress(a.listner))
	if err != nil {
		err = fmt.Errorf("upload bucket %s error", err)
		return
	}
	//生成文件下载连接，以http get请求方式，三天后过期
	return bucket.SignURL(fileName, oss.HTTPGet, 60*60*24*3)
}
