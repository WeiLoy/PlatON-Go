// Copyright 2021 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashkey-chain/hashkey-chain/common/sort"
	"math"
	"math/big"
	"sync"

	"github.com/hashkey-chain/hashkey-chain/x/gov"

	"github.com/hashkey-chain/hashkey-chain/common/hexutil"

	"github.com/hashkey-chain/hashkey-chain/rlp"

	"github.com/hashkey-chain/hashkey-chain/crypto"

	"github.com/hashkey-chain/hashkey-chain/core/snapshotdb"

	"github.com/hashkey-chain/hashkey-chain/p2p/discover"

	"github.com/hashkey-chain/hashkey-chain/x/staking"

	"github.com/hashkey-chain/hashkey-chain/common"
	"github.com/hashkey-chain/hashkey-chain/common/vm"
	"github.com/hashkey-chain/hashkey-chain/core/types"
	"github.com/hashkey-chain/hashkey-chain/log"
	"github.com/hashkey-chain/hashkey-chain/x/reward"
	"github.com/hashkey-chain/hashkey-chain/x/xcom"
	"github.com/hashkey-chain/hashkey-chain/x/xutil"
)

type RewardMgrPlugin struct {
	db      snapshotdb.DB
	nodeID  discover.NodeID
	nodeADD common.NodeAddress
}

const (
	LessThanFoundationYearDeveloperRate    = 100
	AfterFoundationYearDeveloperRewardRate = 50
	AfterFoundationYearFoundRewardRate     = 50
	RewardPoolIncreaseRate                 = 80 // 80% of fixed-issued tokens are allocated to reward pool each year

)

var (
	rewardOnce  sync.Once
	rm          *RewardMgrPlugin = nil
	millisecond                  = 1000
	minutes                      = 60 * millisecond
)

func RewardMgrInstance() *RewardMgrPlugin {
	rewardOnce.Do(func() {
		log.Info("Init Reward plugin ...")
		sdb := snapshotdb.Instance()
		rm = &RewardMgrPlugin{
			db: sdb,
		}
	})
	return rm
}

// BeginBlock does something like check input params before execute transactions,
// in RewardMgrPlugin it does nothing.
func (rmp *RewardMgrPlugin) BeginBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	return nil
}

// EndBlock will handle reward work, if it's time to settle, reward staking. Then reward worker
// for create new block, this is necessary. At last if current block is the last block at the end
// of year, increasing issuance.
func (rmp *RewardMgrPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	blockNumber := head.Number.Uint64()

	if xutil.IsEndOfEpoch(blockNumber) {
		if err := rmp.runIncreaseIssuance(blockHash, head, state); nil != err {
			return err
		}
	}

	return nil
}

// Confirmed does nothing
func (rmp *RewardMgrPlugin) Confirmed(nodeId discover.NodeID, block *types.Block) error {
	return nil
}

func (rmp *RewardMgrPlugin) SetCurrentNodeID(nodeId discover.NodeID) {
	rmp.nodeID = nodeId
	add, err := xutil.NodeId2Addr(rmp.nodeID)
	if err != nil {
		panic(err)
	}
	rmp.nodeADD = add
}

func (rmp *RewardMgrPlugin) isLessThanFoundationYear(thisYear uint32) bool {
	if thisYear < xcom.PlatONFoundationYear()-1 {
		return true
	}
	return false
}

func (rmp *RewardMgrPlugin) addPlatONFoundation(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) {
	platonFoundationIncr := percentageCalculation(currIssuance, uint64(allocateRate))
	state.AddBalance(xcom.PlatONFundAccount(), platonFoundationIncr)
}

func (rmp *RewardMgrPlugin) addCommunityDeveloperFoundation(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) {
	developerFoundationIncr := percentageCalculation(currIssuance, uint64(allocateRate))
	state.AddBalance(xcom.CDFAccount(), developerFoundationIncr)
}
func (rmp *RewardMgrPlugin) addRewardPoolIncreaseIssuance(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) {
	rewardpoolIncr := percentageCalculation(currIssuance, uint64(allocateRate))
	state.AddBalance(vm.RewardManagerPoolAddr, rewardpoolIncr)
}

