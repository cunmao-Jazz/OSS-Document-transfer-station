package store

//oss 存储适配器
//定义统一接口,降低耦合度
type OSSUploader interface {
	Upload(bucketName,objectKey,fileName string) (downloadUrl string,err error)
}