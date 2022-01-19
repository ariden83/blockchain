package recaptcha

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"time"

	"github.com/ariden83/blockchain/cmd/web/config"
)

type Captcha struct {
	log    *zap.Logger
	cfg    config.ReCaptcha
	client *http.Client
}

func New(cfg config.ReCaptcha, log *zap.Logger) *Captcha {
	if cfg.SiteKey == "" || cfg.SecretKey == "" {
		return nil
	}

	return &Captcha{
		cfg: cfg,
		log: log.With(zap.String("service", "recaptcha")),
		client: &http.Client{
			Timeout: cfg.Timeout * time.Second,
		},
	}
}

type Input struct {
	// Required. The shared key between your site and reCAPTCHA.
	Secret string `json:"secret"`
	// Required. The user response token provided by the reCAPTCHA client-side integration on your site.
	Response string `json:"response"`
	// Optional. The user's IP address.
	RemoteIP string `json:"remoteip"`
}

type Output struct {
	Success     bool      `json:"success"`
	ChallengeTs time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	Error       []string  `json:"error-codes"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
}

// https://developers.google.com/recaptcha/docs/v3?authuser=1
// https://www.google.com/recaptcha/admin/site/505443814/setup
// https://www.google.com/recaptcha/admin/site/505443814 (analytics)
// https://www.google.com/recaptcha/admin/site/505443814/settings (settings)
func (c *Captcha) Verify(token, remoteIP string) (valid bool) {
	c.log.Debug("verifyCaptcha")
	URL, err := url.Parse(c.cfg.URL)
	if err != nil {
		c.log.Error("fail to parse URL", zap.Error(err))
		return
	}

	req, err := http.NewRequest(http.MethodPost, URL.String(), nil)
	if err != nil {
		c.log.Error("fail to set New Request to google API", zap.Error(err))
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	q := req.URL.Query()
	q.Add("secret", c.cfg.SecretKey)
	q.Add("response", token)
	q.Add("remoteip", remoteIP)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("fail to call google api", zap.Error(err))
		return
	}

	defer resp.Body.Close()

	var captchaResp Output
	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&captchaResp); err != nil {
		c.log.Error("fail to call google api", zap.Error(err))
		return
	}
	if captchaResp.Success {
		valid = true
		c.log.Info("recaptcha is valid")
	} else {
		c.log.Error(fmt.Sprintf("google captcha api return error: %+v", captchaResp.Error), zap.Error(err))
	}
	return
}
