package provider

import (
	"encoding/base64"
	"fmt"

	"github.com/Splinter0/identity/bankid"
)

// Provider

type BankIDProvider struct {
	Client             *bankid.BankIDRP
	DefaultUserDetails []UserDetailType
}

func (provider *BankIDProvider) Authenticate(request BankIDAuthenticationRequest) {
	if !request.SameDevice {
		return // TODO: qr code
	}
	// TODO: build requirements
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
	fmt.Println(rawRequest.UserVisibleData)
	fmt.Println(resp)
	fmt.Println(provider.Client.GenerateLaunchURL(resp, "http://localhost", true))
	data := provider.Client.StartCollection(resp)
	fmt.Println(data)
}

func (Provider *BankIDProvider) buildUserVisibeData(details []UserDetailType, message string) string {
	// TODO: create general config system
	data := "# On behalf of Amazing AB\nAmazing AB would like to use BankID to access your following details"
	for _, detail := range details {
		data += fmt.Sprintf("\n+ *%s*", detail) // TODO: build using description
	}
	data += "\nMessage from Amazing AB:\n \"" + message + "\""

	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Request

type BankIDAuthenticationRequest struct {
	RequestedDetails []UserDetailType // Leave empty for defaults
	Mobile           bool
	SameDevice       bool
	UserIp           string
	MessageForUser   string
}

func (request *BankIDAuthenticationRequest) GetRequestedUserDetails() []UserDetailType {
	return request.RequestedDetails
}

// Response

type BankIDAuthenticationResponse struct {
	rawResponse   *bankid.AuthResponse
	userDetailMap map[UserDetailType]interface{}
}

func (response *BankIDAuthenticationResponse) GetUserDetail(udt UserDetailType) interface{} {
	return response.userDetailMap[udt]
}
