
# Counter self hosted

An attempt to self host https://github.com/ihucos/counter.dev

In Progress


# Quickstart

## Install
```
$ curl https://github.com/ihucos/counter.dev-selfhosted/releases/download/0.1/cntr-linux-amd64 > /usr/local/bin/cntr
$ chmod +x /usr/local/bin/cntr
```

## Create user
```
$ cntr createuser --redis-url redis://localhost:6379 admin
Password for new user:
```

## Serve
```
$ cntr serve --redis-url redis://localhost:6379 --bind :80
```

## Go to UI
Visit the fired up server, login and follow the integrations steps there.
