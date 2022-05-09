package webhook

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Webhook struct {
	client *resty.Client
	url    string
}

func New(url string) *Webhook {
	return &Webhook{
		client: resty.New(),
		url:    url,
	}
}

func (w Webhook) SetToken(token string) {
	w.client.SetAuthToken(token)
}

func (w Webhook) SetBasicAuth(username, password string) {
	w.client.SetBasicAuth(username, password)
}

func (w Webhook) SetRootCertificate(ca string) {
	w.client.SetRootCertificateFromString(ca)
}

func (w Webhook) SetClientCertificate(certs ...tls.Certificate) {
	w.client.SetCertificates(certs...)
}

func (w Webhook) Send(ctx context.Context, subject, message string) error {
	resp, err := w.client.R().
		SetBody([]byte(message)).
		SetHeader("Content-Type", "application/json").
		Post(w.url)
	if err != nil {
		return err
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("unexpect response status code %d, response: %v", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}
