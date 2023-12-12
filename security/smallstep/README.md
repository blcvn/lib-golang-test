# ecdsa-ca-poc

## CA Server
https://github.com/smallstep/certificates

1. init CA Server
```shell
step ca init
```

2. Run your certificate authority
```shell
step-ca $(step path)/config/ca.json
```

3. Test POC
```shell
make start
```