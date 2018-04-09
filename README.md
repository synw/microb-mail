# Microb mail

Contact form service for [Microb](https://github.com/synw/microb). Features:

- **Contact form**: send emails from a form
- **Csrf protection** with one time use tokens
- **Record emails** in an sqlite database

As this is a Microb service it has remote commands and records logs into an sqlite database

Requirements: 
[Centrifugo](https://github.com/centrifugal/centrifugo/) and [Redis](http://redis.io/) and the
[Microb http service](https://github.com/synw/microb-http)

#### Install and status

To install you have to compile Microb with the email 
and the http service as for now as no release has been made. 
The dev status is move fast and break things for the moment.

## Usage

Configure `config_email.json`: ex with Postfix:

   ```javascript
   {
	   "to": "email@adress.tld",
	   "host": "localhost",
	   "port": 25,
	   "user": "",
	   "password": "",
	   "db": "mails.sqlite"
   }
   ```

The email form is avalaible at `/mail`

## Cli commands

Todo: count or read the last emails from the database

