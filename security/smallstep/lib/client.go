package smallstep

import (
	"flag"
	"path/filepath"
	"strings"

	"github.com/Hnampk/prometheuslog/flogging"
	"github.com/pkg/errors"
	"github.com/smallstep/cli/flags"
	"github.com/smallstep/cli/utils/cautils"
	"github.com/urfave/cli"
	"go.step.sm/cli-utils/errs"
	"go.step.sm/cli-utils/step"
	"go.step.sm/cli-utils/token"
	"go.step.sm/cli-utils/ui"
	"go.step.sm/crypto/pemutil"
)

var (
	clientLogger = flogging.MustGetLogger("libs.ca.smallstep.client")
)

func GenerateCertificate(subject string) error {
	crtFile, keyFile := subject+".crt", subject+".key"
	tok := ""
	sans := []string{}

	flagSet := &flag.FlagSet{}
	flagSet.String(flags.CaURL.Name, CAURL, "CAURL")
	flagSet.String(flags.NotAfter.Name, CertValidDuration, "NotAfter")

	// if err := flagSet.Set("ca-url", "127.0.0.1:1443"); err != nil {
	// 	return err
	// }
	ctx := cli.NewContext(&cli.App{
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name: "san",
				Usage: `Add <dns|ip|email|uri> Subject Alternative Name(s) (SANs)
	that should be authorized. Use the '--san' flag multiple times to configure
	multiple SANs. The '--san' flag and the '--token' flag are mutually exclusive.`,
			},
			cli.StringFlag{
				Name:  "attestation-ca-url",
				Usage: "The base url of the Attestation CA to use",
			},
			cli.StringFlag{
				Name:  "attestation-ca-root",
				Usage: "The path to the PEM <file> with trusted roots when connecting to the Attestation CA",
			},
			cli.BoolFlag{
				Name:   "attestation-ca-insecure",
				Usage:  "Disables TLS server validation when connecting to the Attestation CA",
				Hidden: true,
			},
			cli.StringFlag{
				Name:  "tpm-storage-directory",
				Usage: "The directory where TPM keys and certificates will be stored",
				Value: filepath.Join(step.Path(), "tpm"),
			},
			flags.TemplateSet,
			flags.TemplateSetFile,
			flags.CaConfig,
			flags.CaURL,
			flags.Root,
			flags.Token,
			flags.Context,
			flags.Provisioner,
			flags.ProvisionerPasswordFile,
			flags.KTY,
			flags.Curve,
			flags.Size,
			flags.NotAfter,
			flags.NotBefore,
			flags.AttestationURI,
			flags.Force,
			flags.Offline,
			flags.PasswordFile,
			flags.KMSUri,
			flags.X5cCert,
			flags.X5cKey,
			flags.X5cChain,
			flags.NebulaCert,
			flags.NebulaKey,
			flags.K8sSATokenPathFlag,
		},
	}, flagSet, nil)

	// certificate flow unifies online and offline flows on a single api
	flow, err := cautils.NewCertificateFlow(ctx)
	if err != nil {
		return err
	}

	if tok == "" {
		// Use the ACME protocol with a different certificate authority.
		if ctx.IsSet("acme") {
			return cautils.ACMECreateCertFlow(ctx, "")
		}
		if tok, err = flow.GenerateToken(ctx, subject, sans); err != nil {
			var acmeTokenErr *cautils.ACMETokenError
			if errors.As(err, &acmeTokenErr) {
				return cautils.ACMECreateCertFlow(ctx, acmeTokenErr.Name)
			}
			return err
		}
	}

	req, pk, err := flow.CreateSignRequest(ctx, tok, subject, sans)
	if err != nil {
		return err
	}

	jwt, err := token.ParseInsecure(tok)
	if err != nil {
		return err
	}

	switch jwt.Payload.Type() {
	case token.JWK: // Validate that subject matches the CSR common name.
		if ctx.String("token") != "" && len(sans) > 0 {
			return errs.MutuallyExclusiveFlags(ctx, "token", "san")
		}
		if !strings.EqualFold(subject, req.CsrPEM.Subject.CommonName) {
			return errors.Errorf("token subject '%s' and argument '%s' do not match", req.CsrPEM.Subject.CommonName, subject)
		}
	case token.OIDC, token.AWS, token.GCP, token.Azure, token.K8sSA:
		// Common name will be validated on the server side, it depends on
		// server configuration.
	default:
		return errors.New("token is not supported")
	}

	if err := flow.Sign(ctx, tok, req.CsrPEM, crtFile); err != nil {
		return err
	}

	_, err = pemutil.Serialize(pk, pemutil.ToFile(keyFile, 0600))
	if err != nil {
		return err
	}

	ui.PrintSelected("Certificate", crtFile)
	ui.PrintSelected("Private Key", keyFile)
	return nil
}
