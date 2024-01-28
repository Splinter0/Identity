package provider

type ProviderCountry string

const SWEDEN ProviderCountry = "Sweden"

type AuthenticationRequest interface {
	GetRequestedUserDetails() []UserDetailType
}

type AuthenticationResponse interface {
	GetUserDetail(UserDetailType) interface{}
}

type Provider interface {
	Authenticate(AuthenticationRequest) AuthenticationResponse
	GetSetUpUserDetails() []UserDetailType
	GetProviderCountry() ProviderCountry
}
