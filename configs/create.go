package configs

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v2"
	"hzer/pkg/integral"
	"hzer/pkg/jwt"
	"hzer/pkg/util"
	"io"
	"os"
	"strings"
)

var (
	//go:embed _config_template.yml
	configTemplate []byte
	Data           Model
	Env            string
)

func InitConfigs() {
	//获取启动目录
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for i, v := range os.Args {
		if v == "-env" {
			if len(os.Args)-1 >= i+1 {
				Env = os.Args[i+1]
			}
		}
	}
	confpath := fmt.Sprintf("%s/%s_config.yml", path, Env)
	fileis, err := integral.PathExists(confpath)
	if err != nil {
		panic(err)
	}
	if !fileis {
		fmt.Println("环境不存在,是否创建新环境配置？(Y/N 默认:Y)：")
		for {
			var ifs string
			fmt.Scanf("%s", &ifs)
			switch strings.ToLower(ifs) {
			case "n":
				break
			case "y":
				err = os.WriteFile(Env+"_config.yml", configTemplate, 0644)
				if err != nil {
					fmt.Printf("配置文件生成失败, Error: %s", err.Error())
					os.Exit(2)
				}
				fmt.Println("配置文件已生成,请修改配置文件后重新运行")
				break
			default:
				fmt.Println("配置文件已生成,请修改配置文件后重新运行")
				break
			}
			os.Exit(0)
		}

	}
	hFile, err := os.Open(confpath)
	defer hFile.Close()
	if err != nil {
		panic(err)
	}
	bFile, err := io.ReadAll(hFile)
	err = yaml.Unmarshal(bFile, &Data)
	if err != nil {
		panic(err)
	}
	if Data.App.SecretKey != "" {
		os.Setenv("HZER_JWT_SECRET_KEY", Data.App.SecretKey)
		jwt.SecretKey = Data.App.SecretKey
	} else {
		jwt.SecretKey = util.Ifs(jwt.SecretKey == "", util.RandomStr(32), jwt.SecretKey)
	}
}
