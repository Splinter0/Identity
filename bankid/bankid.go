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
	"time"

	"github.com/skip2/go-qrcode"
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
		return "https://appapi2.test.bankid.com/"
	}
}

type BankIDRP struct {
	Client      *http.Client
	Version     string
	Environment BankIDEnvironment
}

func NewBankIDRP(env BankIDEnvironment) (*BankIDRP, error) {
	tlsConfig, err := buildTLSConfig(env)
	if err != nil {
		return nil, errors.New("Could not load TLS configuration for BankID: " + err.Error())
	}
	return &BankIDRP{
		Client: &http.Client{
			Transport: &http.Transport{TLSClientConfig: tlsConfig},
		},
		Version:     LATEST_VERSION,
		Environment: env,
	}, nil
}

// URL Building

func (b *BankIDRP) buildUrl(path string) string {
	return b.Environment.getEndpoint() + "/rp/v" + b.Version + path
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

func (b *BankIDRP) GenerateQR(resp *AuthResponse, comms chan []byte) {
	for count := 0; count < 30; count++ {
		hmac := hmac.New(sha256.New, []byte(resp.QrStartSecret))
		hmac.Write([]byte(strconv.Itoa(count)))
		qrAuthCode := hex.EncodeToString(hmac.Sum(nil))
		content := fmt.Sprintf("bankid.%s.%d.%s", resp.QrStartToken, count, qrAuthCode)
		log.Println(content)
		code, err := qrcode.Encode(content, qrcode.Medium, 256)
		if err != nil {
			log.Fatal(err)
		}

		comms <- code

		time.Sleep(1 * time.Second)
	}
}

// Collecting order

func (b *BankIDRP) StartCollection(resp *AuthResponse, messagesChannel chan string) *CollectResponse {
	collectRequest := &CollectRequest{
		OrderRef: resp.OrderRef,
	}
	jsonData, err := json.Marshal(collectRequest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
	req, err := http.NewRequest("POST", b.GetCollectUrl(), bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	for {
		collectResponse := b.collect(req)
		if collectResponse.Status == FAILED || collectResponse.Status == COMPLETE {
			return collectResponse
		}
		messagesChannel <- collectResponse.HintCode.GetMessage()
		time.Sleep(2 * time.Second)
	}
}

func (b *BankIDRP) collect(req *http.Request) *CollectResponse {
	resp, err := b.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var response CollectResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}
	return &response
}

// Utils

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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, response)
	if err != nil {
		log.Fatal(err)
	}
}

func buildTLSConfig(env BankIDEnvironment) (*tls.Config, error) {
	path := fmt.Sprintf("bankid/certificates/%s/", env)
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
