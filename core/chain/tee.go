package chain

import (
	"github.com/CESSProject/sdk-go/core/utils"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/pkg/errors"
)

func (c *chainClient) QueryTeePodr2Puk() ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			println(utils.RecoverError(err))
		}
	}()
	var data TeePodr2Pk

	if !c.IsChainClientOk() {
		c.SetChainState(false)
		return nil, ERR_RPC_CONNECTION
	}
	c.SetChainState(true)

	key, err := types.CreateStorageKey(
		c.metadata,
		TEEWORKER,
		TEEPODR2Pk,
	)
	if err != nil {
		return nil, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return nil, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return nil, ERR_RPC_EMPTY_VALUE
	}

	return []byte(string(data[:])), nil
}

func (c *chainClient) QueryTeeInfoList() ([]TeeWorkerInfo, error) {
	var list []TeeWorkerInfo
	key := createPrefixedKey(TEEWORKER, TEEWORKERMAP)
	keys, err := c.api.RPC.State.GetKeysLatest(key)
	if err != nil {
		return list, errors.Wrap(err, "[GetKeysLatest]")
	}
	set, err := c.api.RPC.State.QueryStorageAtLatest(keys)
	if err != nil {
		return list, errors.Wrap(err, "[QueryStorageAtLatest]")
	}
	for _, elem := range set {
		for _, change := range elem.Changes {
			var teeWorker TeeWorkerInfo
			if err := codec.Decode(change.StorageData, &teeWorker); err != nil {
				println(err)
				continue
			}
			list = append(list, teeWorker)
		}
	}
	return list, nil
}

func (c *chainClient) QueryTeePeerID(puk []byte) (PeerID, error) {
	defer func() {
		if err := recover(); err != nil {
			println(utils.RecoverError(err))
		}
	}()
	var data TeeWorkerInfo

	acc, err := types.NewAccountID(puk)
	if err != nil {
		return PeerID{}, errors.Wrap(err, "[NewAccountID]")
	}

	owner, err := codec.Encode(*acc)
	if err != nil {
		return PeerID{}, errors.Wrap(err, "[EncodeToBytes]")
	}

	if !c.IsChainClientOk() {
		c.SetChainState(false)
		return PeerID{}, ERR_RPC_CONNECTION
	}
	c.SetChainState(true)

	key, err := types.CreateStorageKey(
		c.metadata,
		TEEWORKER,
		TEEWORKERMAP,
		owner,
	)
	if err != nil {
		return PeerID{}, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return PeerID{}, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return PeerID{}, ERR_RPC_EMPTY_VALUE
	}

	return data.PeerId, nil
}