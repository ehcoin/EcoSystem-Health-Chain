// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php


// Package ehc implements the Ecosystem protocol.
package ehc

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/ecosystem/go-ecosystem/random"

	"github.com/ecosystem/go-ecosystem/ca"

	"github.com/ecosystem/go-ecosystem/mc"
	"github.com/ecosystem/go-ecosystem/reelection"

	"github.com/ecosystem/go-ecosystem/accounts"
	"github.com/ecosystem/go-ecosystem/accounts/signhelper"
	"github.com/ecosystem/go-ecosystem/blkconsensus/blkverify"
	"github.com/ecosystem/go-ecosystem/blockgenor"
	"github.com/ecosystem/go-ecosystem/broadcastTx"
	"github.com/ecosystem/go-ecosystem/common"
	"github.com/ecosystem/go-ecosystem/common/hexutil"
	"github.com/ecosystem/go-ecosystem/consensus"
	"github.com/ecosystem/go-ecosystem/consensus/clique"
	"github.com/ecosystem/go-ecosystem/consensus/ehcash"
	"github.com/ecosystem/go-ecosystem/core"
	"github.com/ecosystem/go-ecosystem/core/bloombits"
	"github.com/ecosystem/go-ecosystem/core/rawdb"
	"github.com/ecosystem/go-ecosystem/core/types"
	"github.com/ecosystem/go-ecosystem/core/vm"
	"github.com/ecosystem/go-ecosystem/depoistInfo"
	"github.com/ecosystem/go-ecosystem/ehc/downloader"
	"github.com/ecosystem/go-ecosystem/ehc/filters"
	"github.com/ecosystem/go-ecosystem/ehc/gasprice"
	"github.com/ecosystem/go-ecosystem/ehcdb"
	"github.com/ecosystem/go-ecosystem/event"
	"github.com/ecosystem/go-ecosystem/hd"
	"github.com/ecosystem/go-ecosystem/internal/ehcapi"
	"github.com/ecosystem/go-ecosystem/log"
	"github.com/ecosystem/go-ecosystem/miner"
	"github.com/ecosystem/go-ecosystem/node"
	"github.com/ecosystem/go-ecosystem/p2p"
	"github.com/ecosystem/go-ecosystem/params"
	"github.com/ecosystem/go-ecosystem/rlp"
	"github.com/ecosystem/go-ecosystem/rpc"
	"github.com/ecosystem/go-ecosystem/topnode"
	"github.com/ecosystem/go-ecosystem/verifier"

	"sync"
)

var MsgCenter *mc.Center

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *core.ChainIndexer)
}

// Ecosystem implements the Ecosystem full node service.
type Ecosystem struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan chan bool // Channel for shutting down the Ecosystem

	// Handlers
	txPool          *core.TxPool
	blockchain      *core.BlockChain
	protocolManager *ProtocolManager
	lesServer       LesServer

	// DB interfaces
	chainDb ehcdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	APIBackend *EthAPIBackend

	miner     *miner.Miner
	gasPrice  *big.Int
	ehcbase common.Address

	networkId     uint64
	netRPCService *ehcapi.PublicNetAPI

	broadTx *broadcastTx.BroadCast //YY

	//algorithm
	ca         *ca.Identity //node传进来的
	msgcenter  *mc.Center   //node传进来的
	hd         *hd.HD       //node传进来的
	signHelper *signhelper.SignHelper

	reelection   *reelection.ReElection //换届服务
	random       *random.Random
	topNode      *topnode.TopNodeService
	blockgen     *blockgenor.BlockGenor
	blockVerify  *blkverify.BlockVerify
	leaderServer *verifier.LeaderIdentity

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and ehcbase)
}

func (s *Ecosystem) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

