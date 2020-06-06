# lazy-rest

Lazy REST Server in Go.

PUT a JSON on any endpoint, and it will be returned if you GET that endpoint.

POST to an endpoint and it expects an id field in your JSON.
You can then retrieve that JSON at endpoint/id.

See the [tests](./test/requests.http) for some examples.

## Docker

The image is available on docker hub [here](https://hub.docker.com/r/akleinloog/lazy-rest)


## Version tagging

```
git tag -a v0.1.0 -m "Version 0.1.0"

git push origin --tags
```
