[![Build Status](https://travis-ci.org/benjefferies/gogo-garage-opener.svg?branch=master)](https://travis-ci.org/benjefferies/gogo-garage-opener)
# gogo-garage-opener
Go implementation of a Raspberry Pi garage door opener

## Features

### Open garage door using an app

[Garage Opener on Play Store](https://play.google.com/store/apps/details?id=uk.echosoft.garage.opener) ([source](https://github.com/benjefferies/gogo-garage-opener-react-native))

### Use one time pin

To use a one time pin go to http://localhost:8080/user/one-time-pin/abcd1234. The pin at the end is the generated pin, once the open button has been pressed the pin will be marked as used.

### Open garage door notification
The application can be configured to notify users which have accounts via their email address if the garage door has been left open for a configurable period.
The CLI argument `-notification=15m` configures the app to notify all users if the door has been left open longer than the configuration duration.
It uses [AWS SES](https://aws.amazon.com/documentation/ses/) as an SMTP service for sending the emails.
To configure the application to use your SES account you will need to set the environmental variables $AWS_ACCESS_KEY_ID, $AWS_SECRET_KEY and $AWS_SES_ENDPOINT environmental variables. See [go-ses](https://github.com/sourcegraph/go-ses#running-tests)

### Autoclose
Auto close the door if it's left open between 10PM and 8AM.

## Guide
### Prerequisites

* Garage door opener I have an [EcoStar Liftronic 500](https://www.amazon.co.uk/gp/product/B00520C7M2/ref=oh_aui_detailpage_o03_s00?ie=UTF8&psc=1) but any model which will allow you to hook up a switch will work
* Raspberry Pi (I am using model B) wired up to internet
* [Relay switch](https://www.amazon.co.uk/gp/product/B00J4FTWO2/ref=oh_aui_detailpage_o00_s00?ie=UTF8&psc=1)
* Normal open [Magnetic switch](https://www.amazon.co.uk/gp/product/B0056K5ZC2/ref=oh_aui_detailpage_o00_s00?ie=UTF8&psc=1)
* Small [solderless breadboard](https://www.amazon.co.uk/gp/product/B0040Z4QGA/ref=oh_aui_detailpage_o09_s00?ie=UTF8&psc=1)
* Clone the repository
* [Docker](https://docs.docker.com/engine/installation/) installed
* [Docker compose](https://withblue.ink/2019/07/13/yes-you-can-run-docker-on-raspbian.html)

### To build
Look to [Dockerfile](./Dockerfile) for the latest instructions for building

#### Auth0 Setup
Authentication has been setup using Auth0. More details can be found [here](./auth0/README.md)

#### Running

The easier way to use gogo-garage-opener is by using the [docker image](https://cloud.docker.com/u/benjjefferies/repository/docker/benjjefferies/gogo-garage-opener) built on every commit

1. Clone this repo on the raspberry pi
1. Update [.env](.env)
    * AWS_ACCESS_KEY_ID is your AWS access key used for sending email notifications
    * AWS_SECRET_KEY  is your AWS secret key used for sending email notifications
    * AWS_SES_ENDPOINT is the SES endpoint used for sending email notifications (more information https://aws.amazon.com/getting-started/tutorials/send-an-email/)
    * AS is the domain of the authorisation server
    * RS is the domain of the resource server (domain set for API in) and used in Caddy (reverse proxy
    * RELAY the pin number used to toggle the relay
    * SWITCH the pin number used to read from the contact switch
1. Run the application `docker-compose up -d`

#### Hardware

I will describe how to set up your Raspberry Pi and the peripherals referring to the GPIO pins set up my Raspberry PI. Check that they are the same for your version. [GPIO documentation](https://www.raspberrypi.org/documentation/usage/gpio/)

Wiring up the relay

1. Connect up a (black) GPIO wire to pin 2 which is a 5v output
1. Connect up a (white) GPIO wire to pin 6 which is a ground
1. Connect up a (grey) GPIO wire to pin 18 which is a GPIO pin. This will be used in "out" mode to toggle the relay
1. Connect up the (black) 5v output GPIO pin 2 wire to the VCC on the relay
1. Connect up the (white) ground GPIO pin 6 wire to the GND on the relay
1. Connect up the (grey) GPIO pin 18 wire to IN1 on the relay
1. Connect the positive from the garage door to the NO1 connector on the relay
1. Connect the negative from the garage to the COM1 connector on the relay

It should now look like the images below

![Pins right view](https://i.imgur.com/jrU1R6c.jpg)
![Relay switch](https://i.imgur.com/6KsMJDC.jpg)
![Relay switch garage door opener wires](https://i.imgur.com/P8KZ5Vj.jpg)
![Garage door opener wiring](https://i.imgur.com/UYSarP8.jpg)

Wiring up the magnetic switch

1. Fix the magnet to the garage door
1. Fix the sensor to the garage door frame
1. Connect up a (brown) GPIO wire to pin 15 which is a GPIO pin. This will be used in "in" mode to read from the sensor
1. Connect up a (red) GPIO wire to pin 9 which a ground
1. Connect up the (red) ground GPIO pin 9 wire to a terminal strip on the breadboard
1. Connect up the (brown/blue) GPIO pin 15 wire to a different terminal strip on the breadboard
1. Connect one wire of the magnetic switch to the terminal strip for GPIO pin 9 and the other to the terminal strip for GPIO pin 15

![Pins left view](https://i.imgur.com/uSChY65.jpg)
![Breadboard with magnetic switch wiring](https://i.imgur.com/DVqXEzu.jpg)

#### Future

* Open garage door via location i.e. Automatically open garage within 100 metres of your garage. I have found a nice way to do this using [Automate](http://llamalab.com/automate/).
