package plugin

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// Provides an API interface to obtain data related to the economic model
type PublicPPOSAPI struct {
	snapshotDB snapshotdb.DB
}

func NewPublicPPOSAPI() *PublicPPOSAPI {
	return &PublicPPOSAPI{snapshotdb.Instance()}
}

// Get node list of zero-out blocks
func (p *PublicPPOSAPI) GetWaitSlashingNodeList() string {
	list, err := slash.getWaitSlashingNodeList(0, common.ZeroHash)
	if nil != err || len(list) == 0 {
		return ""
	}
	return fmt.Sprintf("%+v", list)
}

var nodeList []discover.NodeID

func (p *PublicPPOSAPI) SetValidatorList(nodeListEncode string) error {
	var tempNodeList []discover.NodeID
	decodeByte, err := hexutil.Decode(nodeListEncode)
	if nil != err {
		return err
	}
	if err := rlp.DecodeBytes(decodeByte, &tempNodeList); nil != err {
		panic(err)
	}
	nodeList = tempNodeList
	log.Info("SetValidatorList success", "nodeList", fmt.Sprintf("%+v", nodeList))
	return nil
}
