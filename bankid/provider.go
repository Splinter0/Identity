package bankid

import (
	"encoding/base64"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

const SESSION_TIMEOUT = 30

type BankIDProvider struct {
	Client *BankIDRP
	Cache  *cache.Cache
}

type BankIDAuthenticationRequest struct {
	SameDevice     bool
	UserAgent      string
	UserIp         string
	MessageForUser string
	RedirectURL    string
}

type BankIDAuthenticationResponse struct {
	LaunchURL      string `json:"launchUrl,omitempty"`
	QrCodeData     string `json:"qrCodeData,omitempty"`
	Success        bool   `json:"success"`
	Message        string `json:"message,omitempty"`
	TransactionKey string `json:"-"`
}

type BankIDStatusResponse struct {
	Message string        `json:"message"`
	Status  CollectStatus `json:"status"`
	Data    interface{}   `json:"data,omitempty"`
}

type BankIDTransaction struct {
	SameDevice bool
	Mobile     bool
	UserIp     string
	QrCodeData []string
	OrderRef   string
	StartedAt  time.Time
}

func NewBankIDProvider(config *BankIDConfig) *BankIDProvider {
	rp, err := NewBankIDRP(config)
	if err != nil {
		log.Fatal(err)
	}
	c := cache.New((SESSION_TIMEOUT+1)*time.Second, 1*time.Minute)

	return &BankIDProvider{
		Client: rp,
		Cache:  c,
	}
}

func (provider *BankIDProvider) Authenticate(request BankIDAuthenticationRequest) BankIDAuthenticationResponse {
	var policy CertificatePolicy
	isMobile := IsMobileUserAgent(request.UserAgent)
	if isMobile || !request.SameDevice {
		policy = Mobile
	} else {
		policy = OnFile
	}
	rawRequest := AuthRequest{
		EndUserIp:             request.UserIp,
		UserVisibleDataFormat: "simpleMarkdownV1",
		UserVisibleData:       provider.buildUserVisibeData(request.MessageForUser),
		Requirement: &AuthRequestRequirements{
			CertificatePolicies: []string{
				provider.Client.GetCertPolicyString(
					policy,
				),
			},
		},
	}
	resp := provider.Client.DoAuth(rawRequest)
	if resp.ErrorCode != "" {
		return BankIDAuthenticationResponse{
			Success: false,
			Message: resp.Details,
		}
	}
	response := BankIDAuthenticationResponse{
		TransactionKey: md5sum(resp.OrderRef),
		Success:        true,
	}
	qrCodeData := []string{}
	if !request.SameDevice {
		qrCodeData = provider.Client.GenerateQRData(resp)
		response.QrCodeData = qrCodeData[0]
	} else {
		response.LaunchURL = provider.Client.GenerateLaunchURL(resp, request.RedirectURL, true)
	}
	provider.setTransaction(
		response.TransactionKey,
		BankIDTransaction{
			SameDevice: request.SameDevice,
			Mobile:     isMobile,
			UserIp:     request.UserIp,
			QrCodeData: qrCodeData,
			OrderRef:   resp.OrderRef,
			StartedAt:  time.Now(),
		},
	)

	return response
}

func (provider *BankIDProvider) Status(transactionKey string) BankIDStatusResponse {
	transaction, ok := provider.getTransaction(transactionKey)
	if !ok {
		return BankIDStatusResponse{
			Message: "Transaction not found",
			Status:  FAILED,
		}
	}

	collectedData := provider.Client.DoCollection(transaction.OrderRef)
	if collectedData.Status == FAILED {
		return BankIDStatusResponse{
			Message: collectedData.HintCode.GetMessage(),
			Status:  FAILED,
		}
	} else if collectedData.Status == PENDING {
		var data interface{}
		if !transaction.SameDevice {
			frame := int(time.Since(transaction.StartedAt).Seconds()) % SESSION_TIMEOUT
			data = map[string]string{
				"qrData": transaction.QrCodeData[frame],
			}
		}
		return BankIDStatusResponse{
			Message: collectedData.HintCode.GetMessage(),
			Status:  PENDING,
			Data:    data,
		}
	}

	// If for some reason we have a weird status
	if collectedData.Status != COMPLETE {
		return BankIDStatusResponse{
			Message: "Authentication failed",
			Status:  FAILED,
		}
	}

	if transaction.SameDevice && transaction.UserIp != collectedData.CompletionData.Device.IpAddress {
		return BankIDStatusResponse{
			Message: "BankID transaction was not completed using the same device",
			Status:  FAILED,
		}
	}

	return BankIDStatusResponse{
		Message: "Success!",
		Status:  COMPLETE,
		Data:    collectedData.CompletionData,
	}
}

func (provider *BankIDProvider) Cancel(transactionKey string) {
	transaction, ok := provider.getTransaction(transactionKey)
	if !ok {
		// Means it's already expired and it does not matter
		return
	}
	provider.Client.Cancel(transaction.OrderRef)
}

func (provider *BankIDProvider) buildUserVisibeData(message string) string {
	return base64.StdEncoding.EncodeToString([]byte(message))
}

func (provider *BankIDProvider) setTransaction(key string, transaction BankIDTransaction) {
	provider.Cache.Set(key, transaction, cache.DefaultExpiration)
}

func (provider *BankIDProvider) getTransaction(key string) (BankIDTransaction, bool) {
	value, found := provider.Cache.Get(key)
	if !found {
		return BankIDTransaction{}, false
	}

	transaction, ok := value.(BankIDTransaction)
	if !ok {
		return BankIDTransaction{}, false
	}

	return transaction, true
}
