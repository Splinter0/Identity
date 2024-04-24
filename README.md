# Identity

This project is a collection of implementations of eIDs and other identification/authentication providers that work on the premise of "something you have" and "something you are".

The focus of this repository is to provide a Golang implementation of various providers with a **focus on security**.

**Implemented APIs**
- Swedish BankID RP API

## Usage

The project can be used as a skeleton starter for your integration with the providers you need. It can be customized, styled and hooked into at your preference.
You can edit the `templates/` and `static/` files to fit your needs and then deploy this project as a microservice that can then be used to authenticate your users.

Under the hood a [gin](https://github.com/gin-gonic/gin) server is used to handle requests.

Alternatively you can use each provider such as: `github.com/Splinter0/identity/bankid` as a library in your Golang project to integrate the providers that way.

The service requires a configuration file `config.yml` where your application needs are specified, however this is only needed if you deploy this as a standalone service. If you are using this repository as a library you can specify the configuration programmatically.

```yml
service: "Company AB"
providers: ["bankid"]
bankid:
  env: "test"
  version: "6.0"
  certificateFolder: "bankid/certificates/"
  domain: "example.app"
  visibleMessage: "Log into an amazing company"
```

- The `service` key is a global key that defines the name of your application
- The `providers` key is used to defined which providers should be enabled for this deployment
- Each provider has its own configuration which can be defined by the provider's name (for example `bankid`) following the parameters required for that specific provider.

## Swedish BankID

The Swedish BankID RP API allows users to log in using the BankID app. By default the environment, set by the config key `env`, is set to `"test"`, this means that the test servers of BankID are being used.

The project ships with the test certificates needed to integrate with the test servers, in production you will have to request your own certificates from BankID. You can follow the guides [here](https://www.bankid.com/utvecklare/test) on how you can set up your BankID testing environment.

You can find a deployment of this at [http://bankid.mastersplinter.work](http://bankid.mastersplinter.work) where you can test it out provided you have configured your BankID device in test mode, you can learn how to do that [here](https://www.bankid.com/en/utvecklare/test/skaffa-testbankid/testbankid-konfiguration)

### Configuration

- `env` -> sets the current environment, can be set to `test` or `prod`
- `version` -> BankID API version to use
- `certificateFolder` -> where your certificates to communicate with BankID's API are stored
- `domain` -> the domain in which the app will run under
- `visibleMessage` -> the message your users will see when logging in with BankID

### Security Note

This repository is part of a security research project that specifically looks into the security of different eID solutions. Specifically in BankID's case, security features such as `certificate policies` and ip address checks have been implemented to serve as a guideline on how to securely implement this provider.

You can read more about this research [here](https://mastersplinter.work/research/bankid/)