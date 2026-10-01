package airflow

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Auth struct {
	CaCert   string
	Username string
	Password string
}

func (auth *Auth) getTlsConfigs() (*tls.Config, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: false,
	}

	if auth.CaCert != "" {
		caCertContent, err := ioutil.ReadFile(auth.CaCert)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to read patroni CA certificate file: %s", err.Error()))
		}
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(caCertContent)
		if !ok {
			return nil, errors.New("Failed to parse patroni CA certificat")
		}
		(*tlsConf).RootCAs = roots
	}

	return tlsConf, nil
}

type ClientConfig struct {
	ServerAddress     string
	Auth              Auth
	ConnectionTimeout time.Duration
	RequestTimeout    time.Duration
}

type Client struct {
	conf        *ClientConfig
	handle      *http.Client
	accessToken string
	logger      Logger
}

func NewClient(conf *ClientConfig) (*Client, error) {
	tlsConf, tlsConfErr := conf.Auth.getTlsConfigs()
	if tlsConfErr != nil {
		return nil, tlsConfErr
	}

	return &Client{
		handle: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig:       tlsConf,
				TLSHandshakeTimeout:   conf.ConnectionTimeout,
				IdleConnTimeout:       conf.RequestTimeout,
				ResponseHeaderTimeout: conf.RequestTimeout,
			},
		},
		conf: conf,
	}, nil
}

func (cli *Client) SetLogger(logger Logger) {
	cli.logger = logger
}

func (cli *Client) LogRequest(ctx context.Context, method string, url string, body string) {
	if cli.logger != nil {
		cli.logger.Info(ctx, fmt.Sprintf("REQUEST -> %s %s", method, url))
		if body != "" {
			cli.logger.Debug(ctx, fmt.Sprintf("\n*****REQUEST BODY*****\n%s\n**********************\n", body))
		}
	}
}

func (cli *Client) BuildUrl(path string) (string, error) {
	baseURL, err := url.Parse(cli.conf.ServerAddress)
	if err != nil {
		return "", fmt.Errorf("Invalid server address: %w", err)
	}
	
	return baseURL.ResolveReference(&url.URL{Path: path}).String(), nil
}