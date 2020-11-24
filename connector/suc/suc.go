package suc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/dexidp/dex/connector"
	vv "github.com/go-playground/validator/v10"
	resty "github.com/go-resty/resty/v2"

	"github.com/dexidp/dex/pkg/log"
)

type Config struct {
	RootUrl string `json:"rootUrl"`

	// Required if LDAP host does not use TLS.
	InsecureNoSSL bool `json:"insecureNoSSL"`

	// Don't verify the CA.
	InsecureSkipVerify bool `json:"insecureSkipVerify"`

	// Path to a trusted root certificate file.
	RootCA string `json:"rootCA"`

	// Path to a client cert file generated by rootCA.
	ClientCert string `json:"clientCert"`
	// Path to a client private key file generated by rootCA.
	ClientKey string `json:"clientKey"`
	// Base64 encoded PEM data containing root CAs.
	RootCAData []byte `json:"rootCAData"`

	// BindDN and BindPW for an application service account. The connector uses these
	// credentials to search for users and groups.
	BindDN string `json:"bindDN"`
	BindPW string `json:"bindPW"`

	// UsernamePrompt allows users to override the username attribute (displayed
	// in the username/password prompt). If unset, the handler will use
	// "Username".
	UsernamePrompt string `json:"usernamePrompt"`

	SystemCode string `json:"systemCode"`

	InstanceCode string `json:"instanceCode"`

	APIUser string `json:"api_user"`

	APIPassword string `json:"api_password"`

	LDAPAddress []string `json:"ldap_address"`
}

type sucConnector struct {
	Config

	tlsConfig *tls.Config

	httpClient *resty.Client

	logger log.Logger
}

var (
	_ connector.PasswordConnector = (*sucConnector)(nil)
	_ connector.RefreshConnector  = (*sucConnector)(nil)
)

var validate = vv.New()

// Open returns an authentication strategy using LDAP.
func (c *Config) Open(id string, logger log.Logger) (connector.Connector, error) {
	conn, err := c.OpenConnector(logger)
	if err != nil {
		return nil, err
	}
	return connector.Connector(conn), nil
}

// OpenConnector is the same as Open but returns a type with all implemented connector interfaces.
func (c *Config) OpenConnector(logger log.Logger) (interface {
	connector.Connector
	connector.PasswordConnector
	connector.RefreshConnector
}, error) {
	return c.openConnector(logger)
}

func (c *sucConnector) validate(st interface{}) error {
	if validate != nil {
		if err := validate.Struct(st); err != nil {
			c.logger.Errorf("Validate Struct Error: %s", err)
			return err
		}
	}
	return nil
}

func (c *Config) openConnector(logger log.Logger) (*sucConnector, error) {
	requiredFields := []struct {
		name string
		val  string
	}{
		{"rootUrl", c.RootUrl},
	}

	for _, field := range requiredFields {
		if field.val == "" {
			return nil, fmt.Errorf("ldap: missing required field %q", field.name)
		}
	}

	var (
		err error
	)

	client := newHttpClient()

	tlsConfig := &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify}
	if c.RootCA != "" || len(c.RootCAData) != 0 {
		data := c.RootCAData
		if len(data) == 0 {
			if data, err = ioutil.ReadFile(c.RootCA); err != nil {
				return nil, fmt.Errorf("suc: read ca file: %v", err)
			}
			client.SetRootCertificate(c.RootCA)
		}
		rootCAs := x509.NewCertPool()
		if !rootCAs.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("suc: no certs found in ca file")
		}
		tlsConfig.RootCAs = rootCAs
	}

	if c.ClientKey != "" && c.ClientCert != "" {
		cert, err := tls.LoadX509KeyPair(c.ClientCert, c.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("suc: load client cert failed: %v", err)
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}
	client.SetTLSClientConfig(tlsConfig)
	return &sucConnector{*c, tlsConfig, client, logger}, nil
}

func (c *sucConnector) Prompt() string {
	return c.UsernamePrompt
}

func newHttpClient() *resty.Client {
	return resty.New().SetRetryCount(3).
		SetHeader("Content-Type", "application/json").
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(3 * time.Second).
		SetTimeout(15 * time.Second)
}

func (c *sucConnector) ensureHttpClient() *resty.Client {
	if c.httpClient != nil {
		return c.httpClient
	}
	client := newHttpClient()
	client.SetTLSClientConfig(c.tlsConfig)
	c.httpClient = client
	return client
}

func (c *sucConnector) doLoginAuth(ctx context.Context, username, password string) (sucLoginResp, error) {
	client := c.ensureHttpClient()

	domain := "local"

	uu := ""

	ss := strings.Split(username, "\\")
	if len(ss) == 2 {
		domain = ss[0]
		uu = ss[1]
	}else {
		uu = username
	}

	//payload := sucLoginAuthReq{Domain: "local", Account: username, Password: password}
	payload := map[string]string{
		"domain":   domain,
		"account":  uu,
		"password": password,
	}
	var r sucLoginResp
	url := fmt.Sprintf("%s%s", c.RootUrl, "/accounts/login")
	resp, err := client.R().SetContext(ctx).SetQueryParams(payload).Get(url)
	if err != nil {
		c.logger.Errorf("Post Auth Login Error: %s", err.Error())
		return r, err
	}
	body := resp.Body()
	if err := json.Unmarshal(body, &r); err != nil {
		return r, err
	}
	if !r.Success {
		return r, errors.Errorf("Auth Fail, Message: %s, Data: %s", r.Message, r.Data)
	}
	c.logger.Debugf("Do Auth Login, Resp: Status Code %d, Data: %s", resp.StatusCode(), string(body))
	return r, nil
}

func (c *sucConnector) Login(ctx context.Context, s connector.Scopes, username, password string) (connector.Identity, bool, error) {
	identity := connector.Identity{}

	if password == "" || username == "" {
		return identity, false, nil
	}

	if userInfo, err := c.doLoginAuth(ctx, username, password); err != nil {
		return identity, false, err
	} else {
		identity.UserID = strconv.Itoa(userInfo.Code)
		identity.Username = username
		identity.Email = username
		identity.EmailVerified = false
		identity.PreferredUsername = username
	}

	return identity, true, nil
}

//suc 不需要refresh，直接返回
func (c *sucConnector) Refresh(ctx context.Context, s connector.Scopes, ident connector.Identity) (connector.Identity, error) {
	return ident, nil
}
