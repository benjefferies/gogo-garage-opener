# gogo-garage-opener
Go implementation of a Raspberry Pi garage door opener

Also see [gogo-garage-opener-ui](https://github.com/benjefferies/gogo-garage-opener-ui) implemented using ionics framework

Also see [Compile for arm] (https://gist.github.com/steeve/6905542)

#### Guide
##### Create user
To create a user to use API's or login via the app [gogo-garage-opener-ui](https://github.com/benjefferies/gogo-garage-opener-ui) run the app with `--email` and `--password` arguments e.g.
`gogo-garage-opener --email benjefferies@example.com --password secret`

##### Use one time pin
To use a one time pin go to http://localhost:8080/user/one-time-pin/060cd65f-f700-4bf4-80df-ae8f78c38696. The UUID at the end is the generated pin, once the open button has been pressed the pin will be marked as used.


#### Features

* Open garage door using [gogo-garage-opener-ui](https://github.com/benjefferies/gogo-garage-opener-ui)
* Generate one time pins to allow someone temporary access to your garage i.e. A delivery man

#### Future

* Add tests
* Open garage door via location i.e. Automatically open garage within 100 metres of your garage