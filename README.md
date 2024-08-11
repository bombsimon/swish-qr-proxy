# Swish QR proxy

A simple proxy to convert HTTP GET requests to HTTP POST requests to the Swish
API and generate a QR code.

## Generate QR code

```sh
https://api.swish.nu/qr/v2/prefilled
```

```json
{
  "size": 240,
  "border": 1,
  "payee": "0701111111",
  "color": true,
  "message": {
    "value": "Swish via Garmin",
    "editable": true
  },
  "amount": {
    "value": 100,
    "editable": true
  }
}
```

## Deploy

```sh
NAME=gcr.io/zippy-cab-252712/swish-proxy

gcloud builds submit --tag $NAME
gcloud beta run deploy --image $NAME --platform managed
```
