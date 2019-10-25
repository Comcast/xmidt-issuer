package xhttpserver

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeerVerifyError(t *testing.T) {
	var (
		assert = assert.New(t)
		err    = PeerVerifyError{Reason: "expected"}
	)

	assert.Equal("expected", err.Error())
}

func testConfiguredPeerVerifierSuccess(t *testing.T) {
	testData := []struct {
		peerCert x509.Certificate
		options  PeerVerifyOptions
	}{
		{
			peerCert: x509.Certificate{
				DNSNames: []string{"test.foobar.com"},
			},
			options: PeerVerifyOptions{
				DNSSuffixes: []string{"foobar.com"},
			},
		},
		{
			peerCert: x509.Certificate{
				DNSNames: []string{"first.foobar.com", "second.something.net"},
			},
			options: PeerVerifyOptions{
				DNSSuffixes: []string{"another.thing.org", "something.net"},
			},
		},
		{
			peerCert: x509.Certificate{
				Subject: pkix.Name{
					CommonName: "A Great Organization",
				},
			},
			options: PeerVerifyOptions{
				CommonNames: []string{"A Great Organization"},
			},
		},
		{
			peerCert: x509.Certificate{
				Subject: pkix.Name{
					CommonName: "A Great Organization",
				},
			},
			options: PeerVerifyOptions{
				CommonNames: []string{"First Organization Doesn't Match", "A Great Organization"},
			},
		},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert  = assert.New(t)
				require = require.New(t)

				verifier = NewConfiguredPeerVerifier(record.options)
			)

			require.NotNil(verifier)
			assert.NoError(verifier.Verify(&record.peerCert, nil))
		})
	}
}

func testConfiguredPeerVerifierFailure(t *testing.T) {
	testData := []struct {
		peerCert x509.Certificate
		options  PeerVerifyOptions
	}{
		{
			peerCert: x509.Certificate{},
			options: PeerVerifyOptions{
				DNSSuffixes: []string{"foobar.net"},
				CommonNames: []string{"For Great Justice"},
			},
		},
		{
			peerCert: x509.Certificate{
				DNSNames: []string{"another.company.com"},
			},
			options: PeerVerifyOptions{
				DNSSuffixes: []string{"foobar.net"},
				CommonNames: []string{"For Great Justice"},
			},
		},
		{
			peerCert: x509.Certificate{
				Subject: pkix.Name{
					CommonName: "Villains For Hire",
				},
			},
			options: PeerVerifyOptions{
				DNSSuffixes: []string{"foobar.net"},
				CommonNames: []string{"For Great Justice"},
			},
		},
		{
			peerCert: x509.Certificate{
				DNSNames: []string{"another.company.com"},
				Subject: pkix.Name{
					CommonName: "Villains For Hire",
				},
			},
			options: PeerVerifyOptions{
				DNSSuffixes: []string{"foobar.net"},
				CommonNames: []string{"For Great Justice"},
			},
		},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert  = assert.New(t)
				require = require.New(t)

				verifier = NewConfiguredPeerVerifier(record.options)
			)

			require.NotNil(verifier)
			err := verifier.Verify(&record.peerCert, nil)
			assert.Error(err)
			require.IsType(PeerVerifyError{}, err)
			assert.Equal(&record.peerCert, err.(PeerVerifyError).Certificate)
		})
	}
}

func TestConfiguredPeerVerifier(t *testing.T) {
	t.Run("Success", testConfiguredPeerVerifierSuccess)
	t.Run("Failure", testConfiguredPeerVerifierFailure)
}

func TestNewConfiguredPeerVerifier(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert := assert.New(t)
		assert.Nil(NewConfiguredPeerVerifier(PeerVerifyOptions{}))
	})
}

func TestPeerVerifiers(t *testing.T) {
	t.Run("VerifyPeerCertificate", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			assert := assert.New(t)
			assert.NoError(PeerVerifiers{}.VerifyPeerCertificate(nil, nil))
		})
	})

	t.Run("Verify", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			assert := assert.New(t)
			assert.NoError(PeerVerifiers{}.Verify(nil, nil))
		})
	})
}

func testNewTlsConfigNil(t *testing.T) {
	assert := assert.New(t)
	tc, err := NewTlsConfig(nil)
	assert.Nil(tc)
	assert.NoError(err)
}

