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
	"github.com/hashkey-chain/hashkey-chain/x/gov"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashkey-chain/hashkey-chain/x/xcom"

	"github.com/hashkey-chain/hashkey-chain/core/snapshotdb"

	"github.com/hashkey-chain/hashkey-chain/common"

	"github.com/hashkey-chain/hashkey-chain/common/vm"
)

func TestIncreaseIssuance(t *testing.T) {
	var plugin = RewardMgrInstance()

	mockDB := buildStateDB(t)

	initIncreaseIssuanceRatio := xcom.IncreaseIssuanceRatio()
	gov.InitGenesisGovernParam(common.ZeroHash, snapshotdb.Instance(), 2048)

	thisYear, lastYear := uint32(1), uint32(0)

	genesisIssue := new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18))
	SetYearEndCumulativeIssue(mockDB, 0, genesisIssue)

	lastIssue := GetHistoryCumulativeIssue(mockDB, lastYear)

	mockDB.AddBalance(vm.RestrictingContractAddr, genesisIssue)

	if err := plugin.increaseIssuance(thisYear, lastYear, mockDB, 1, common.ZeroHash); nil != err {
		t.Fatal(err)
	}

	newIssue := GetHistoryCumulativeIssue(mockDB, thisYear)

	increaseIssuanceRatio, err := gov.GovernIncreaseIssuanceRatio(1, common.ZeroHash)
	if nil != err {
		t.Fatal(err)
	}

	tmp := new(big.Int).Sub(newIssue, lastIssue)
	assert.Equal(t, increaseIssuanceRatio, initIncreaseIssuanceRatio)
	assert.Equal(t, tmp, new(big.Int).Div(new(big.Int).Mul(lastIssue, big.NewInt(int64(initIncreaseIssuanceRatio))), big.NewInt(int64(10000))))

	if plugin.isLessThanFoundationYear(thisYear) {
		mockDB.GetBalance(xcom.CDFAccount())

	} else {
		mockDB.GetBalance(xcom.CDFAccount())
		mockDB.GetBalance(xcom.PlatONFundAccount())
	}

}

func TestZeroIncreaseIssuance(t *testing.T) {
	var plugin = RewardMgrInstance()

	_, genesis, _ := newChainState()

	mockDB := buildStateDB(t)

	gov.InitGenesisGovernParam(common.ZeroHash, snapshotdb.Instance(), 2048)

	if err := snapshotdb.Instance().NewBlock(blockNumber, genesis.Hash(), common.ZeroHash); nil != err {
		t.Fatal(err)
	}
	defer func() {
		snapshotdb.Instance().Clear()
	}()

	if err := gov.SetGovernParam(gov.ModuleReward, gov.KeyIncreaseIssuanceRatio, "", "0", 0, common.ZeroHash); nil != err {
		t.Fatal(err)
	}

	thisYear, lastYear := uint32(1), uint32(0)

	genesisIssue := new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18))
	SetYearEndCumulativeIssue(mockDB, 0, genesisIssue)

	lastIssue := GetHistoryCumulativeIssue(mockDB, lastYear)

	mockDB.AddBalance(vm.RestrictingContractAddr, genesisIssue)

	if err := plugin.increaseIssuance(thisYear, lastYear, mockDB, 1, common.ZeroHash); nil != err {
		t.Fatal(err)
	}

	newIssue := GetHistoryCumulativeIssue(mockDB, thisYear)

	increaseIssuanceRatio, err := gov.GovernIncreaseIssuanceRatio(1, common.ZeroHash)
	if nil != err {
		t.Fatal(err)
	}

	tmp := new(big.Int).Sub(newIssue, lastIssue)
	assert.Equal(t, increaseIssuanceRatio, uint16(0))
	assert.Equal(t, tmp.Uint64(), uint64(0))

}

func TestCDFAccountOneYearIncreaseIssuance(t *testing.T) {
	var plugin = RewardMgrInstance()

	mockDB := buildStateDB(t)

	gov.InitGenesisGovernParam(common.ZeroHash, snapshotdb.Instance(), 2048)

	thisYear, lastYear := uint32(1), uint32(0)

	genesisIssue := new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18))
	SetYearEndCumulativeIssue(mockDB, 0, genesisIssue)

	lastIssue := GetHistoryCumulativeIssue(mockDB, lastYear)

	mockDB.AddBalance(vm.RestrictingContractAddr, genesisIssue)

	CDFAccountBalance := mockDB.GetBalance(xcom.CDFAccount())
	if err := plugin.increaseIssuance(thisYear, lastYear, mockDB, 1, common.ZeroHash); nil != err {
		t.Fatal(err)
	}

	newIssue := GetHistoryCumulativeIssue(mockDB, thisYear)

	currIssue := new(big.Int).Sub(newIssue, lastIssue)

	currCDFAccountBalance := new(big.Int).Sub(mockDB.GetBalance(xcom.CDFAccount()), CDFAccountBalance)
	rewardpoolIncr := percentageCalculation(currIssue, uint64(RewardPoolIncreaseRate))
	assert.Equal(t, currCDFAccountBalance, new(big.Int).Sub(currIssue, rewardpoolIncr))

}

func TestCDFAccountTenYearIncreaseIssuance(t *testing.T) {
	var plugin = RewardMgrInstance()

	mockDB := buildStateDB(t)

	gov.InitGenesisGovernParam(common.ZeroHash, snapshotdb.Instance(), 2048)

	thisYear, lastYear := uint32(10), uint32(0)

	genesisIssue := new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18))
	SetYearEndCumulativeIssue(mockDB, 0, genesisIssue)

	lastIssue := GetHistoryCumulativeIssue(mockDB, lastYear)

	mockDB.AddBalance(vm.RestrictingContractAddr, genesisIssue)

	CDFAccountBalance := mockDB.GetBalance(xcom.CDFAccount())
	PlatONFundAccountBalance := mockDB.GetBalance(xcom.PlatONFundAccount())
	if err := plugin.increaseIssuance(thisYear, lastYear, mockDB, 1, common.ZeroHash); nil != err {
		t.Fatal(err)
	}

	newIssue := GetHistoryCumulativeIssue(mockDB, thisYear)

	currIssue := new(big.Int).Sub(newIssue, lastIssue)

	currCDFAccountBalance := new(big.Int).Sub(mockDB.GetBalance(xcom.CDFAccount()), CDFAccountBalance)
	currPlatONFundAccountBalance := new(big.Int).Sub(mockDB.GetBalance(xcom.PlatONFundAccount()), PlatONFundAccountBalance)

	lessBalance := new(big.Int).Sub(currIssue, percentageCalculation(currIssue, uint64(RewardPoolIncreaseRate)))
	assert.Equal(t, currCDFAccountBalance, percentageCalculation(lessBalance, uint64(AfterFoundationYearDeveloperRewardRate)))
	assert.Equal(t, currPlatONFundAccountBalance, percentageCalculation(lessBalance, uint64(AfterFoundationYearFoundRewardRate)))

}
