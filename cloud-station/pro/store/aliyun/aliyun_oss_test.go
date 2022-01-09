package aliyun_test

import (
	"os"
	"testing"
	"gitee.com/wenqirui5600129/go-share-examples/cloud-station/pro/store/aliyun"
	"github.com/stretchr/testify/assert"
)

var (
	bucketName string
	ak         string
	sk         string
	eq         string
)

//测试用例编写
func TestUploadFile(t *testing.T) {
	should := assert.New(t)
	uploader, err := aliyun.NewUploader(eq, ak, sk)
	//调用assert库对返回值进行校验
	//没有错误则判断通过
	if should.NoError(err) {
		downloadURL, err := uploader.Upload(bucketName, "1.txt", "1.txt")
		//没有返回error且downloadURL不为空则通过
		if should.NoError(err) {
			should.NotEmpty(downloadURL)
		}
	}

}

func init() {
	//将测试用的值通过环境变量的方式注入到变量中
	eq = os.Getenv("ALI_OSS_ENDPOINT")
	ak = os.Getenv("ALI_AK")
	sk = os.Getenv("ALI_SK")
	bucketName = os.Getenv("ALI_OSS_BUCKETNAME")
}
