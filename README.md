
# Self-Hosted [Counter Web Analytics](https://github.com/ihucos/counter.dev)

This is the official way to self host [counter.dev](https://counter.dev/).

Please note:
- The self-hosted version is maturing and can be currenlty considered **beta**
- Due to complexity constraints in this iteration the archiving functionality is not implemented.
- Please create an issue if there is anything

# Quickstart

## 1. Install
```
$ curl https://github.com/ihucos/counter.dev-selfhosted/releases/download/0.2/cntr-linux-amd64 > /usr/local/bin/cntr
$ chmod +x /usr/local/bin/cntr
```

## 2. Create user
```
$ cntr createuser --redis-url redis://localhost:6379 --utc-offset 2 admin
Password for new user:
```
("Utc offset" is your timezones utc offset.)


## 3. Serve
```
$ cntr serve --redis-url redis://localhost:6379 --bind :80
```

## 4. Go to UI
Visit the fired up server, login and follow the integrations steps for your website there.

# Screenshots

<img width="1440" alt="Screenshot 2023-03-26 at 21 21 36" src="https://user-images.githubusercontent.com/2066372/227825413-307290db-2d38-4443-adbb-e22df6304c73.png">

<img width="1440" alt="Screenshot 2023-03-26 at 21 24 30" src="https://user-images.githubusercontent.com/2066372/227825733-118fb7c8-c1af-4af0-8bc9-7f38b0af53c0.png">

# Technical Description

Redis is the primary database. With Golang all static assets are nicely inside a single binary which you just need to run at your server. Everything basically works exactly as in the [hosted version of this project](https://github.com/ihucos/counter.dev), consult it's README.


# I forgot my password.

Keep calm and don't email me. Run this at your server:
```
$ cntr chgpwd --redis-url redis://localhost:6379 youruser
```
