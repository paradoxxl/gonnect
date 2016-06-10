package server

import (
	"crypto/x509"
	"math/big"
	"crypto/x509/pkix"
	"time"
	"crypto/rsa"
	"log"
	"crypto/rand"
	"fmt"
	"os"
	"encoding/pem"
	"crypto/ecdsa"
)

type certParams struct{
	Subject subject
	Issuer issuer
}
type subject struct{
	Country string
	Organization string
	OrganizationalUnit string
}
type issuer struct{
	Country string
	Organization string
	OrganizationalUnit string
	Locality string
	Province string
	StreetAddress string
	PostalCode string
	SerialNumber string
	CommonName string
}

var DefaultCert = certParams{
	Subject:subject{
		Country:"Switzerland",
		Organization:"fhnw",
		OrganizationalUnit:"imvs",
	},
	Issuer:issuer{
		Country:"Switzerland",
		Organization:"fhnw",
		OrganizationalUnit:"imvs",
		Locality:"WIndisch",
		Province:"Brugg",
		StreetAddress:"ergendwo",
		PostalCode:"kA",
		SerialNumber:"3z33833838",
		CommonName:"Gonnect",
	},
}
func GenCert(params certParams,keyFile,pemFile string){
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1337),
		Subject: pkix.Name{
			Country:            []string{params.Subject.Country},
			Organization:       []string{params.Subject.Organization},
			OrganizationalUnit: []string{params.Subject.OrganizationalUnit},
		},
		Issuer: pkix.Name{
			Country:            []string{params.Issuer.Country},
			Organization:       []string{params.Issuer.Organization},
			OrganizationalUnit: []string{params.Issuer.OrganizationalUnit},
			Locality:           []string{params.Issuer.Locality},
			Province:           []string{params.Issuer.Province},
			StreetAddress:      []string{params.Issuer.StreetAddress},
			PostalCode:         []string{params.Issuer.PostalCode},
			SerialNumber:       params.Issuer.SerialNumber,
			CommonName:         params.Issuer.CommonName,
		},
		SignatureAlgorithm:    x509.SHA512WithRSA,
		PublicKeyAlgorithm:    x509.ECDSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, 10),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}


	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("create cert failed %#v", err)
		return
	}


	pub := &priv.PublicKey
	derBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		log.Fatalf("create cert failed %#v", err)
		return
	}


	certOut, err := os.Create(pemFile)
	if err != nil {
		log.Fatalf("failed to open cert.pem for writing: %s", err)
		return
	}

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	log.Print("written cert.pem\n")


	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("failed to open key.pem for writing:", err)
		return
	}

	pem.Encode(keyOut, pemBlockForKey(priv))
	keyOut.Close()
	log.Print("written key.pem\n")


}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}