package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


var (
	vers bool
	ossProvider string
	ossEndpoint string
	aliAccessKey string
	aliSecretKey string
)
//cobra第三方库编写cmd方式
//创建一个总root节点
var RootCmd = &cobra.Command{
	//节点名称
	Use: "cloud-station-cli",
	//节点短描述
	Short: "cloud-station-cli 文件中转服务",
	//节点长描述
	Long: `cloud-station-cli ...`,
	//需要在该节点运行的内容
	RunE: func(cmd *cobra.Command, args []string) error {
		//如果vers为true则打印版本信息
		if vers {
			fmt.Println("version: 1.0.0")
			return nil
		}
		return fmt.Errorf("no flafs find")
	},
}

//cmd启动器
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

//命令行参数，赋值的定义和初始化
func init() {
	//全局标志PersistentFlags 所有的子节点都可以使用  Flags仅本节点可用
	//第一个形参是cmd 传参传给哪个变量接受，第二个代表 该 参数的名称，第三个代表 参数指令，第四个 代表参数的默认值，第五个代表参数的帮助描述
	RootCmd.PersistentFlags().StringVarP(&ossProvider, "oss_provider", "p", "aliyun", "the oss provider [aliyun/qcloud]")
	// RootCmd.PersistentFlags().StringVarP(&aliAccessKey, "ali_access_id", "i", "", "the ali oss access id")
	// RootCmd.PersistentFlags().StringVarP(&aliSecretKey, "ali_secret_key", "k", "", "the ali oss access key")
	RootCmd.Flags().BoolVarP(&vers, "version", "v", false, "the cloud-station-cli version")
}