func testNewTlsConfigNoCertificateFile(t *testing.T) {
	assert := assert.New(t)
	tc, err := NewTlsConfig(&Tls{KeyFile: "key.pem"})
	assert.Nil(tc)
	assert.Equal(ErrTlsCertificateRequired, err)
}

func testNewTlsConfigNoKeyFile(t *testing.T) {
	assert := assert.New(t)
	tc, err := NewTlsConfig(&Tls{CertificateFile: "cert.pem"})
	assert.Nil(tc)
	assert.Equal(ErrTlsCertificateRequired, err)
}

func testNewTlsConfigLoadCertificateError(t *testing.T) {
	assert := assert.New(t)
	tc, err := NewTlsConfig(&Tls{
		CertificateFile: "nosuch",
		KeyFile:         "nosuch",
	})

	assert.Nil(tc)
	assert.Error(err)
}

func testNewTlsConfigSimple(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)

		tc, err = NewTlsConfig(&Tls{
			CertificateFile: "server.cert",
			KeyFile:         "server.key",
		})
	)

	require.NoError(err)
	require.NotNil(tc)

	assert.Zero(tc.MinVersion)
	assert.Zero(tc.MaxVersion)
	assert.Empty(tc.ServerName)
	assert.Equal([]string{"http/1.1"}, tc.NextProtos)
	assert.NotEmpty(tc.NameToCertificate)
}

func testNewTlsConfigWithoutClientCACertificateFile(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)

		tc, err = NewTlsConfig(&Tls{
			MinVersion:      1,
			MaxVersion:      3,
			ServerName:      "test",
			CertificateFile: "server.cert",
			KeyFile:         "server.key",
			NextProtos:      []string{"http/1.0"},
		})
	)

	require.NoError(err)
	require.NotNil(tc)

	assert.Equal(uint16(1), tc.MinVersion)
	assert.Equal(uint16(3), tc.MaxVersion)
	assert.Equal("test", tc.ServerName)
	assert.Equal([]string{"http/1.0"}, tc.NextProtos)
	assert.NotEmpty(tc.NameToCertificate)
}

func testNewTlsConfigWithClientCACertificateFile(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)

		tc, err = NewTlsConfig(&Tls{
			CertificateFile:         "server.cert",
			KeyFile:                 "server.key",
			ClientCACertificateFile: "server.cert",
		})
	)

	require.NoError(err)
	require.NotNil(tc)

	assert.Zero(tc.MinVersion)
	assert.Zero(tc.MaxVersion)
	assert.Empty(tc.ServerName)
	assert.Equal([]string{"http/1.1"}, tc.NextProtos)
	assert.NotEmpty(tc.NameToCertificate)
	assert.NotNil(tc.ClientCAs)
	assert.Equal(tls.RequireAndVerifyClientCert, tc.ClientAuth)
}

func testNewTlsConfigLoadClientCACertificateError(t *testing.T) {
	var (
		assert = assert.New(t)

		tc, err = NewTlsConfig(&Tls{
			CertificateFile:         "server.cert",
			KeyFile:                 "server.key",
			ClientCACertificateFile: "nosuch",
		})
	)

	assert.Nil(tc)
	assert.Error(err)
}

func testNewTlsConfigAppendClientCACertificateError(t *testing.T) {
	var (
		assert = assert.New(t)

		tc, err = NewTlsConfig(&Tls{
			CertificateFile:         "server.cert",
			KeyFile:                 "server.key",
			ClientCACertificateFile: "server.key", // not a certificate, but still valid PEM
		})
	)

	assert.Nil(tc)
	assert.Equal(ErrUnableToAddClientCACertificate, err)
}

func TestNewTlsConfig(t *testing.T) {
	t.Run("Nil", testNewTlsConfigNil)
	t.Run("NoCertificateFile", testNewTlsConfigNoCertificateFile)
	t.Run("NoKeyFile", testNewTlsConfigNoKeyFile)
	t.Run("LoadCertificateError", testNewTlsConfigLoadCertificateError)
	t.Run("Simple", testNewTlsConfigSimple)
	t.Run("WithoutClientCACertificateFile", testNewTlsConfigWithoutClientCACertificateFile)
	t.Run("WithClientCACertificateFile", testNewTlsConfigWithClientCACertificateFile)
	t.Run("LoadClientCACertificateError", testNewTlsConfigLoadClientCACertificateError)
	t.Run("AppendClientCACertificateError", testNewTlsConfigAppendClientCACertificateError)
}
