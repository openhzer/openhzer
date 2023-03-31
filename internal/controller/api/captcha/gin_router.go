package captcha

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"hzer/configs"
	"hzer/internal/response"
	"hzer/pkg/logger"
)

//configJsonBody json request body.
type configJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

var store = base64Captcha.DefaultMemStore

//getCaptcha
// @Summary 获取验证码
// @Description 用于预防人机的验证码
// @Accept json
// @Tags public
// @Produce  json
// @Router /api/captcha [get]
// @Success 200 {object} response.Message{Data=captcha.reqCaptcha}
func getCaptcha(c *gin.Context) {
	// 	{
	// 		ShowLineOptions: [],
	// 		CaptchaType: "string",
	// 		Id: '',
	// 		VerifyValue: '',
	// 		DriverAudio: {
	// 			Length: 6,
	// 			Language: 'zh'
	// 		},
	// 		DriverString: {
	// 			Height: 60,
	// 			Width: 240,
	// 			ShowLineOptions: 0,
	// 			NoiseCount: 0,
	// 			Source: "1234567890qwertyuioplkjhgfdsazxcvbnm",
	// 			Length: 6,
	// 			Fonts: ["wqy-microhei.ttc"],
	// 			BgColor: {R: 0, G: 0, B: 0, A: 0},
	// 		},
	// 		DriverMath: {
	// 			Height: 60,
	// 			Width: 240,
	// 			ShowLineOptions: 0,
	// 			NoiseCount: 0,
	// 			Length: 6,
	// 			Fonts: ["wqy-microhei.ttc"],
	// 			BgColor: {R: 0, G: 0, B: 0, A: 0},
	// 		},
	// 		DriverChinese: {
	// 			Height: 60,
	// 			Width: 320,
	// 			ShowLineOptions: 0,
	// 			NoiseCount: 0,
	// 			Source: "设想,你在,处理,消费者,的音,频输,出音,频可,能无,论什,么都,没有,任何,输出,或者,它可,能是,单声道,立体声,或是,环绕立,体声的,,不想要,的值",
	// 			Length: 2,
	// 			Fonts: ["wqy-microhei.ttc"],
	// 			BgColor: {R: 125, G: 125, B: 0, A: 118},
	// 		},
	// 		DriverDigit: {
	// 			Height: 80,
	// 			Width: 240,
	// 			Length: 5,
	// 			MaxSkew: 0.7,
	// 			DotCount: 80
	// 		}
	// 	},
	// 	blob: "",
	// 	loading: false
	// }
	param := configJsonBody{
		Id:          "",
		CaptchaType: configs.Data.App.Captcha.Type,
		VerifyValue: "",
		DriverAudio: &base64Captcha.DriverAudio{},
		DriverString: &base64Captcha.DriverString{
			Length:          4,
			Height:          40,
			Width:           120,
			ShowLineOptions: 2,
			NoiseCount:      0,
			Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		},
		DriverChinese: &base64Captcha.DriverChinese{},
		DriverMath:    &base64Captcha.DriverMath{},
		DriverDigit:   &base64Captcha.DriverDigit{},
	}
	var driver base64Captcha.Driver

	//create base64 encoding captcha
	switch param.CaptchaType {
	case "audio":
		driver = param.DriverAudio
	case "string":
		driver = param.DriverString.ConvertFonts()
	case "math":
		driver = param.DriverMath.ConvertFonts()
	case "chinese":
		driver = param.DriverChinese.ConvertFonts()
	default:
		driver = param.DriverDigit
	}
	ca := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := ca.Generate()
	if err != nil {
		response.FailJson(c, response.CaptchaError, false)
		logger.Error(err)
		return
	}
	req := reqCaptcha{
		Data: b64s,
		Id:   id,
	}
	response.SuccessJson(c, "验证码请尽快使用", req)
}

func VerifyCaptcha(id string, code string) bool {
	return store.Verify(id, code, true)
}
