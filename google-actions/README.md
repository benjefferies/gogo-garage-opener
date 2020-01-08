# Google actions
This is a simple proxy for gogo-garage-opener to allow you to control your garage door as a smart device using google assistant.

## Setup
1. Create google action
  - Name the google action
1. Set up account linking
  - Set client_id, client_secret from [auth0](../auth0/README.md)
  - Set authorize and token URL from https://<AUTH0-DOMAIN>/.well-known/openid-configuration
1. Deploy proxy function
  - Install firebase
  - Login `firebase login`
  - Switch to google actions project `firebase use <GOOGLE_ACTIONS_PROJECT_ID>`
  - Deploy `firebase deploy`
1. Set up google home
  - Go to google home app on device
  - Add new device
  - Set up device
  - Click `Have something already set up?`
  - Search for your app and log in

Now you can say "ok google, is the garage open?"

## Troubleshooting
Firebase action errors with getaddrinfo EAI_AGAIN <DOMAIN> at GetAddrInfoReqWrap.onlookup [as oncomplete] (dns.js:67:26)

You'll need to change the firebase plan to pay as you go (blaze) to make external network calls.
