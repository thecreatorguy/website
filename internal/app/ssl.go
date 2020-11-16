package app

import (
	"crypto/tls"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

type keypairReloader struct {
	certMu   sync.RWMutex
	cert     *tls.Certificate
	certPath string
	keyPath  string
}

func NewKeypairReloader(certPath, keyPath string) (*keypairReloader, error) { 
	result := &keypairReloader{
		certPath: certPath,
		keyPath:  keyPath,
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	result.cert = &cert
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP)
		for range c {
			logrus.Printf("Received SIGHUP, reloading TLS certificate and key from %v and %v", certPath, keyPath)
			if err := result.reload(); err != nil {
				logrus.WithError(err).Warn("Keeping old TLS certificate because the new one could not be loaded")
			}
		}
	}()
	return result, nil
}

func (kpr *keypairReloader) reload() error { 
	newCert, err := tls.LoadX509KeyPair(kpr.certPath, kpr.keyPath)
	if err != nil {
			return err
	}
	kpr.certMu.Lock()
	defer kpr.certMu.Unlock()
	kpr.cert = &newCert
	return nil
}

func (kpr *keypairReloader) GetCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) { 
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			kpr.certMu.RLock()
			defer kpr.certMu.RUnlock()
			return kpr.cert, nil
	}
}