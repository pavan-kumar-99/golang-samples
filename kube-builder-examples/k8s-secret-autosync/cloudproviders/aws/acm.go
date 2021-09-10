package main

import (
	"fmt"

	aws "github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
)

var sess *acm.ACM

func init() {
	sess = acm.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})))
}

func main() {
	cert_arn := "arn:aws:acm:us-east-1:530774763960:certificate/a7f7970c-409d-4b67-a035-5f18823f0201"
	certinput := &acm.ExportCertificateInput{
		CertificateArn: &cert_arn,
		Passphrase:     []byte("Pass123"),
	}
	exported_cert, err := sess.ExportCertificate(certinput)
	if err != nil {
		fmt.Println(err.Error())
	}
	cert := exported_cert.Certificate
	cert_chain := exported_cert.CertificateChain
	pem_key := exported_cert.PrivateKey
	fmt.Println(cert, cert_chain, pem_key)
}
