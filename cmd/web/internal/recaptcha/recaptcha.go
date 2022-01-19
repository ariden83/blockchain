package recaptcha

import (
	"bytes"
	"fmt"
	"github.com/ariden83/blockchain/cmd/web/config"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Captcha struct {
	log    *zap.Logger
	cfg    *config.ReCaptcha
	client *http.Client
}

func New(cfg *config.ReCaptcha, log *zap.Logger) *Captcha {
	return &Captcha{
		cfg: cfg,
		log: log.With(zap.String("service", "recaptcha")),
		client: &http.Client{
			Timeout: cfg.Timeout * time.Second,
		},
	}
}

func (c *Captcha) verifyCaptcha(token string) (err error) {

	c.log.Debug("verifyCaptcha")

	var URL *url.URL
	URL, err = url.Parse(c.cfg.URL)
	if err != nil {
		c.log.Error("fail to parse URL", zap.Error(err))
		return
	}
	parameters := url.Values{}
	URL.RawQuery = parameters.Encode()

	data := url.Values{}
	data.Set("client_id", `Lazy Test`)

	req, err := http.NewRequest(http.MethodPost, URL.String(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		c.log.Error("fail to set New Request to google API", zap.Error(err))
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	if err != nil {
		c.log.Error("fail to call google api", zap.Error(err))
		return
	}
	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("fail to call google api", zap.Error(err))
		return
	}

	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.log.Error("fail to call google api", zap.Error(err))
		return
	}
	resp.Body.Close()
	if err != nil {
		c.log.Error("fail to call google api", zap.Error(err))
		return
	}
	fmt.Println(string(f))
	return nil
}
