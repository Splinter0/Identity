package provider

import "github.com/gin-gonic/gin"

type ProviderCountry string

const SWEDEN ProviderCountry = "Sweden"

type AuthenticationRequest interface {
	GetRequestedUserDetails() []UserDetailType
	GetOrderToken() string
}

type AuthenticationResponse interface {
	GetUserDetail(UserDetailType) interface{}
	GetUserDetailMap() map[UserDetailType]interface{}
}

type Provider interface {
	Authenticate(AuthenticationRequest) AuthenticationResponse
	GetSetUpUserDetails() []UserDetailType
	GetProviderCountry() ProviderCountry
	GetName() string
	GetRoutes() map[string]func(*gin.Context)
}
