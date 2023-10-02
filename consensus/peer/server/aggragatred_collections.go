package server

import "github.com/hyperledger/fabric-protos-go/ledger/rwset"

type seqAndDataModel struct {
	seq       uint64
	dataModel rwset.TxReadWriteSet_DataModel
}

// Below map temporarily stores the private data that have passed the corresponding collection policy.
// outer map is from seqAndDataModel to inner map,
// and innner map is from namespace to []*rwset.CollectionPvtReadWriteSet
type aggregatedCollections map[seqAndDataModel]map[string][]*rwset.CollectionPvtReadWriteSet

// addCollection adds private data based on seq, namespace, and collection.
func (ac aggregatedCollections) addCollection(seqInBlock uint64, dm rwset.TxReadWriteSet_DataModel, namespace string, col *rwset.CollectionPvtReadWriteSet) {
	seq := seqAndDataModel{
		dataModel: dm,
		seq:       seqInBlock,
	}
	if _, exists := ac[seq]; !exists {
		ac[seq] = make(map[string][]*rwset.CollectionPvtReadWriteSet)
	}
	ac[seq][namespace] = append(ac[seq][namespace], col)
}

// asPrivateDataMap converts aggregatedCollections to map[uint64]*rwset.TxPvtReadWriteSet
// as defined in BlockAndPrivateData protobuf message.
func (ac aggregatedCollections) asPrivateDataMap() map[uint64]*rwset.TxPvtReadWriteSet {
	pvtDataMap := make(map[uint64]*rwset.TxPvtReadWriteSet)
	for seq, ns := range ac {
		// create a txPvtReadWriteSet and add collection data to it
		txPvtRWSet := &rwset.TxPvtReadWriteSet{
			DataModel: seq.dataModel,
		}

		for namespaceName, cols := range ns {
			txPvtRWSet.NsPvtRwset = append(txPvtRWSet.NsPvtRwset, &rwset.NsPvtReadWriteSet{
				Namespace:          namespaceName,
				CollectionPvtRwset: cols,
			})
		}

		pvtDataMap[seq.seq] = txPvtRWSet
	}
	return pvtDataMap
}
