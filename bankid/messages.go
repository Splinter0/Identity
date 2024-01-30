package bankid

var messages = map[string]string{
	"RFA1":  "Start your BankID app.",
	"RFA3":  "Action cancelled. Please try again.",
	"RFA4":  "An identification or signing for this personal number is already started. Please try again.",
	"RFA5":  "Internal error. Please try again.",
	"RFA6":  "Action cancelled.",
	"RFA8":  "The BankID app is not responding. Please check that it's started and that you have internet access. If you don't have a valid BankID you can get one from your bank. Try again.",
	"RFA9":  "Enter your security code in the BankID app and select Identify or Sign.",
	"RFA13": "Trying to start your BankID app.",
	"RFA14": "Searching for BankID, it may take a little while... If a few seconds have passed and still no BankID has been found, you probably don't have a BankID which can be used for this identification/signing on this computer. If you have a BankID card, please insert it into your card reader. If you don't have a BankID you can get one from your bank. If you have a BankID on another device you can start the BankID app on that device.",
	"RFA15": "Searching for BankID:s, it may take a little while... If a few seconds have passed and still no BankID has been found, you probably don't have a BankID which can be used for this identification/signing on this computer. If you have a BankID card, please insert it into your card reader. If you don't have a BankID you can get one from your bank.",
	"RFA16": "The BankID you are trying to use is blocked or too old. Please use another BankID or get a new one from your bank.",
	"RFA17": "The BankID app couldn't be found on your computer or mobile device. Please install it and get a BankID from your bank. Install the app from your app store or https://install.bankid.com.",
	"RFA18": "Start the BankID app.",
	"RFA19": "Would you like to identify yourself or sign with a BankID on this computer, or with a Mobile BankID?",
	"RFA20": "Would you like to identify yourself or sign with a BankID on this computer, or with a BankID on another device?",
	"RFA21": "Identification or signing in progress.",
	"RFA22": "Unknown error. Please try again.",
	"RFA23": "Process your machine-readable travel document using the BankID app.",
}

type HintCode string

const (
	UNKNOWN                 HintCode = "unknown"
	OUTSTANDING_TRANSACTION HintCode = "outstandingTransaction"
	NO_CLIENT               HintCode = "noClient"
	STARTED                 HintCode = "started"
	USER_MRTD               HintCode = "userMrtd"
	USER_CALL_CONFRIRM      HintCode = "userCallConfirm"
	USER_SIGN               HintCode = "userSign"
	EXPIRED_TRANSACTION     HintCode = "expiredTransaction"
	CERTIFICATE_ERR         HintCode = "certificateErr"
	USER_CANCEL             HintCode = "userCancel"
	CANCELLED               HintCode = "cancelled"
	START_FAILED            HintCode = "startFailed"
	USER_DECLINED_CALL      HintCode = "userDeclinedCall"
)

func (hc HintCode) GetMessage() string {
	switch hc {
	case OUTSTANDING_TRANSACTION:
		return messages["RFA13"]
	case NO_CLIENT:
		return messages["RFA1"]
	case STARTED:
		return messages["RFA15"]
	case USER_MRTD:
		return messages["RFA23"]
	case USER_SIGN:
		return messages["RFA9"]
	case EXPIRED_TRANSACTION:
		return messages["RFA8"]
	case CERTIFICATE_ERR:
		return messages["RFA16"]
	case USER_CANCEL:
		return messages["RFA6"]
	case CANCELLED:
		return messages["RFA3"]
	case START_FAILED:
		return messages["RFA17"]
	default:
		return messages["RFA22"]
	}
}
