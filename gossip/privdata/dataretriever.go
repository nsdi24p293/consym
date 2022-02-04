/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	protosgossip "github.com/Yunpeng-J/fabric-protos-go/gossip"
	"github.com/Yunpeng-J/fabric-protos-go/ledger/rwset"
	"github.com/Yunpeng-J/HLF-2.2/core/committer"
	"github.com/Yunpeng-J/HLF-2.2/core/ledger"
	"github.com/Yunpeng-J/HLF-2.2/core/transientstore"
	"github.com/Yunpeng-J/HLF-2.2/gossip/privdata/common"
	"github.com/Yunpeng-J/HLF-2.2/gossip/util"
	"github.com/pkg/errors"
)

//go:generate mockery -dir . -name RWSetScanner -case underscore -output mocks/

// RWSetScanner is the local interface used to generate mocks for foreign interface.
type RWSetScanner interface {
	transientstore.RWSetScanner
}

// StorageDataRetriever defines an API to retrieve private date from the storage.
type StorageDataRetriever interface {
	// CollectionRWSet retrieves for give digest relevant private data if
	// available otherwise returns nil, bool which is true if data fetched from ledger and false if was fetched from transient store, and an error
	CollectionRWSet(dig []*protosgossip.PvtDataDigest, blockNum uint64) (Dig2PvtRWSetWithConfig, bool, error)
}

type dataRetriever struct {
	logger    util.Logger
	store     *transientstore.Store
	committer committer.Committer
}

// NewDataRetriever constructing function for implementation of the
// StorageDataRetriever interface
func NewDataRetriever(channelID string, store *transientstore.Store, committer committer.Committer) StorageDataRetriever {
	return &dataRetriever{
		logger:    logger.With("channel", channelID),
		store:     store,
		committer: committer,
	}
}

// CollectionRWSet retrieves for give digest relevant private data if
// available otherwise returns nil, bool which is true if data fetched from ledger and false if was fetched from transient store, and an error
func (dr *dataRetriever) CollectionRWSet(digests []*protosgossip.PvtDataDigest, blockNum uint64) (Dig2PvtRWSetWithConfig, bool, error) {
	height, err := dr.committer.LedgerHeight()
	if err != nil {
		// if there is an error getting info from the ledger, we need to try to read from transient store
		return nil, false, errors.Wrap(err, "wasn't able to read ledger height")
	}
	if height <= blockNum {
		dr.logger.Debug("Current ledger height ", height, "is below requested block sequence number",
			blockNum, "retrieving private data from transient store")
	}

	if height <= blockNum { // Check whenever current ledger height is equal or below block sequence num.
		results := make(Dig2PvtRWSetWithConfig)
		for _, dig := range digests {
			filter := map[string]ledger.PvtCollFilter{
				dig.Namespace: map[string]bool{
					dig.Collection: true,
				},
			}
			pvtRWSet, err := dr.fromTransientStore(dig, filter)
			if err != nil {
				dr.logger.Errorf("couldn't read from transient store private read-write set, "+
					"digest %+v, because of %s", dig, err)
				continue
			}
			results[common.DigKey{
				Namespace:  dig.Namespace,
				Collection: dig.Collection,
				TxId:       dig.TxId,
				BlockSeq:   dig.BlockSeq,
				SeqInBlock: dig.SeqInBlock,
			}] = pvtRWSet
		}

		return results, false, nil
	}
	// Since ledger height is above block sequence number private data is might be available in the ledger
	results, err := dr.fromLedger(digests, blockNum)
	return results, true, err
}

