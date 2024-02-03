package bankid

// https://www.bankid.com/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/auth

// Requests

type AuthRequest struct {
	EndUserIp             string                   `json:"endUserIp"`
	UserVisibleData       string                   `json:"userVisibleData,omitempty"`
	UserNonVisibleData    string                   `json:"userNonVisibleData,omitempty"`
	UserVisibleDataFormat string                   `json:"userVisibleDataFormat,omitempty"`
	Requirement           *AuthRequestRequirements `json:"requirement,omitempty"`
}

type CardReaderClass string

const (
	ComputerOrReader CardReaderClass = "class1" // default
	OnlyReader       CardReaderClass = "class2"
)

type CertificatePolicy string

const (
	OnFile    CertificatePolicy = "1.2.752.78.1.1"
	SmartCard CertificatePolicy = "1.2.752.78.1.2"
	Mobile    CertificatePolicy = "1.2.752.78.1.5"
)

func (c CertificatePolicy) getTest() string {
	switch c {
	case OnFile:
		return "1.2.3.4.5"
	case SmartCard:
		return "1.2.3.4.10"
	case Mobile:
		return "1.2.3.4.25"
	}
	return "1.2.752.60.1.6"
}

type AuthRequestRequirements struct {
	PinCode             bool            `json:"pinCode,omitempty"`
	Mrtd                bool            `json:"mrtd,omitempty"`
	CardReader          CardReaderClass `json:"cardReader,omitempty"`
	PersonalNumber      string          `json:"personalNumber,omitempty"`
	CertificatePolicies []string        `json:"certificatePolicies,omitempty"`
}

type CallInitiator string

const (
	UserCallInitiator CallInitiator = "user"
	RPCallInitiator   CallInitiator = "RP"
)

type PhoneAuthRequest struct {
	PersonalNumber        string                   `json:"personalNumber"`
	CallInitiator         CallInitiator            `json:"callInitiator"`
	UserVisibleData       string                   `json:"userVisibleData,omitempty"`
	UserNonVisibleData    string                   `json:"userNonVisibleData,omitempty"`
	UserVisibleDataFormat string                   `json:"userVisibleDataFormat,omitempty"`
	Requirements          *AuthRequestRequirements `json:"requirement,omitempty"`
}

type OrderRequest struct {
	OrderRef string `json:"orderRef"`
}

type CollectRequest OrderRequest

type CancelRequest OrderRequest

// Responses

type BankIDResponse struct {
	ErrorCode string `json:"errorCode"`
	Details   string `json:"details"`
}

type AuthResponse struct {
	BankIDResponse
	OrderRef       string `json:"orderRef"`
	AutoStartToken string `json:"autoStartToken"`
	QrStartToken   string `json:"qrStartToken"`
	QrStartSecret  string `json:"qrStartSecret"`
}

type CollectStatus string

const (
	PENDING  CollectStatus = "pending"
	COMPLETE CollectStatus = "complete"
	FAILED   CollectStatus = "failed"
)

type CollectResponse struct {
	BankIDResponse
	CollectRequest
	Status         CollectStatus         `json:"status"`
	HintCode       HintCode              `json:"hintCode"`
	CompletionData CollectCompletionData `json:"completionData"`
}

type CollectCompletionData struct {
	User            CompletionDataUser   `json:"user"`
	Device          CompletionDataDevice `json:"device"`
	BankIdIssueDate string               `json:"bankIdIssueDate"`
	StepUp          CompletionDataStepUp `json:"stepUp"`
	Signature       string               `json:"signature"`
	OcspResponse    string               `json:"ocspResponse"`
}

type CompletionDataUser struct {
	PersonalNumber string `json:"personalNumber"`
	Name           string `json:"name"`
	GivenName      string `json:"givenName"`
	Surname        string `json:"surname"`
}

type CompletionDataDevice struct {
	IpAddress string `json:"ipAddress"`
	UHI       string `json:"uhi"`
}

type CompletionDataStepUp struct {
	Mrtd bool `json:"mrtd"`
}
