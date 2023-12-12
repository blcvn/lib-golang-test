# smallstep CA

- sample configuration in `config.yaml` file:

```yaml
CA:
    URL: "127.0.0.1:1443"
    CertValidDuration: "8760h" # valid duration for new certificate. 8760h = 365 days
```

## Setup CA server

1. init

```bash
docker run -it -v ./step:/home/step smallstep/step-ca step ca init
```