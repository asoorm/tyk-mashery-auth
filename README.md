Mashery Signature Validator
===========================

Tyk gRPC plugin written in GoLang to handle Mashery X-Signature validation.

CLI
---

```
tyk-mashery-auth --help
Usage of tyk-mashery-auth:
  -debug
        enable debug mode
  -header_auth string
        header location to look for auth token (default "Authorization")
  -header_signature string
        header location to look for signature (default "X-Signature")
  -listen string
        listen address e.g. :9000 | /tmp/foo.sock (default ":9000")
  -network string
        network mode e.g. tcp | unix (default "tcp")
  -secret string
        shared secret (default "4321knj8fqgm5ffq64tdzifato6fb5p5rkqze933ehivqelctivti8qs0xnzmpq3")
  -skew int
        allowed clock skew in seconds (default 300)
  -token string
        token used for generating debug logs (default "foo")
```

1. Download the src & install:

```bash
go get -u github.com/asoorm/tyk-mashery-auth
```

2. Start the `tyk-mashery-auth` service:

```bash
# defaults
tyk-mashery-auth

# grpc listen on tcp://0.0.0.0:9001
tyk-mashery-auth --network tcp --listen :9001

# grpc listen on unix socket
tyk-mashery-auth --network unix --listen /tmp/foo.sock

# shared secret `mysharedsecret`
tyk-mashery-auth --secret mysharedsecret

# turn on debug mode
tyk-mashery-auth --debug

# set the allowed clock-skew to +/- 10 minutes
tyk-mashery-auth --skew 600

# override the default authorization header key & auth signature header keys
tyk-mashery-auth --header_auth Api-Auth --header_signature X-My-Signature
```

3. Modify `tyk.conf` in the gateway to point to the gRPC server and restart the service:

e.g. if server is listening on local unix socket `unix:///tmp/foo.sock`
```
"coprocess_options": {
    "enable_coprocess": true,
    "coprocess_grpc_server": "unix:///tmp/foo.sock"
},
```

e.g. if server is hosted on different machine `tcp://IP:PORT`

```
"coprocess_options": {
    "enable_coprocess": true,
    "coprocess_grpc_server": "tcp://127.0.0.1:9001"
},
```

4. Configure Tyk API to use the gRPC signature validator as a `pre` plugin:

This will ensure that the signature is validated prior to the key being checked by the authentication middleware. You could also validate the signature after auth check if you require.

Modify the api definition `custom_middleware.driver` to specify `grpc`
Modify the api definition `custom_middleware.pre[]` array to include the `ValidateSignature` hook

```json
{
  "custom_middleware": {
    "pre": [
      {
        "name": "ValidateSignature"
      }
    ],
    "driver": "grpc"
  }
}
```

5. Save the API definition and when you send API requests via the gateway, the gateway will pass responsibility for
validating the signature to the `tyk-mashery-auth` plugin.

```bash
curl http://localhost:8080/sha/get \
  -H 'Authorization: foo' \
  -H 'X-Signature: e08f7ab275ad200f041d5af0ba6bb51525905899b2bdf1825c9ea5d578ca1161'
```
