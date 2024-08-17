# Swish QR proxy

A proxy to convert HTTP GET requests to HTTP POST requests to the Swish API and
generate a QR code.

> [!IMPORTANT]
> This app is backed by a public API hosted by [Swish]. This API can at any time
> be removed, changed or blocked for whatever reason which will result in this
> app not working untul an alternative solution is fixed.

## Deploy

```sh
NAME=gcr.io/zippy-cab-252712/swish-proxy

gcloud builds submit --tag $NAME
gcloud beta run deploy --image $NAME --platform managed
```

## Disclaimer

This app is not affiliated with [Swish]. Swish and its logo are trademarks or
registered trademarks of Swish. All rights to the name, logo, and API belong to
Swish.

[Swish]: https://www.swish.nu
