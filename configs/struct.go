package configs

type Model struct {
	App      App      `yaml:"App"`
	Database Database `yaml:"Database"`
}

type App struct {
	AdminMain string  `yaml:"AdminMain"` //管理员面板入口
	Host      string  `yaml:"Host"`      //绑定ip
	Port      int     `yaml:"Port"`      //绑定端口
	SecretKey string  `yaml:"SecretKey"` //app密钥，请保存不要被泄漏
	Captcha   Captcha `yaml:"Captcha"`   //验证码类型
	DS        CS      `yaml:"CS"`        //安全请求
}

type Captcha struct {
	Enable bool   `yaml:"Enable"`
	Type   string `yaml:"Type"`
}

type CS struct {
	Enable bool   `yaml:"Enable"`
	Key    string `yaml:"Key"`
}

type Database struct {
	Redis Redis `yaml:"Redis"`
	Mysql Mysql `yaml:"Mysql"`
}

type Redis struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Password string `yaml:"Password"`
}

type Mysql struct {
	Enable     bool   `yaml:"Enable"`
	Host       string `yaml:"Host"`
	Port       int    `yaml:"Port"`
	DataName   string `yaml:"DataName"`
	UserName   string `yaml:"UserName"`
	Password   string `yaml:"Password"`
	Charset    string `yaml:"Charset"`
	RedisCache bool   `yaml:"RedisCache"`
}
