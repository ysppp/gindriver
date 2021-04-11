package config

var Config = struct {
	AppName string

	DB struct {
		Name  string
		User  string
		Pass  string
		Host  string
		Port  uint
		Param string
	}

	Redis struct {
		Host  string
		Port  uint
		Index int
	}

	Oss struct {
		AccessKeyId     string
		AccessKeySecret string
		EndPoint        string
		BucketName      string
	}

	ListenAddr     string
	RPID           string
	RPOrigin       string
	RPDisplayName  string
	UploadLocation string
}{}
