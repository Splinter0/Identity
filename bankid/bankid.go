package bankid

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const LATEST_VERSION = "6.0"

type BankIDEnvironment string

const (
	PRODUCTION BankIDEnvironment = "prod"
	TEST       BankIDEnvironment = "test"
)

func (e BankIDEnvironment) getEndpoint() string {
	switch e {
	case PRODUCTION:
		return "https://appapi2.bankid.com"
	default:
		return "https://appapi2.test.bankid.com"
	}
}

type BankIDRP struct {
	Client *http.Client
	Config *BankIDConfig
}

type BankIDConfig struct {
	Env               BankIDEnvironment `yaml:"env"`
	Version           string            `yaml:"version"`
	CertificateFolder string            `yaml:"certificateFolder"`
	Domain            *string           `yaml:"domain"`
	VisibleMessage    string            `yaml:"visibleMessage"`
}

func NewBankIDRP(config *BankIDConfig) (*BankIDRP, error) {
	tlsConfig, err := buildTLSConfig(config)
	if err != nil {
		return nil, errors.New("Could not load TLS configuration for BankID: " + err.Error())
	}
	return &BankIDRP{
		Client: &http.Client{
			Transport: &http.Transport{TLSClientConfig: tlsConfig},
		},
		Config: config,
	}, nil
}

// URL Building

func (b *BankIDRP) buildUrl(path string) string {
	return b.Config.Env.getEndpoint() + "/rp/v" + b.Config.Version + path
}

func (b *BankIDRP) GetAuthUrl() string {
	return b.buildUrl("/auth")
}

func (b *BankIDRP) GetCollectUrl() string {
	return b.buildUrl("/collect")
}

func (b *BankIDRP) GetPhoneAuthUrl() string {
	return b.buildUrl("/phone/auth")
}

func (b *BankIDRP) GetCancelUrl() string {
	return b.buildUrl("/cancel")
}

// Authentication methods

func (b *BankIDRP) DoAuth(ar AuthRequest) *AuthResponse {
	var response AuthResponse
	b.post(b.GetAuthUrl(), &ar, &response)
	return &response
}

func (b *BankIDRP) DoPhoneAuth(par PhoneAuthRequest) *AuthResponse {
	var response AuthResponse
	b.post(b.GetPhoneAuthUrl(), &par, &response)
	return &response
}

// Launching

func (b *BankIDRP) GenerateLaunchURL(resp *AuthResponse, returnURL string, appLink bool) string {
	if appLink {
		return fmt.Sprintf("bankid:///?autostarttoken=%s&redirect=%s", resp.AutoStartToken, returnURL)
	}

	return fmt.Sprintf("https://app.bankid.com/?autostarttoken=%s&redirect=%s", resp.AutoStartToken, returnURL)
}

func (b *BankIDRP) GenerateQRData(resp *AuthResponse) (qrCodeData []string) {
	for count := 0; count < 30; count++ {
		hmac := hmac.New(sha256.New, []byte(resp.QrStartSecret))
		hmac.Write([]byte(strconv.Itoa(count)))
		qrAuthCode := hex.EncodeToString(hmac.Sum(nil))
		qrCodeData = append(qrCodeData, fmt.Sprintf("bankid.%s.%d.%s", resp.QrStartToken, count, qrAuthCode))
	}

	return qrCodeData
}

// Collecting order

func (b *BankIDRP) DoCollection(orderRef string) *CollectResponse {
	collectRequest := &CollectRequest{
		OrderRef: orderRef,
	}
	var response CollectResponse
	b.post(b.GetCollectUrl(), collectRequest, &response)
	return &response
}

// Cancelling

func (b *BankIDRP) Cancel(orderRef string) {
	cancelRequest := &CancelRequest{
		OrderRef: orderRef,
	}
	b.post(b.GetCancelUrl(), cancelRequest, nil)
}

// Utils

func (b *BankIDRP) GetCertPolicyString(policy CertificatePolicy) string {
	if b.Config.Env == TEST {
		return policy.getTest()
	}

	return string(policy)
}

func (b *BankIDRP) post(url string, request, response interface{}) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := b.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if response == nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, response)
	if err != nil {
		log.Fatal(err)
	}
}

func buildTLSConfig(config *BankIDConfig) (*tls.Config, error) {
	path := fmt.Sprintf("%s%s/", config.CertificateFolder, config.Env)
	cert, err := tls.LoadX509KeyPair(
		path+"cert.pem",
		path+"key.pem",
	)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(path + "ca-cert.pem")
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}, nil
}
