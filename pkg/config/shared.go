// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package config

import (
	"encoding/hex"

	"github.com/TheThingsNetwork/ttn/pkg/errors"
	"github.com/TheThingsNetwork/ttn/pkg/log"
)

// Base represents base component configuration.
type Base struct {
	Config []string `name:"config" shorthand:"c" description:"Location of the config files"`
	Log    Log      `name:"log"`
	Silent bool     `name:"silent" shorthand:"s" description:"Suppress all log messages"`
}

// Log represents configuration for the logger.
type Log struct {
	Level log.Level `name:"level" description:"The minimum level log messages must have to be shown"`
}

// Sentry represents configuration for error tracking using Sentry.
type Sentry struct {
	DSN string `name:"dsn" description:"Sentry Data Source Name"`
}

// TLS represents TLS configuration.
type TLS struct {
	Certificate string `name:"certificate" description:"Location of TLS certificate"`
	Key         string `name:"key" description:"Location of TLS private key"`
}

// Cluster represents clustering configuration.
type Cluster struct {
	Join              []string `name:"join" description:"Addresses of cluster peers to join"`
	Name              string   `name:"name" description:"Name of the current cluster peer (default: $HOSTNAME)"`
	Address           string   `name:"address" description:"Address to use for cluster communication"`
	IdentityServer    string   `name:"identity-server" description:"Address for the identity server"`
	GatewayServer     string   `name:"gateway-server" description:"Address for the gateway server"`
	NetworkServer     string   `name:"network-server" description:"Address for the network server"`
	ApplicationServer string   `name:"application-server" description:"Address for the application server"`
	JoinServer        string   `name:"join-server" description:"Address for the join server"`
	TLS               bool     `name:"tls" description:"Do cluster gRPC over TLS"`
}

// GRPC represents gRPC listener configuration.
type GRPC struct {
	Listen    string `name:"listen" description:"Address for the TCP gRPC server to listen on"`
	ListenTLS string `name:"listen-tls" description:"Address for the TLS gRPC server to listen on"`
}

// Cookie represents cookie configuration.
type Cookie struct {
	HashKey  string `name:"hash-key" description:"Key for cookie contents verification"`
	BlockKey string `name:"block-key" description:"Key for cookie contents encryption"`
}

// Keys parses the hash and block keys from the Cookie config from Hex to bytes.
// It returns an error if the keys do not have a valid length.
func (c Cookie) Keys() (hash []byte, block []byte, err error) {
	if c.HashKey != "" {
		hash, err = hex.DecodeString(c.HashKey)
		if err != nil {
			return nil, nil, err
		}
		switch len(hash) {
		case 16, 24, 32:
		default:
			return nil, nil, errors.New("Invalid length for cookie hash key: must be 16, 24 or 32 bytes")
		}
	}
	if c.BlockKey != "" {
		block, err = hex.DecodeString(c.BlockKey)
		if err != nil {
			return nil, nil, err
		}
		switch len(block) {
		case 16, 24, 32:
		default:
			return nil, nil, errors.New("Invalid length for cookie block key: must be 16, 24 or 32 bytes")
		}
	}
	return
}

// HTTP represents the HTTP and HTTPS server configuration.
type HTTP struct {
	Listen    string `name:"listen" description:"Address for the HTTP server to listen on"`
	ListenTLS string `name:"listen-tls" description:"Address for the HTTPS server to listen on"`
	Cookie    Cookie `name:"cookie"`
	PProf     bool   `name:"pprof" description:"Expose pprof over HTTP"`
}

// Identity represents identity configuration.
type Identity struct {
	Servers map[string]string `name:"servers" description:"TTN Identity Servers (id=https://...)"`
	Keys    map[string]string `name:"keys" description:"TTN Identity Server Public Keys (id=/path/to/...)"`
}

// Redis represents Redis configuration.
type Redis struct {
	Address  string `name:"address" description:"Address of the Redis server"`
	Database int    `name:"database" description:"Redis database to use"`
	Prefix   string `name:"prefix" description:"Prefix for Redis keys"`
}

// RemoteProviderConfig represents remote config provider configuration(see Viper documentation).
type RemoteProviderConfig struct {
	Name     string `name:"name" description:"Name of the config on the remote without the extension"`
	Provider string `name:"provider" description:"Remote config provider name"`
	Endpoint string `name:"endpoint" description:"Endpoint where the remote config provider is accessible"`
	Path     string `name:"path" description:"Path where to look for the config on the remote"`
	KeyRing  string `name:"keyring" description:"Optional path to secret keyring for initializing a secure connection to remote config provider"`
}

// ServiceBase represents base service configuration.
type ServiceBase struct {
	Base         `name:",squash"`
	Cluster      Cluster               `name:"cluster"`
	Redis        Redis                 `name:"redis"`
	GRPC         GRPC                  `name:"grpc"`
	HTTP         HTTP                  `name:"http"`
	TLS          TLS                   `name:"tls"`
	Identity     Identity              `name:"identity"`
	RemoteConfig *RemoteProviderConfig `name:"remote-config"`
	Sentry       Sentry                `name:"sentry"`
}

// IsValid returns wether or not the remote config is valid or not.
func (c RemoteProviderConfig) IsValid() bool {
	return c.Provider != "" && c.Endpoint != "" && c.Path != ""
}
