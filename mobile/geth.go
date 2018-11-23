// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php


// Contains all the wrappers from the node package to support client side node
// management on mobile platforms.

package gehc

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ecosystem/go-ecosystem/core"
	"github.com/ecosystem/go-ecosystem/ehc"
	"github.com/ecosystem/go-ecosystem/ehc/downloader"
	"github.com/ecosystem/go-ecosystem/ehcclient"
	"github.com/ecosystem/go-ecosystem/ehcstats"
	"github.com/ecosystem/go-ecosystem/internal/debug"
	"github.com/ecosystem/go-ecosystem/les"
	"github.com/ecosystem/go-ecosystem/node"
	"github.com/ecosystem/go-ecosystem/p2p"
	"github.com/ecosystem/go-ecosystem/p2p/nat"
	"github.com/ecosystem/go-ecosystem/params"
	whisper "github.com/ecosystem/go-ecosystem/whisper/whisperv6"
)

// NodeConfig represents the collection of configuration values to fine tune the Gehc
// node embedded into a mobile process. The available values are a subset of the
// entire API provided by go-ecosystem to reduce the maintenance surface and dev
// complexity.
type NodeConfig struct {
	// Bootstrap nodes used to establish connectivity with the rest of the network.
	BootstrapNodes *Enodes

	// MaxPeers is the maximum number of peers that can be connected. If this is
	// set to zero, then only the configured static and trusted peers can connect.
	MaxPeers int

	// EcosystemEnabled specifies whether the node should run the Ecosystem protocol.
	EcosystemEnabled bool

	// EcosystemNetworkID is the network identifier used by the Ecosystem protocol to
	// decide if remote peers should be accepted or not.
	EcosystemNetworkID int64 // uint64 in truth, but Java can't handle that...

	// EcosystemGenesis is the genesis JSON to use to seed the blockchain with. An
	// empty genesis state is equivalent to using the mainnet's state.
	EcosystemGenesis string

	// EcosystemDatabaseCache is the system memory in MB to allocate for database caching.
	// A minimum of 16MB is always reserved.
	EcosystemDatabaseCache int

	// EcosystemNetStats is a netstats connection string to use to report various
	// chain, transaction and node stats to a monitoring server.
	//
	// It has the form "nodename:secret@host:port"
	EcosystemNetStats string

	// WhisperEnabled specifies whether the node should run the Whisper protocol.
	WhisperEnabled bool

	// Listening address of pprof server.
	PprofAddress string
}

// defaultNodeConfig contains the default node configuration values to use if all
// or some fields are missing from the user's specified list.
var defaultNodeConfig = &NodeConfig{
	BootstrapNodes:        FoundationBootnodes(),
	MaxPeers:              25,
	EcosystemEnabled:       true,
	EcosystemNetworkID:     1,
	EcosystemDatabaseCache: 16,
}

// NewNodeConfig creates a new node option set, initialized to the default values.
func NewNodeConfig() *NodeConfig {
	config := *defaultNodeConfig
	return &config
}

// Node represents a Gehc Ecosystem node instance.
type Node struct {
	node *node.Node
}

// NewNode creates and configures a new Gehc node.
func NewNode(datadir string, config *NodeConfig) (stack *Node, _ error) {
	// If no or partial configurations were specified, use defaults
	if config == nil {
		config = NewNodeConfig()
	}
	if config.MaxPeers == 0 {
		config.MaxPeers = defaultNodeConfig.MaxPeers
	}
	if config.BootstrapNodes == nil || config.BootstrapNodes.Size() == 0 {
		config.BootstrapNodes = defaultNodeConfig.BootstrapNodes
	}

	if config.PprofAddress != "" {
		debug.StartPProf(config.PprofAddress)
	}

	// Create the empty networking stack
	nodeConf := &node.Config{
		Name:        clientIdentifier,
		Version:     params.Version,
		DataDir:     datadir,
		KeyStoreDir: filepath.Join(datadir, "keystore"), // Mobile should never use internal keystores!
		P2P: p2p.Config{
			NoDiscovery:      true,
			DiscoveryV5:      true,
			BootstrapNodesV5: config.BootstrapNodes.nodes,
			ListenAddr:       ":0",
			NAT:              nat.Any(),
			MaxPeers:         config.MaxPeers,
		},
	}
	rawStack, err := node.New(nodeConf)
	if err != nil {
		return nil, err
	}

	debug.Memsize.Add("node", rawStack)

	var genesis *core.Genesis
	if config.EcosystemGenesis != "" {
		// Parse the user supplied genesis spec if not mainnet
		genesis = new(core.Genesis)
		if err := json.Unmarshal([]byte(config.EcosystemGenesis), genesis); err != nil {
			return nil, fmt.Errorf("invalid genesis spec: %v", err)
		}
		// If we have the testnet, hard code the chain configs too
		if config.EcosystemGenesis == TestnetGenesis() {
			genesis.Config = params.TestnetChainConfig
			if config.EcosystemNetworkID == 1 {
				config.EcosystemNetworkID = 3
			}
		}
	}
	// Register the Ecosystem protocol if requested
	if config.EcosystemEnabled {
		ehcConf := ehc.DefaultConfig
		ehcConf.Genesis = genesis
		ehcConf.SyncMode = downloader.LightSync
		ehcConf.NetworkId = uint64(config.EcosystemNetworkID)
		ehcConf.DatabaseCache = config.EcosystemDatabaseCache
		if err := rawStack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
			return les.New(ctx, &ehcConf)
		}); err != nil {
			return nil, fmt.Errorf("ecosystem init: %v", err)
		}
		// If netstats reporting is requested, do it
		if config.EcosystemNetStats != "" {
			if err := rawStack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
				var lesServ *les.LightEcosystem
				ctx.Service(&lesServ)

				return ehcstats.New(config.EcosystemNetStats, nil, lesServ)
			}); err != nil {
				return nil, fmt.Errorf("netstats init: %v", err)
			}
		}
	}
	// Register the Whisper protocol if requested
	if config.WhisperEnabled {
		if err := rawStack.Register(func(*node.ServiceContext) (node.Service, error) {
			return whisper.New(&whisper.DefaultConfig), nil
		}); err != nil {
			return nil, fmt.Errorf("whisper init: %v", err)
		}
	}
	return &Node{rawStack}, nil
}

// Start creates a live P2P node and starts running it.
func (n *Node) Start() error {
	return n.node.Start()
}

// Stop terminates a running node along with all it's services. In the node was
// not started, an error is returned.
func (n *Node) Stop() error {
	return n.node.Stop()
}

// GetEcosystemClient retrieves a client to access the Ecosystem subsystem.
func (n *Node) GetEcosystemClient() (client *EcosystemClient, _ error) {
	rpc, err := n.node.Attach()
	if err != nil {
		return nil, err
	}
	return &EcosystemClient{ehcclient.NewClient(rpc)}, nil
}

// GetNodeInfo gathers and returns a collection of metadata known about the host.
func (n *Node) GetNodeInfo() *NodeInfo {
	return &NodeInfo{n.node.Server().NodeInfo()}
}

// GetPeersInfo returns an array of metadata objects describing connected peers.
func (n *Node) GetPeersInfo() *PeerInfos {
	return &PeerInfos{n.node.Server().PeersInfo()}
}
