package entry

import (
	"fmt"

	api2captcha "github.com/2captcha/2captcha-go"
)

const (
	SITE_KEY = "6Lcn_iATAAAAABvxzR18hrQhMZh_A6b4hPrZQKv2"
)

func (s *Session) createCaptchaSolver() error {

	s.Log.Info("solving captcha...")

	TwoClient := api2captcha.NewClient(s.UserConfig.TwoCapKey)

	cap := api2captcha.ReCaptcha{
		SiteKey:   SITE_KEY,
		Url:       "https://www.motosport.com/win",
		Invisible: false,
		Action:    "verify",
	}

	req := cap.ToRequest()
	code, _, err := TwoClient.Solve(req)
	if err != nil {
		return fmt.Errorf("error solving captcha: %w", err)
	}

	s.state.ReCap = code
	// s.Log.Info("captcha solved: \n", s.state.ReCap)
	return nil

}
