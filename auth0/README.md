# Auth0

Auth0 is Authentication/Authorisation as a service and has a free tier offering which will be suffice for our needs.

## Setup
1. Register an account at https://auth0.com/
1. Create a new API https://auth0.com/docs/apis
   * Name: gogo-garage-opener
   * Identifier: https://<YOUR_GARAGE_DOOR_DOMAIN>/api
   * Signing Algorithm: RS256
1. Create a react native application https://auth0.com/docs/quickstart/native/react-native/00-login
   * Name: gogo-garage-opener-app
   * Type: Native
   * Allowed Callback URLs: uk.echosoft.garage.opener://<YOUR_AUTH0_DOMAIN>/android/uk.echosoft.garage.opener/callback

## App set up
1. Install https://play.google.com/store/apps/details?id=uk.echosoft.garage.opener
1. Update settings
   * Garage Opener Domain: This should be set to the RS/<YOUR_GARAGE_DOOR_DOMAIN>. The domain pointing to the Raspberry pi
   * Auth0 Domain: This should be set to the AS/<YOUR_AUTH0_DOMAIN>
   * Client ID: This should be set to the client ID found in Auth0 under your react native application
   * Auth0 API Audience: This should be set to the API Audience found on the API page in Auth0

## Rules
This rule allows you to use a social login for convenience but still lock down to specific users.
1. Go to Rules
1. Create Rule
1. Empty Rule
1. Copy and paste [allowed-emails.js](allowed-emails.js) into the code block and save
