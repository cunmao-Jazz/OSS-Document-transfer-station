package cmd

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"gitee.com/wenqirui5600129/go-share-examples/cloud-station/pro/store"
	"gitee.com/wenqirui5600129/go-share-examples/cloud-station/pro/store/aliyun"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var (
	fileName   string
	bucketName string
)

//定义cmd子节点
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "上传文件到中转站",
	Long:  `上传文件到中转站`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//将密钥进行加密输入
		getSecretKeyFromInput()
		//实例化oss对象
		p, err := getProvider()
		if err != nil {
			return err
		}
		//如果传输文件为空则报错
		if fileName == "" {
			return fmt.Errorf("file name required")
		}
		//获取当前时间的时间戳
		day := time.Now().Format("20060102")
		//获取当前系统的主机名
		hn, err := os.Hostname()
		//如果os获取当前主机名失败,则通过请求服务的返回信息截取
		if err != nil {
			ipAddr := getOutBindIp()
			//如果返回为空则标记未知
			if ipAddr == "" {
				hn = "unknown"
			} else {
				hn = ipAddr
			}
		}

		//将时间戳,主机名,文件名进行拼接
		objecyKey := fmt.Sprintf("%s/%s/%s", day, hn, fileName)
		//调用upload方法进行文件上传
		downloadURL, err := p.Upload(bucketName, objecyKey, fileName)
		if err != nil {
			return err
		}
		//打印文件上传成功
		fmt.Printf("文件:%s上传成功\n", fileName)

		// 打印下载链接
		fmt.Printf("下载链接: %s\n", downloadURL)
		fmt.Println("\n注意: 文件下载有效期为1天, 中转站保存时间为3天, 请及时下载")

		return nil
	},
}

func getOutBindIp() string {
	//测试网络通信
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	//以冒号为分隔符，返回第一个字段
	addr := strings.Split(conn.LocalAddr().String(), ":")
	if len(addr) == 0 {
		return ""
	}

	return addr[0]
}

func getSecretKeyFromInput(){
	//将输入accessKey加密
	access := &survey.Password{
		Message : "请输入access key",
	}
	//将输入secretKey加密
	survey.AskOne(access,&aliAccessKey)
	secret := &survey.Password{
		Message : "请输入secret key",
	}
	survey.AskOne(secret,&aliSecretKey)
}

func getProvider() (p store.OSSUploader, err error) {
	//判断输入 是什么云厂商，如果是阿里云，则创建oss实例进行返回，如果是其他则返回nil和error
	switch ossProvider {
	case "aliyun":
		p, err = aliyun.NewUploader(ossEndpoint, aliAccessKey, aliSecretKey)
		return
	case "qcloud":
		return nil, fmt.Errorf("not impl")
	default:
		return nil, fmt.Errorf("unknown oss privier options [aliyun/qcloud]")
	}
}

func init() {
	uploadCmd.PersistentFlags().StringVarP(&bucketName, "bucket_name", "b", "devcloud-stations", "存储桶名称")
	uploadCmd.PersistentFlags().StringVarP(&fileName, "file_name", "f", "", "上传文件的名称")
	uploadCmd.PersistentFlags().StringVarP(&ossEndpoint, "bucket_endpoint", "e", "", "upload oss endpoint")
	//将子节点拼接到父节点后
	RootCmd.AddCommand(uploadCmd)
}