// New creates a new Ecosystem object (including the
// initialisation of the common Ecosystem object)
func New(ctx *node.ServiceContext, config *Config) (*Ecosystem, error) {
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run ehc.Ecosystem in light sync mode, use les.LightEcosystem")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	ehc := &Ecosystem{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		ca:             ctx.Ca,
		msgcenter:      ctx.MsgCenter,
		hd:             ctx.HD,
		signHelper:     ctx.SignHelper,

		engine:        CreateConsensusEngine(ctx, &config.Ethash, chainConfig, chainDb),
		shutdownChan:  make(chan bool),
		networkId:     config.NetworkId,
		gasPrice:      config.GasPrice,
		ehcbase:     config.Etherbase,
		bloomRequests: make(chan chan *bloombits.Retrieval),
		bloomIndexer:  NewBloomIndexer(chainDb, params.BloomBitsBlocks),
	}
	log.Info("Initialising Ecosystem protocol", "versions", ProtocolVersions, "network", config.NetworkId)

	if !config.SkipBcVersionCheck {
		bcVersion := rawdb.ReadDatabaseVersion(chainDb)
		if bcVersion != core.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run gehc upgradedb.\n", bcVersion, core.BlockChainVersion)
		}
		rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
	}
	var (
		vmConfig    = vm.Config{EnablePreimageRecording: config.EnablePreimageRecording}
		cacheConfig = &core.CacheConfig{Disabled: config.NoPruning, TrieNodeLimit: config.TrieCache, TrieTimeLimit: config.TrieTimeout}
	)
	ehc.blockchain, err = core.NewBlockChain(chainDb, cacheConfig, ehc.chainConfig, ehc.engine, vmConfig)
	if err != nil {
		return nil, err
	}
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		ehc.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	ehc.bloomIndexer.Start(ehc.blockchain)

	ca.SetTopologyReader(ehc.blockchain.TopologyStore())

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	ehc.txPool = core.NewTxPool(config.TxPool, ehc.chainConfig, ehc.blockchain, ctx.GetConfig().DataDir)

	if ehc.protocolManager, err = NewProtocolManager(ehc.chainConfig, config.SyncMode, config.NetworkId, ehc.eventMux, ehc.txPool, ehc.engine, ehc.blockchain, chainDb, ctx.MsgCenter); err != nil {
		return nil, err
	}
	//ehc.protocolManager.Msgcenter = ctx.MsgCenter
	MsgCenter = ctx.MsgCenter
	ehc.miner, err = miner.New(ehc.blockchain, ehc.chainConfig, ehc.EventMux(), ehc.engine, ehc.blockchain.DPOSEngine(), ehc.hd)
	if err != nil {
		return nil, err
	}
	ehc.miner.SetExtra(makeExtraData(config.ExtraData))

	//algorithm
	dbDir := ctx.GetConfig().DataDir
	ehc.reelection, err = reelection.New(ehc.blockchain, dbDir)
	if err != nil {
		return nil, err
	}
	ehc.random, err = random.New(ehc.msgcenter)
	if err != nil {
		return nil, err
	}

	ehc.APIBackend = &EthAPIBackend{ehc, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	ehc.APIBackend.gpo = gasprice.NewOracle(ehc.APIBackend, gpoParams)
	depoistInfo.NewDepositInfo(ehc.APIBackend)
	ehc.broadTx = broadcastTx.NewBroadCast(ehc.APIBackend) //YY

	ehc.leaderServer, err = verifier.NewLeaderIdentityService(ehc, "leader服务")

	ehc.topNode = topnode.NewTopNodeService(ehc.blockchain.DPOSEngine())
	topNodeInstance := topnode.NewTopNodeInstance(ehc.signHelper, ehc.hd)
	ehc.topNode.SetValidatorReader(ehc.blockchain)
	ehc.topNode.SetTopNodeStateInterface(topNodeInstance)
	ehc.topNode.SetValidatorAccountInterface(topNodeInstance)
	ehc.topNode.SetMessageSendInterface(topNodeInstance)
	ehc.topNode.SetMessageCenterInterface(topNodeInstance)

	if err = ehc.topNode.Start(); err != nil {
		return nil, err
	}

	ehc.blockgen, err = blockgenor.New(ehc)
	if err != nil {
		return nil, err
	}

	ehc.blockVerify, err = blkverify.NewBlockVerify(ehc)
	if err != nil {
		return nil, err
	}

	return ehc, nil
}

func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"gehc",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}
	return extra
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (ehcdb.Database, error) {
	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*ehcdb.LDBDatabase); ok {
		db.Meter("ehc/db/chaindata/")
	}
	return db, nil
}

// CreateConsensusEngine creates the required type of consensus engine instance for an Ecosystem service
func CreateConsensusEngine(ctx *node.ServiceContext, config *ehcash.Config, chainConfig *params.ChainConfig, db ehcdb.Database) consensus.Engine {
	// If proof-of-authority is requested, set it up
	if chainConfig.Clique != nil {
		return clique.New(chainConfig.Clique, db)
	}
	// Otherwise assume proof-of-work
	switch config.PowMode {
	case ehcash.ModeFake:
		log.Warn("Ethash used in fake mode")
		return ehcash.NewFaker()
	case ehcash.ModeTest:
		log.Warn("Ethash used in test mode")
		return ehcash.NewTester()
	case ehcash.ModeShared:
		log.Warn("Ethash used in shared mode")
		return ehcash.NewShared()
	default:
		engine := ehcash.New(ehcash.Config{
			CacheDir:       ctx.ResolvePath(config.CacheDir),
			CachesInMem:    config.CachesInMem,
			CachesOnDisk:   config.CachesOnDisk,
			DatasetDir:     config.DatasetDir,
			DatasetsInMem:  config.DatasetsInMem,
			DatasetsOnDisk: config.DatasetsOnDisk,
		})
		engine.SetThreads(-1) // Disable CPU mining
		return engine
	}
}

