package provider

import (
	"encoding/base64"
	"fmt"

	"github.com/Splinter0/identity/bankid"
)

// Provider

type BankIDConfig struct {
	CompanyName     string
	RedirectBaseUrl string
}

type BankIDProvider struct {
	Client             *bankid.BankIDRP
	DefaultUserDetails []UserDetailType
	LaunchURLChannel   chan string
	QRCodeChannel      chan []byte
	MessageChannel     chan string
	ResponseChannel    chan BankIDAuthenticationResponse
	Config             BankIDConfig
}

func (provider *BankIDProvider) Authenticate(request BankIDAuthenticationRequest) {
	// TODO: build requirements with certificate authorities
	var details []UserDetailType
	if len(request.RequestedDetails) == 0 {
		details = provider.DefaultUserDetails
	} else {
		details = request.RequestedDetails
	}
	rawRequest := bankid.AuthRequest{
		EndUserIp:             request.UserIp,
		UserVisibleDataFormat: "simpleMarkdownV1",
		UserVisibleData:       provider.buildUserVisibeData(details, request.MessageForUser),
	}
	resp := provider.Client.DoAuth(rawRequest)
	if !request.SameDevice {
		go provider.Client.GenerateQR(resp, provider.QRCodeChannel)
	} else {
		provider.LaunchURLChannel <- provider.Client.GenerateLaunchURL(resp, provider.Config.RedirectBaseUrl, true)
	}
	data := provider.Client.StartCollection(resp, provider.MessageChannel)
	if data.Status == bankid.FAILED {
		provider.ResponseChannel <- BankIDAuthenticationResponse{
			Success: false,
			Message: data.HintCode.GetMessage(),
		}
		return
	}
	provider.ResponseChannel <- BankIDAuthenticationResponse{
		UserDetailMap: provider.buildUserDetailMap(details, data.CompletionData),
		Success:       true,
		Message:       "",
	}
	fmt.Println(data)
}

func (provider *BankIDProvider) GetName() string {
	return "bankid"
}

func (provider *BankIDProvider) buildUserVisibeData(details []UserDetailType, message string) string {
	// TODO: create general config system
	data := fmt.Sprintf(
		"# On behalf of %s\n%s would like to use BankID to access your following details",
		provider.Config.CompanyName,
		provider.Config.CompanyName,
	)
	for _, detail := range details {
		data += fmt.Sprintf("\n+ *%s*", detail) // TODO: build using description
	}
	data += fmt.Sprintf(
		"\nMessage from %s:\n \"%s\"",
		provider.Config.CompanyName,
		message,
	)

	return base64.StdEncoding.EncodeToString([]byte(data))
}

func (provider *BankIDProvider) buildUserDetailMap(details []UserDetailType, completionData bankid.CollectCompletionData) map[UserDetailType]interface{} {
	userDetailMap := make(map[UserDetailType]interface{})
	for _, d := range details {
		var value interface{}
		switch d {
		case FIRST_NAME:
			value = completionData.User.GivenName
		case LAST_NAME:
			value = completionData.User.Surname
		default:
			continue
		}
		userDetailMap[d] = value
	}

	return userDetailMap
}

// Routes

/*
func (provider *BankIDProvider) GetRoutes() map[string]func(*gin.Context) {

}

func sameDeviceAction(c *gin.Context) {

*/

// Request

type BankIDAuthenticationRequest struct {
	RequestedDetails []UserDetailType // Leave empty for defaults
	OrderToken       string
	SameDevice       bool
	UserIp           string
	MessageForUser   string
}

func (request *BankIDAuthenticationRequest) GetRequestedUserDetails() []UserDetailType {
	return request.RequestedDetails
}

// Response

type BankIDAuthenticationResponse struct {
	UserDetailMap map[UserDetailType]interface{}
	Success       bool
	Message       string
}

func (response *BankIDAuthenticationResponse) GetUserDetail(udt UserDetailType) interface{} {
	return response.UserDetailMap[udt]
}
