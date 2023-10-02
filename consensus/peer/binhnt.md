# Peer contains 2 service:
    - Delivery Service: 
        + peer use this service to connect to orderer (ordering service) => get new block event
        + After verify block => broadcast block in gossip service
        + call Deliverer to  deliverBlocks
        + deliverBlocks: used to pull out blocks from the ordering service to distributed them across peers
    - Deliver server: 
        + Handle delivery request 
        + call deliverBlocks (same function but different execute) => call ChainManager => GetChain (from DeliverChainManager)
        => DeliverChainManager (contain Peer)
        => Peer contain list of Channel => Return chain
        => Chain => chain.Reader 


# Ledger 
    + FileLedger : 
        + blockStore FileLedgerBlockStore
            + AddBlock(block *cb.Block) error
            + GetBlockchainInfo() (*cb.BlockchainInfo, error)
            + RetrieveBlocks(startBlockNumber uint64) (ledger.+ + ResultsIterator, error)
        + signal 
        => Iterator(startPosition *ab.SeekPosition): fileLedgerIterator{ledger: fl, blockNumber: startingBlockNumber, commonIterator: iterator}

        

    + fileLedgerIterator: 
        + ledger         *FileLedger
        + blockNumber    uint64
        + commonIterator ledger.ResultsIterator

        + 	iterator, err := fl.blockStore.RetrieveBlocks(startingBlockNumber)

    + fileLedgerFactory: Create a ledger 

    + RetrieveBlocks
    

# Block structure 


# CLient => Server: TLS 

# Check block signer 