// APIs return the collection of RPC services the ecosystem package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *Ecosystem) APIs() []rpc.API {
	apis := ehcapi.GetAPIs(s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicEcosystemAPI(s),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.APIBackend, false),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *Ecosystem) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Ecosystem) Etherbase() (eb common.Address, err error) {
	s.lock.RLock()
	ehcbase := s.ehcbase
	s.lock.RUnlock()

	if ehcbase != (common.Address{}) {
		return ehcbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			ehcbase := accounts[0].Address

			s.lock.Lock()
			s.ehcbase = ehcbase
			s.lock.Unlock()

			log.Info("Etherbase automatically configured", "address", ehcbase)
			return ehcbase, nil
		}
	}
	return common.Address{}, fmt.Errorf("ehcbase must be explicitly specified")
}

func (s *Ecosystem) StartMining(local bool) error {
	eb, err := s.Etherbase()
	if err != nil {
		log.Error("Cannot start mining without ehcbase", "err", err)
		return fmt.Errorf("ehcbase missing: %v", err)
	}
	if clique, ok := s.engine.(*clique.Clique); ok {
		wallet, err := s.accountManager.Find(accounts.Account{Address: eb})
		if wallet == nil || err != nil {
			log.Error("Etherbase account unavailable locally", "err", err)
			return fmt.Errorf("signer missing: %v", err)
		}
		clique.Authorize(eb, wallet.SignHash)
	}
	if local {
		// If local (CPU) mining is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU mining on mainnet is ludicrous
		// so none will ever hit this path, whereas marking sync done on CPU mining
		// will ensure that private networks work in single miner mode too.
		atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)
	}
	go s.miner.Start()
	return nil
}

func (s *Ecosystem) StopMining()         { s.miner.Stop() }
func (s *Ecosystem) IsMining() bool      { return s.miner.Mining() }
func (s *Ecosystem) Miner() *miner.Miner { return s.miner }

func (s *Ecosystem) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *Ecosystem) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *Ecosystem) TxPool() *core.TxPool               { return s.txPool }
func (s *Ecosystem) EventMux() *event.TypeMux           { return s.eventMux }
func (s *Ecosystem) Engine() consensus.Engine           { return s.engine }
func (s *Ecosystem) DPOSEngine() consensus.DPOSEngine   { return s.blockchain.DPOSEngine() }
func (s *Ecosystem) ChainDb() ehcdb.Database            { return s.chainDb }
func (s *Ecosystem) IsListening() bool                  { return true } // Always listening
func (s *Ecosystem) EthVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *Ecosystem) NetVersion() uint64                 { return s.networkId }
func (s *Ecosystem) Downloader() *downloader.Downloader { return s.protocolManager.downloader }
func (s *Ecosystem) CA() *ca.Identity                   { return s.ca }
func (s *Ecosystem) MsgCenter() *mc.Center              { return s.msgcenter }
func (s *Ecosystem) SignHelper() *signhelper.SignHelper { return s.signHelper }
func (s *Ecosystem) ReElection() *reelection.ReElection { return s.reelection }
func (s *Ecosystem) HD() *hd.HD                         { return s.hd }
func (s *Ecosystem) TopNode() *topnode.TopNodeService   { return s.topNode }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *Ecosystem) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.SubProtocols
	}
	return append(s.protocolManager.SubProtocols, s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// Ecosystem protocol implementation.
func (s *Ecosystem) Start(srvr *p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = ehcapi.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers
	if s.config.LightServ > 0 {
		if s.config.LightPeers >= srvr.MaxPeers {
			return fmt.Errorf("invalid peer config: light peer count (%d) >= total peer count (%d)", s.config.LightPeers, srvr.MaxPeers)
		}
		maxPeers -= s.config.LightPeers
	}
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}
	//s.broadTx.Start()//YY
	return nil
}
func (s *Ecosystem) FetcherNotify(hash common.Hash, number uint64) {
	ids := ca.GetRolesByGroup(common.RoleValidator | common.RoleBroadcast)
	for _, id := range ids {
		peer := s.protocolManager.Peers.Peer(id.String()[:16])
		if peer == nil {
			log.Info("==========YY===========", "get PeerID is nil by Validator ID:id", id.String(), "Peers:", s.protocolManager.Peers.peers)
			continue
		}
		s.protocolManager.fetcher.Notify(id.String()[:16], hash, number, time.Now(), peer.RequestOneHeader, peer.RequestBodies)
	}
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Ecosystem protocol.
func (s *Ecosystem) Stop() error {
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}
	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	s.broadTx.Stop() //YY
	close(s.shutdownChan)

	return nil
}
