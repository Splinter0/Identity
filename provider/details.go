package provider

type UserDetailType string

const (
	FIRST_NAME UserDetailType = "firstName"
	LAST_NAME  UserDetailType = "lastName"
	AGE        UserDetailType = "age"
	AGE_COMP   UserDetailType = "ageComparison"
	IP_ADDRESS UserDetailType = "ipAddress"
	DEVICE_ID  UserDetailType = "deviceId"
	ISSUE_DATE UserDetailType = "issueDate"
	COUNTRY_N  UserDetailType = "countryNumber"
)

type UserDetail struct {
	Type       UserDetailType
	Value      interface{}
	Requested  bool
	Allowed    bool
	Fullfilled bool
}

// TODO: description builder

type FirstName struct {
	UserDetail
	Value string
}
