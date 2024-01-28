# BankID weaknesses

- One can grab a valid autostarttoken from another client transaction, modify the redirect and then send it to someone. This is an open redirect. `bankid:///?autostarttoken=d1195a87-7a87-45ec-bac6-ec8e98a2776c&redirect=http://localhost`
- Also if I capture the link and I send it to someone else, what happens then? Detection is probably on the RP side. https://docs.swedenconnect.se/technical-framework/latest/12_-_BankID_Profile_for_the_Swedish_eID_Framework.html
- DoS via phone authentication?
- Stealing data from a compromised machine