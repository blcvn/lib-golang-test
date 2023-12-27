# smallstep CA

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

## sample configuration in `config.yaml` file:

```yaml
CA:
    URL: "127.0.0.1:1443"
    CertValidDuration: "8760h" # valid duration for new certificate. 8760h = 365 days
```