// increaseIssuance used for increase issuance at the end of each year
func (rmp *RewardMgrPlugin) increaseIssuance(thisYear, lastYear uint32, state xcom.StateDB, blockNumber uint64, blockHash common.Hash) error {
	var currIssuance *big.Int
	//issuance increase
	{
		histIssuance := GetHistoryCumulativeIssue(state, lastYear)
		increaseIssuanceRatio, err := gov.GovernIncreaseIssuanceRatio(blockNumber, blockHash)
		if nil != err {
			log.Error("Failed to increaseIssuance, call GovernIncreaseIssuanceRatio is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"histIssuance", histIssuance, "err", err)
			return err
		}

		tmp := new(big.Int).Mul(histIssuance, big.NewInt(int64(increaseIssuanceRatio)))
		currIssuance = tmp.Div(tmp, big.NewInt(10000))

		// Restore the cumulative issue at this year end
		histIssuance.Add(histIssuance, currIssuance)
		SetYearEndCumulativeIssue(state, thisYear, histIssuance)
		log.Debug("Call EndBlock on reward_plugin: increase issuance", "thisYear", thisYear, "addIssuance", currIssuance, "hit", histIssuance)

	}
	rewardpoolIncr := percentageCalculation(currIssuance, uint64(RewardPoolIncreaseRate))
	state.AddBalance(vm.RewardManagerPoolAddr, rewardpoolIncr)
	lessBalance := new(big.Int).Sub(currIssuance, rewardpoolIncr)
	if rmp.isLessThanFoundationYear(thisYear) {
		log.Debug("Call EndBlock on reward_plugin: increase issuance to developer", "thisYear", thisYear, "developBalance", lessBalance)
		rmp.addCommunityDeveloperFoundation(state, lessBalance, LessThanFoundationYearDeveloperRate)
	} else {
		log.Debug("Call EndBlock on reward_plugin: increase issuance to developer and hskchain", "thisYear", thisYear, "develop and hskchain Balance", lessBalance)
		rmp.addCommunityDeveloperFoundation(state, lessBalance, AfterFoundationYearDeveloperRewardRate)
		rmp.addPlatONFoundation(state, lessBalance, AfterFoundationYearFoundRewardRate)
	}
	balance := state.GetBalance(vm.RewardManagerPoolAddr)
	SetYearEndBalance(state, thisYear, balance)
	return nil
}

func (rmp *RewardMgrPlugin) getBlockMinderAddress(blockHash common.Hash, head *types.Header) (discover.NodeID, common.NodeAddress, error) {
	if blockHash == common.ZeroHash {
		return rmp.nodeID, rmp.nodeADD, nil
	}
	pk := head.CachePublicKey()
	if pk == nil {
		return discover.ZeroNodeID, common.ZeroNodeAddr, errors.New("failed to get the public key of the block producer")
	}
	return discover.PubkeyID(pk), crypto.PubkeyToNodeAddress(*pk), nil
}

// Calculation percentage ,  input 100,10    cal:  100*10/100 = 10
func percentageCalculation(mount *big.Int, rate uint64) *big.Int {
	ratio := new(big.Int).Mul(mount, big.NewInt(int64(rate)))
	return new(big.Int).Div(ratio, common.Big100)
}

// SetYearEndCumulativeIssue used for set historical cumulative increase at the end of the year
func SetYearEndCumulativeIssue(state xcom.StateDB, year uint32, total *big.Int) {
	yearEndIncreaseKey := reward.GetHistoryIncreaseKey(year)
	state.SetState(vm.RewardManagerPoolAddr, yearEndIncreaseKey, total.Bytes())
}

// GetHistoryCumulativeIssue used for get the cumulative issuance of a certain year in history
func GetHistoryCumulativeIssue(state xcom.StateDB, year uint32) *big.Int {
	var issue = new(big.Int)
	histIncreaseKey := reward.GetHistoryIncreaseKey(year)
	bIssue := state.GetState(vm.RewardManagerPoolAddr, histIncreaseKey)
	log.Trace("show history cumulative issue", "lastYear", year, "amount", issue.SetBytes(bIssue))
	return issue.SetBytes(bIssue)
}

func SetYearEndBalance(state xcom.StateDB, year uint32, balance *big.Int) {
	yearEndBalanceKey := reward.HistoryBalancePrefix(year)
	state.SetState(vm.RewardManagerPoolAddr, yearEndBalanceKey, balance.Bytes())
}

func GetYearEndBalance(state xcom.StateDB, year uint32) *big.Int {
	var balance = new(big.Int)
	yearEndBalanceKey := reward.HistoryBalancePrefix(year)
	bBalance := state.GetState(vm.RewardManagerPoolAddr, yearEndBalanceKey)
	log.Trace("show balance of reward pool at last year end", "lastYear", year, "amount", balance.SetBytes(bBalance))
	return balance.SetBytes(bBalance)
}

func (rmp *RewardMgrPlugin) runIncreaseIssuance(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	yes, err := xcom.IsYearEnd(blockHash, head.Number.Uint64())
	if nil != err {
		log.Error("Failed to execute runIncreaseIssuance function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return err
	}
	if yes {
		yearNumber, err := LoadChainYearNumber(blockHash, rmp.db)
		if nil != err {
			return err
		}
		yearNumber++
		if err := StorageChainYearNumber(blockHash, rmp.db, yearNumber); nil != err {
			log.Error("Failed to execute runIncreaseIssuance function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return err
		}
		incIssuanceTime, err := xcom.LoadIncIssuanceTime(blockHash, rmp.db)
		if nil != err {
			return err
		}
		if err := rmp.increaseIssuance(yearNumber, yearNumber-1, state, head.Number.Uint64(), blockHash); nil != err {
			return err
		}
		if err := xcom.StorageIncIssuanceTime(blockHash, rmp.db, incIssuanceTime+int64(xcom.AdditionalCycleTime()*uint64(minutes))); nil != err {
			log.Error("storage incIssuanceTime fail", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return err
		}
		incIssuanceTime, err = xcom.LoadIncIssuanceTime(blockHash, rmp.db)
		if nil != err {
			return err
		}
		log.Info("Call CalcEpochReward, increaseIssuance successful", "currBlockNumber", head.Number, "currBlockHash", blockHash, "currBlockTime", head.Time,
			"yearNumber", yearNumber, "incIssuanceTime", incIssuanceTime, "yearEndBalance", GetYearEndBalance(state, yearNumber))
	}
	return nil
}

func StorageChainYearNumber(hash common.Hash, snapshotDB snapshotdb.DB, yearNumber uint32) error {
	if err := snapshotDB.Put(hash, reward.ChainYearNumberKey, common.Uint32ToBytes(yearNumber)); nil != err {
		log.Error("Failed to execute StorageChainYearNumber function", "hash", hash.TerminalString(), "yearNumber", yearNumber, "err", err)
		return err
	}
	return nil
}

func LoadChainYearNumber(hash common.Hash, snapshotDB snapshotdb.DB) (uint32, error) {
	chainYearNumberByte, err := snapshotDB.Get(hash, reward.ChainYearNumberKey)
	if nil != err {
		if err == snapshotdb.ErrNotFound {
			log.Info("Data obtained for the first year", "hash", hash.TerminalString())
			return 0, nil
		}
		log.Error("Failed to execute LoadChainYearNumber function", "hash", hash.TerminalString(), "key", string(reward.ChainYearNumberKey), "err", err)
		return 0, err
	}
	return common.BytesToUint32(chainYearNumberByte), nil
}