func (dr *dataRetriever) fromLedger(digests []*protosgossip.PvtDataDigest, blockNum uint64) (Dig2PvtRWSetWithConfig, error) {
	filter := make(map[string]ledger.PvtCollFilter)
	for _, dig := range digests {
		if _, ok := filter[dig.Namespace]; !ok {
			filter[dig.Namespace] = make(ledger.PvtCollFilter)
		}
		filter[dig.Namespace][dig.Collection] = true
	}

	pvtData, err := dr.committer.GetPvtDataByNum(blockNum, filter)
	if err != nil {
		return nil, errors.Errorf("wasn't able to obtain private data, block sequence number %d, due to %s", blockNum, err)
	}

	results := make(Dig2PvtRWSetWithConfig)
	for _, dig := range digests {
		dig := dig
		pvtRWSetWithConfig := &util.PrivateRWSetWithConfig{}
		for _, data := range pvtData {
			if data.WriteSet == nil {
				dr.logger.Warning("Received nil write set for collection tx in block", data.SeqInBlock, "block number", blockNum)
				continue
			}

			// private data doesn't hold rwsets for namespace and collection or
			// belongs to different transaction
			if !data.Has(dig.Namespace, dig.Collection) || data.SeqInBlock != dig.SeqInBlock {
				continue
			}

			pvtRWSet := dr.extractPvtRWsets(data.WriteSet.NsPvtRwset, dig.Namespace, dig.Collection)
			pvtRWSetWithConfig.RWSet = append(pvtRWSetWithConfig.RWSet, pvtRWSet...)
		}

		confHistoryRetriever, err := dr.committer.GetConfigHistoryRetriever()
		if err != nil {
			return nil, errors.Errorf("cannot obtain configuration history retriever, for collection <%s>"+
				" txID <%s> block sequence number <%d> due to <%s>", dig.Collection, dig.TxId, dig.BlockSeq, err)
		}

		configInfo, err := confHistoryRetriever.MostRecentCollectionConfigBelow(dig.BlockSeq, dig.Namespace)
		if err != nil {
			return nil, errors.Errorf("cannot find recent collection config update below block sequence = %d,"+
				" collection name = <%s> for chaincode <%s>", dig.BlockSeq, dig.Collection, dig.Namespace)
		}

		if configInfo == nil {
			return nil, errors.Errorf("no collection config update below block sequence = <%d>"+
				" collection name = <%s> for chaincode <%s> is available ", dig.BlockSeq, dig.Collection, dig.Namespace)
		}
		configs := extractCollectionConfig(configInfo.CollectionConfig, dig.Collection)
		if configs == nil {
			return nil, errors.Errorf("no collection config was found for collection <%s>"+
				" namespace <%s> txID <%s>", dig.Collection, dig.Namespace, dig.TxId)
		}
		pvtRWSetWithConfig.CollectionConfig = configs
		results[common.DigKey{
			Namespace:  dig.Namespace,
			Collection: dig.Collection,
			TxId:       dig.TxId,
			BlockSeq:   dig.BlockSeq,
			SeqInBlock: dig.SeqInBlock,
		}] = pvtRWSetWithConfig
	}

	return results, nil
}

func (dr *dataRetriever) fromTransientStore(dig *protosgossip.PvtDataDigest, filter map[string]ledger.PvtCollFilter) (*util.PrivateRWSetWithConfig, error) {
	results := &util.PrivateRWSetWithConfig{}
	it, err := dr.store.GetTxPvtRWSetByTxid(dig.TxId, filter)
	if err != nil {
		return nil, errors.Errorf("was not able to retrieve private data from transient store, namespace <%s>"+
			", collection name %s, txID <%s>, due to <%s>", dig.Namespace, dig.Collection, dig.TxId, err)
	}
	defer it.Close()

	maxEndorsedAt := uint64(0)
	for {
		res, err := it.Next()
		if err != nil {
			return nil, errors.Errorf("error getting next element out of private data iterator, namespace <%s>"+
				", collection name <%s>, txID <%s>, due to <%s>", dig.Namespace, dig.Collection, dig.TxId, err)
		}
		if res == nil {
			return results, nil
		}
		rws := res.PvtSimulationResultsWithConfig
		if rws == nil {
			dr.logger.Debug("Skipping nil PvtSimulationResultsWithConfig received at block height", res.ReceivedAtBlockHeight)
			continue
		}
		txPvtRWSet := rws.PvtRwset
		if txPvtRWSet == nil {
			dr.logger.Debug("Skipping empty PvtRwset of PvtSimulationResultsWithConfig received at block height", res.ReceivedAtBlockHeight)
			continue
		}

		colConfigs, found := rws.CollectionConfigs[dig.Namespace]
		if !found {
			dr.logger.Error("No collection config was found for chaincode", dig.Namespace, "collection name",
				dig.Namespace, "txID", dig.TxId)
			continue
		}

		configs := extractCollectionConfig(colConfigs, dig.Collection)
		if configs == nil {
			dr.logger.Error("No collection config was found for collection", dig.Collection,
				"namespace", dig.Namespace, "txID", dig.TxId)
			continue
		}

		pvtRWSet := dr.extractPvtRWsets(txPvtRWSet.NsPvtRwset, dig.Namespace, dig.Collection)
		if rws.EndorsedAt >= maxEndorsedAt {
			maxEndorsedAt = rws.EndorsedAt
			results.CollectionConfig = configs
		}
		results.RWSet = append(results.RWSet, pvtRWSet...)
	}
}

func (dr *dataRetriever) extractPvtRWsets(pvtRWSets []*rwset.NsPvtReadWriteSet, namespace string, collectionName string) []util.PrivateRWSet {
	pRWsets := []util.PrivateRWSet{}

	// Iterate over all namespaces
	for _, nsws := range pvtRWSets {
		// and in each namespace - iterate over all collections
		if nsws.Namespace != namespace {
			dr.logger.Debug("Received private data namespace ", nsws.Namespace, " instead of ", namespace, " skipping...")
			continue
		}
		for _, col := range nsws.CollectionPvtRwset {
			// This isn't the collection we're looking for
			if col.CollectionName != collectionName {
				dr.logger.Debug("Received private data collection ", col.CollectionName, " instead of ", collectionName, " skipping...")
				continue
			}
			// Add the collection pRWset to the accumulated set
			pRWsets = append(pRWsets, col.Rwset)
		}
	}

	return pRWsets
}
