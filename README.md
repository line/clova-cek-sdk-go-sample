# clova-cek-sdk-go-sample

## How to Setup

Clone this repository.
```bash
git clone https://github.com/line/clova-cek-sdk-go-sample.git
```

Dependency
```bash
go get github.com/line/clova-cek-sdk-go/cek
go get github.com/line/line-bot-sdk-go/linebot
```

Set ``EXTENSION_ID``
```bash
export EXTENSION_ID=<YOUR_EXTENSION_ID>
```

If you test local environment, set ``DEBUG_MODE=true`` to avoid request validation.
```bash
export DEBUG_MODE=true
```

## Deploy to Heroku
Install Heroku CLI: [The Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli)

Create Heroku app.
```bash
heroku create
```

Setup govendor
```bash
govendor init
govendor fetch github.com/line/clova-cek-sdk-go/cek
govendor fetch github.com/line/line-bot-sdk-go/linebot
```

Set environment variables
```bash
heroku config:set EXTENSION_ID=<YOUR_EXTENSION_ID>
```

If you use Messaging API
```bash
heroku config:set CHANNEL_SECRET=<YOUR_CHANNEL_SECRET>
heroku config:set CHANNEL_ACCESS_TOKEN=<YOUR_CHANNEL_ACCESS_TOKEN>
```

Deploy
```bash
git push heroku master
```

