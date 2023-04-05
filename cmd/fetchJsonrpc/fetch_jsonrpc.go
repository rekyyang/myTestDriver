package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	jsoniter "github.com/json-iterator/go"
	jsonrpc "github.com/node-real/go-pkg/jsonrpc2"
	"github.com/node-real/go-pkg/log"
)

const (
	StartBlkNo1             = 0x10
	StartBlkNo2             = 0x10000
	StartBlkNo3             = 0x832087
	StartBlkNo4             = 0x85373a
	StartBlkNoAfterShanghai = 0x85f31a

	BlkRange = 100
)

// EvmtestABI is the input ABI used to generate the binding from.
const EvmtestABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ReceivedFromFallback\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ReceivedFromReceive\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"Response\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"balance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"balanceThis\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"call\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"b\",\"type\":\"string\"}],\"name\":\"check_EncodeDecode\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"check_EtherUnits\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num_\",\"type\":\"uint256\"}],\"name\":\"check_Require_1\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num_\",\"type\":\"uint256\"}],\"name\":\"check_Require_2\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"check_Revert_1\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"check_Revert_2\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"check_TimeUnits\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"k\",\"type\":\"uint256\"}],\"name\":\"check_addmod\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"check_balance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"check_ecrecover\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"k\",\"type\":\"uint256\"}],\"name\":\"check_mulmod\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"check_ripemd160\",\"outputs\":[{\"internalType\":\"bytes20\",\"name\":\"\",\"type\":\"bytes20\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"check_sender_balance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"check_sha256\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"delegatecall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_block_coinbase\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_block_difficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_block_gaslimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"get_block_hash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_block_number\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_block_timestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_gasleft\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_msg_data\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_msg_sender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_msg_sig\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_tx_gasprice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_tx_origin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"num\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"send\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"sendViaCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"sendViaSend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"sendViaTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"staticcall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeC\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeCreationCode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeRuntimeCode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"value\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

// EvmtestBin is the compiled bytecode used for deploying new contracts.
var EvmtestBin = "0x6080604052600060015534801561001557600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555061278a806100656000396000f3fe6080604052600436106200029b5760003560e01c80637a9f75d11162000163578063a815af0911620000c7578063c70cf32d1162000085578063c70cf32d1462000f66578063c7e472461462000fd1578063cfcc3abe14620010eb578063efc9b76b1462001145578063f55332ab14620011db578063fede0c921462001222576200030e565b8063a815af091462000ddb578063b17790701462000e46578063b46300ec1462000ee6578063b69ef8a81462000f0a578063bad49bbd1462000f38576200030e565b80638a4068dd11620001215780638a4068dd1462000c5c57806395e8c5d01462000c7c57806395f146a21462000caa5780639724283f1462000cc457806398ba1d461462000d17578063a1e5855d1462000d45576200030e565b80637a9f75d11462000a965780637f8465291462000ac457806382b4b2351462000af2578063830c29ae1462000bd55780638a0c75481462000c2a576200030e565b80633ea7ef55116200020b578063450925a211620001c9578063450925a214620008a05780634e70b1dc14620008ce578063636e082b14620008fc57806367e404ce1462000951578063731a4f3d14620009ab57806374be48061462000a41576200030e565b80633ea7ef55146200075f5780633f1f12d1146200078d5780633f2a403014620007f95780633fa4f245146200085857806341c0e1b51462000886576200030e565b80632d4ca02611620002595780632d4ca026146200058a5780632d72f306146200068b5780632de12dc814620006a557806333cbec4014620006ff5780633df1abd0146200072d576200030e565b80630252971b146200037b5780631896a12914620004115780631f0e645214620004505780631f6b2719146200048f578063271e6a1014620004e9576200030e565b366200030e577fb3f547d62a45bf96295784e4cc876e7adb29056716a0f9f37bc885fd0569bc403334604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a1005b7f03133419d7ad034f47baed0f1b033b92b6fe8c574a97dd06cc2d4142a339e29f3334604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a1005b3480156200038857600080fd5b506200039362001250565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015620003d5578082015181840152602081019050620003b8565b50505050905090810190601f168015620004035780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156200041e57600080fd5b506200044e600480360360208110156200043757600080fd5b81019080803590602001909291905050506200128d565b005b3480156200045d57600080fd5b506200048d600480360360208110156200047657600080fd5b810190808035906020019092919050505062001308565b005b3480156200049c57600080fd5b50620004a76200131a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b348015620004f657600080fd5b506200050162001322565b604051808315151515815260200180602001828103825283818151815260200191508051906020019080838360005b838110156200054d57808201518184015260208101905062000530565b50505050905090810190601f1680156200057b5780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b3480156200059757600080fd5b506200065760048036036020811015620005b057600080fd5b8101908080359060200190640100000000811115620005ce57600080fd5b820183602082011115620005e157600080fd5b803590602001918460018302840111640100000000831117156200060457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506200154a565b60405180826bffffffffffffffffffffffff19166bffffffffffffffffffffffff1916815260200191505060405180910390f35b3480156200069857600080fd5b50620006a362001644565b005b348015620006b257600080fd5b50620006bd62001649565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156200070c57600080fd5b506200071762001651565b6040518082815260200191505060405180910390f35b3480156200073a57600080fd5b506200074562001659565b604051808215151515815260200191505060405180910390f35b3480156200076c57600080fd5b50620007776200169e565b6040518082815260200191505060405180910390f35b3480156200079a57600080fd5b50620007a5620016a6565b60405180827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b3480156200080657600080fd5b5062000856600480360360408110156200081f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050620016d1565b005b3480156200086557600080fd5b506200087062001821565b6040518082815260200191505060405180910390f35b3480156200089357600080fd5b506200089e62001827565b005b348015620008ad57600080fd5b50620008b862001840565b6040518082815260200191505060405180910390f35b348015620008db57600080fd5b50620008e662001848565b6040518082815260200191505060405180910390f35b3480156200090957600080fd5b506200094f600480360360208110156200092257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506200184e565b005b3480156200095e57600080fd5b50620009696200189b565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b348015620009b857600080fd5b50620009c3620018c1565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101562000a05578082015181840152602081019050620009e8565b50505050905090810190601f16801562000a335780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801562000a4e57600080fd5b5062000a946004803603602081101562000a6757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050620018ee565b005b34801562000aa357600080fd5b5062000aae620019a1565b6040518082815260200191505060405180910390f35b34801562000ad157600080fd5b5062000adc620019a9565b6040518082815260200191505060405180910390f35b34801562000aff57600080fd5b5062000bbf6004803603602081101562000b1857600080fd5b810190808035906020019064010000000081111562000b3657600080fd5b82018360208201111562000b4957600080fd5b8035906020019184600183028401116401000000008311171562000b6c57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050620019c8565b6040518082815260200191505060405180910390f35b34801562000be257600080fd5b5062000c286004803603602081101562000bfb57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505062001a73565b005b34801562000c3757600080fd5b5062000c4262001b5c565b604051808215151515815260200191505060405180910390f35b62000c6662001ba2565b6040518082815260200191505060405180910390f35b34801562000c8957600080fd5b5062000c9462001bf4565b6040518082815260200191505060405180910390f35b34801562000cb757600080fd5b5062000cc262001c06565b005b34801562000cd157600080fd5b5062000d016004803603602081101562000cea57600080fd5b810190808035906020019092919050505062001c74565b6040518082815260200191505060405180910390f35b34801562000d2457600080fd5b5062000d2f62001c7f565b6040518082815260200191505060405180910390f35b34801562000d5257600080fd5b5062000d5d62001c87565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101562000d9f57808201518184015260208101905062000d82565b50505050905090810190601f16801562000dcd5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801562000de857600080fd5b5062000e2c6004803603606081101562000e0157600080fd5b8101908080359060200190929190803590602001909291908035906020019092919050505062001cd4565b604051808215151515815260200191505060405180910390f35b34801562000e5357600080fd5b5062000ea46004803603608081101562000e6c57600080fd5b8101908080359060200190929190803560ff169060200190929190803590602001909291908035906020019092919050505062001d0d565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b62000ef062001dd4565b604051808215151515815260200191505060405180910390f35b34801562000f1757600080fd5b5062000f2262001e12565b6040518082815260200191505060405180910390f35b34801562000f4557600080fd5b5062000f5062001e31565b6040518082815260200191505060405180910390f35b34801562000f7357600080fd5b5062000fb76004803603606081101562000f8c57600080fd5b8101908080359060200190929190803590602001909291908035906020019092919050505062001e39565b604051808215151515815260200191505060405180910390f35b34801562000fde57600080fd5b50620010666004803603604081101562000ff757600080fd5b8101908080359060200190929190803590602001906401000000008111156200101f57600080fd5b8201836020820111156200103257600080fd5b803590602001918460018302840111640100000000831117156200105557600080fd5b909192939192939050505062001e72565b6040518083815260200180602001828103825283818151815260200191508051906020019080838360005b83811015620010ae57808201518184015260208101905062001091565b50505050905090810190601f168015620010dc5780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b348015620010f857600080fd5b506200110362001ee1565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156200115257600080fd5b506200115d62001ee9565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156200119f57808201518184015260208101905062001182565b50505050905090810190601f168015620011cd5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6200122060048036036020811015620011f357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505062001f16565b005b3480156200122f57600080fd5b506200123a62002155565b6040518082815260200191505060405180910390f35b60606040518060400160405260018152806020017f4100000000000000000000000000000000000000000000000000000000000000815250905090565b606481101562001305576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600d8152602001807f6e756d5f206973203c203130300000000000000000000000000000000000000081525060200191505060405180910390fd5b50565b60648110156200131757600080fd5b50565b600032905090565b60006060600060603073ffffffffffffffffffffffffffffffffffffffff166040516024016040516020818303038152906040527f731a4f3d000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b60208310620013f95780518252602082019150602081019050602083039250620013d4565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855afa9150503d80600081146200145b576040519150601f19603f3d011682016040523d82523d6000602084013e62001460565b606091505b5091509150818180602001905160208110156200147c57600080fd5b81019080805160405193929190846401000000008211156200149d57600080fd5b83820191506020820185811115620014b457600080fd5b8251866001820283011164010000000082111715620014d257600080fd5b8083526020830192505050908051906020019080838360005b8381101562001508578082015181840152602081019050620014eb565b50505050905090810190601f168015620015365780820380516001836020036101000a031916815260200191505b506040525050508090509350935050509091565b60006003826040516020018082805190602001908083835b6020831062001587578051825260208201915060208101905060208303925062001562565b6001836020036101000a0380198251168184511680821785525050505050509050019150506040516020818303038152906040526040518082805190602001908083835b60208310620015f05780518252602082019150602081019050602083039250620015cb565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa15801562001633573d6000803e3d6000fd5b5050506040515160601b9050919050565b600080fd5b600041905090565b600044905090565b6000600180146200166657fe5b603c80146200167157fe5b610e1080146200167d57fe5b6201518080146200168a57fe5b62093a8080146200169757fe5b6001905090565b60005a905090565b600080357fffffffff0000000000000000000000000000000000000000000000000000000016905090565b600060608373ffffffffffffffffffffffffffffffffffffffff1683604051602401808281526020019150506040516020818303038152906040527f6466414b000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b60208310620017af57805182526020820191506020810190506020830392506200178a565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d806000811462001811576040519150601f19603f3d011682016040523d82523d6000602084013e62001816565b606091505b509150915050505050565b60045481565b3373ffffffffffffffffffffffffffffffffffffffff16ff5b600047905090565b60025481565b8073ffffffffffffffffffffffffffffffffffffffff166108fc6127109081150290604051600060405180830381858888f1935050505015801562001897573d6000803e3d6000fd5b5050565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b606060405180602001620018d590620022e8565b6020820181038252601f19601f82011660405250905090565b60008173ffffffffffffffffffffffffffffffffffffffff166108fc6127109081150290604051600060405180830381858888f193505050509050806200199d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600e8152602001807f4661696c656420746f2073656e6400000000000000000000000000000000000081525060200191505060405180910390fd5b5050565b60003a905090565b60003373ffffffffffffffffffffffffffffffffffffffff1631905090565b60006002826040518082805190602001908083835b6020831062001a025780518252602082019150602081019050602083039250620019dd565b6001836020036101000a038019825116818451168082178552505050505050905001915050602060405180830381855afa15801562001a45573d6000803e3d6000fd5b5050506040513d602081101562001a5b57600080fd5b81019080805190602001909291905050509050919050565b600060608273ffffffffffffffffffffffffffffffffffffffff1661271060405180600001905060006040518083038185875af1925050503d806000811462001ad9576040519150601f19603f3d011682016040523d82523d6000602084013e62001ade565b606091505b50915091508162001b57576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600e8152602001807f4661696c656420746f2073656e6400000000000000000000000000000000000081525060200191505060405180910390fd5b505050565b60006001801462001b6957fe5b64e8d4a51000801462001b7857fe5b66038d7ea4c68000801462001b8957fe5b670de0b6b3a7640000801462001b9b57fe5b6001905090565b60003373ffffffffffffffffffffffffffffffffffffffff166108fc6103e89081150290604051600060405180830381858888f1935050505015801562001bed573d6000803e3d6000fd5b5047905090565b600042421462001c0057fe5b42905090565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600d8152602001807f73696d706c79206661696c65640000000000000000000000000000000000000081525060200191505060405180910390fd5b600081409050919050565b600047905090565b60606000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050905090565b600080828486028162001ce357fe5b0690506000838062001cf157fe5b858709905081811462001d0057fe5b6001925050509392505050565b6000808560405160200180807f19457468657265756d205369676e6564204d6573736167653a0a333200000000815250601c0182815260200191505060405160208183030381529060405280519060200120905060018186868660405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa15801562001dbf573d6000803e3d6000fd5b50505060206040510351915050949350505050565b60003373ffffffffffffffffffffffffffffffffffffffff166108fc6103e89081150290604051600060405180830381858888f19350505050905090565b60003373ffffffffffffffffffffffffffffffffffffffff1631905090565b600045905090565b600080828486018162001e4857fe5b0690506000838062001e5657fe5b858708905081811462001e6557fe5b6001925050509392505050565b600060608062001ec78686868080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506200215d565b905062001ed481620021f7565b9250925050935093915050565b600033905090565b60606040518060200162001efd90620022f6565b6020820181038252601f19601f82011660405250905090565b600060608273ffffffffffffffffffffffffffffffffffffffff163461138890607b60405160240180806020018360ff168152602001828103825260088152602001807f63616c6c20666f6f000000000000000000000000000000000000000000000000815250602001925050506040516020818303038152906040527f24ccab8f000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b6020831062002036578051825260208201915060208101905060208303925062002011565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381858888f193505050503d80600081146200209b576040519150601f19603f3d011682016040523d82523d6000602084013e620020a0565b606091505b50915091507f13848c3e38f8886f3f5d2ad9dff80d8092c2bbb8efd5b887a99c2c6cfc09ac2a8282604051808315151515815260200180602001828103825283818151815260200191508051906020019080838360005b8381101562002114578082015181840152602081019050620020f7565b50505050905090810190601f168015620021425780820380516001836020036101000a031916815260200191505b50935050505060405180910390a1505050565b600043905090565b606082826040516020018083815260200180602001828103825283818151815260200191508051906020019080838360005b83811015620021ac5780820151818401526020810190506200218f565b50505050905090810190601f168015620021da5780820380516001836020036101000a031916815260200191505b509350505050604051602081830303815290604052905092915050565b600060608280602001905160408110156200221157600080fd5b8101908080519060200190929190805160405193929190846401000000008211156200223c57600080fd5b838201915060208201858111156200225357600080fd5b82518660018202830111640100000000821117156200227157600080fd5b8083526020830192505050908051906020019080838360005b83811015620022a75780820151818401526020810190506200228a565b50505050905090810190601f168015620022d55780820380516001836020036101000a031916815260200191505b5060405250505080905091509150915091565b610218806200230583390190565b610238806200251d8339019056fe60806040526004361061004a5760003560e01c80633fa4f2451461004f5780634e70b1dc1461007a5780636466414b146100a557806367e404ce146100d3578063f2c9ecd81461012a575b600080fd5b34801561005b57600080fd5b50610064610155565b6040518082815260200191505060405180910390f35b34801561008657600080fd5b5061008f61015b565b6040518082815260200191505060405180910390f35b6100d1600480360360208110156100bb57600080fd5b8101908080359060200190929190505050610161565b005b3480156100df57600080fd5b506100e86101b3565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561013657600080fd5b5061013f6101d9565b6040518082815260200191505060405180910390f35b60025481565b60005481565b8060008190555033600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503460028190555050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000805490509056fea2646970667358221220d07137050df08e4bb1e8f3cf832b9d912d8209b52aab087f1968d34545d9ab3c64736f6c63430006040033608060405234801561001057600080fd5b50610218806100206000396000f3fe60806040526004361061004a5760003560e01c80633fa4f2451461004f5780634e70b1dc1461007a5780636466414b146100a557806367e404ce146100d3578063f2c9ecd81461012a575b600080fd5b34801561005b57600080fd5b50610064610155565b6040518082815260200191505060405180910390f35b34801561008657600080fd5b5061008f61015b565b6040518082815260200191505060405180910390f35b6100d1600480360360208110156100bb57600080fd5b8101908080359060200190929190505050610161565b005b3480156100df57600080fd5b506100e86101b3565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561013657600080fd5b5061013f6101d9565b6040518082815260200191505060405180910390f35b60025481565b60005481565b8060008190555033600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503460028190555050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000805490509056fea2646970667358221220d07137050df08e4bb1e8f3cf832b9d912d8209b52aab087f1968d34545d9ab3c64736f6c63430006040033a26469706673582212208c2321c7c1cd1ca9caae605ca18ea412defca26199e5564921df0d0fa7146f1664736f6c63430006040033"

var (
	alchemyUrl              = "https://eth-goerli.g.alchemy.com/v2/docs-demo"
	nodeRealProdUrl         = "https://eth-goerli.nodereal.io/v1/d32dc1e5d7554d04832cbf8dbda2c0ff"
	nodeRealQaUrl           = "https://eth-goerli.nodereal.cc/v1/d9fa5c156a3c4102a55af9e97f6eb88f"
	clientAlchemy, _        = jsonrpc.NewClient(jsonrpc.WithURLEndpoint("alchemy_goerli", []string{alchemyUrl}))
	clientNodeReal, _       = jsonrpc.NewClient(jsonrpc.WithURLEndpoint("nodereal_goerli", []string{nodeRealProdUrl}))
	clientNodeRealTracer, _ = jsonrpc.NewClient(jsonrpc.WithURLEndpoint("nodereal_goerli", []string{nodeRealQaUrl}))
)

type RealBlock struct {
	Hash        common.Hash      `json:"hash"             gencodec:"required"`
	ParentHash  common.Hash      `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash      `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address   `json:"miner"            gencodec:"required"`
	Root        common.Hash      `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash      `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash      `json:"receiptsRoot"     gencodec:"required"`
	Bloom       types.Bloom      `json:"logsBloom"        gencodec:"required"`
	Difficulty  hexutil.Big      `json:"difficulty"       gencodec:"required"`
	Number      hexutil.Big      `json:"number"           gencodec:"required"`
	GasLimit    hexutil.Uint64   `json:"gasLimit"         gencodec:"required"`
	GasUsed     hexutil.Uint64   `json:"gasUsed"          gencodec:"required"`
	Time        hexutil.Uint64   `json:"timestamp"        gencodec:"required"`
	Extra       hexutil.Bytes    `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash      `json:"mixHash"`
	Nonce       types.BlockNonce `json:"nonce"`
	Size        hexutil.Uint64   `json:"size"`

	// belows are not belong to header
	TotalDifficulty *hexutil.Big `json:"totalDifficulty,omitempty"  gencodec:"required"`

	// BaseFee was added by EIP-1559 and is ignored in legacy headers.
	BaseFee *hexutil.Big `json:"baseFeePerGas,omitempty"`

	Txs    []RPCTransaction `json:"transactions,omitempty"`
	Uncles []common.Hash    `json:"uncles,omitempty"`

	WithdrawalsHash common.Hash         `json:"withdrawalsRoot,omitempty"`
	Withdrawals     []*types.Withdrawal `json:"withdrawals,omitempty"`
}

type TransactionArgs struct {
	From                 *common.Address `json:"from,omitempty"`
	To                   *common.Address `json:"to,omitempty"`
	Gas                  *hexutil.Uint64 `json:"gas,omitempty"`
	GasPrice             *hexutil.Big    `json:"gasPrice,omitempty"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas,omitempty"`
	Value                *hexutil.Big    `json:"value,omitempty"`
	Nonce                *hexutil.Uint64 `json:"nonce,omitempty"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	// Issue detail: https://github.com/ethereum/go-ethereum/issues/15628
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash        *common.Hash      `json:"blockHash"`
	BlockNumber      *hexutil.Big      `json:"blockNumber"`
	From             common.Address    `json:"from"`
	Gas              hexutil.Uint64    `json:"gas"`
	GasPrice         *hexutil.Big      `json:"gasPrice"`
	GasFeeCap        *hexutil.Big      `json:"maxFeePerGas,omitempty"`
	GasTipCap        *hexutil.Big      `json:"maxPriorityFeePerGas,omitempty"`
	Hash             common.Hash       `json:"hash"`
	Input            hexutil.Bytes     `json:"input"`
	Nonce            hexutil.Uint64    `json:"nonce"`
	To               *common.Address   `json:"to"`
	TransactionIndex *hexutil.Uint64   `json:"transactionIndex"`
	Value            *hexutil.Big      `json:"value"`
	Type             hexutil.Uint64    `json:"type"`
	Accesses         *types.AccessList `json:"accessList,omitempty"`
	ChainID          *hexutil.Big      `json:"chainId,omitempty"`
	V                *hexutil.Big      `json:"v"`
	R                *hexutil.Big      `json:"r"`
	S                *hexutil.Big      `json:"s"`
}

type Content struct {
	Req *jsonrpc.Request
	Rsp *jsonrpc.Response
}

type Block struct {
	Header       types.Header
	Transactions []*common.Hash
}

type TraceFilterMode string

type TraceFilterRequest struct {
	FromBlock   *hexutil.Uint64   `json:"fromBlock,omitempty"`
	ToBlock     *hexutil.Uint64   `json:"toBlock,omitempty"`
	FromAddress []*common.Address `json:"fromAddress,omitempty"`
	ToAddress   []*common.Address `json:"toAddress,omitempty"`
	Mode        TraceFilterMode   `json:"mode,omitempty"`
	After       *uint64           `json:"after,omitempty"`
	Count       *uint64           `json:"count,omitempty"`
}

func main() {
	os.Mkdir("trace_block", os.ModePerm)
	fetchTraceBlock(StartBlkNoAfterShanghai, BlkRange)
	//
	os.Mkdir("trace_replayBlockTransactions", os.ModePerm)
	fetchTraceReplayBlock(StartBlkNoAfterShanghai, 10)
	//
	os.Mkdir("txs", os.ModePerm)
	fetchTransaction(StartBlkNoAfterShanghai, 10)
	//
	os.Mkdir("trace_transaction", os.ModePerm)
	fetchTraceTransaction()
	//
	os.Mkdir("trace_replayTransaction", os.ModePerm)
	fetchTraceReplayTransaction()
	//
	os.Mkdir("trace_call", os.ModePerm)

	os.Mkdir("trace_get", os.ModePerm)
	fetchTraceGet(100, 10)

	os.Mkdir("trace_filter", os.ModePerm)
	fetchTraceFilter(StartBlkNo3, 10)

	//os.Mkdir("trace_call", os.ModePerm)
	//fetchTraceCall(StartBlkNo4, 10)

	validateJsonRpc([]string{
		"trace_block",
		"trace_replayBlockTransactions",
		"trace_transaction",
		"trace_replayTransaction",
		"trace_get",
		"trace_filter",
		"trace_call",
	})

}

func validateJsonRpc(folders []string) {
	failedPath := make([]string, 0)
	errPath := make([]string, 0)

	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"

	for _, folder := range folders {
		files, err := os.ReadDir(folder)
		if err != nil {
			panic(err.Error())
		}
		for _, file := range files {
			pth := path.Join(folder, file.Name())
			fp, err := os.Open(pth)
			if err != nil {
				log.Errorf(err.Error())
				errPath = append(errPath, pth)
				continue
			}
			scanner := bufio.NewScanner(fp)
			buf := make([]byte, 1024*1024*32)
			scanner.Buffer(buf, 1024*1024*32)
			scanner.Scan()

			if err := scanner.Err(); err != nil {
				fmt.Println(err)
				errPath = append(errPath, pth)
				continue
			}
			content := Content{}
			err = jsoniter.Unmarshal(scanner.Bytes(), &content)
			if err != nil {
				fmt.Println(err.Error())
				errPath = append(errPath, pth)
				continue
			}

			respAct, err := clientNodeRealTracer.Call(context.Background(), content.Req, jsonrpc.CallWithHeader(hdr))
			if err != nil {
				fmt.Println(err.Error())
				errPath = append(errPath, pth)
				continue
			}
			rawRespExp, _ := jsoniter.MarshalToString(content.Rsp)
			rawRespAct, _ := jsoniter.MarshalToString(respAct)
			if rawRespExp != rawRespAct {
				failedPath = append(failedPath, pth)
				fmt.Printf("failed case : %s\n", pth)
			} else {
				fmt.Printf("passed case : %s\n", pth)
			}

			fp.Close()
		}
	}
}

// trace_block
func fetchTraceBlock(bnStart, bnRange int) {
	method := "trace_block"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"
	//clientAlchemy, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("alchemy_goerli", []string{"https://eth-goerli.g.alchemy.com/v2/docs-demo"}))
	// block
	// from 0x10 to 0x1000
	for bn := bnStart; bn < bnStart+bnRange; bn++ {
		hx := "0x" + strconv.FormatInt(int64(bn), 16)
		req := jsonrpc.NewRequest(bn, method, hx)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}
}

// trace_replayBlock
func fetchTraceReplayBlock(bnStart, bnRange int) {
	method := "trace_replayBlockTransactions"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"
	//clientAlchemy, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("alchemy_goerli", []string{"https://eth-goerli.g.alchemy.com/v2/docs-demo"}))
	// block
	// from 0x10 to 0x1000
	for bn := bnStart; bn < bnStart+bnRange; bn++ {
		hx := "0x" + strconv.FormatInt(int64(bn), 16)
		fmt.Println(hx)
		req := jsonrpc.NewRequest(bn, method, hx, []string{"trace"})
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	for bn := bnStart; bn < bnStart+bnRange; bn++ {
		hx := "0x" + strconv.FormatInt(int64(bn), 16)
		fmt.Println(hx)
		req := jsonrpc.NewRequest(bn, method, hx, []string{"stateDiff"})
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}
}

// get transactions
func fetchTransaction(bnStart, bnRange int) {
	//noderealUrl := "https://eth-goerli.nodereal.cc/v1/f381061f86f04e2a9490b0986be10a98"
	//clientNodereal, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("nodereal_goerli", []string{noderealUrl}))

	txs := make([]*common.Hash, 0)
	methodGetBlockByNumber := "eth_getBlockByNumber"

	// select some blocks to get transactions
	for bn := bnStart; bn < bnStart+bnRange; bn++ {
		hx := "0x" + strconv.FormatInt(int64(bn), 16)
		fmt.Println(hx)
		req := jsonrpc.NewRequest(bn, methodGetBlockByNumber, hx, false)
		resp, err := clientNodeReal.Call(context.Background(), req)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		//writeJson(fn, resp)
		blk := &Block{}
		err = jsoniter.Unmarshal(resp.Result, blk)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		for _, tx := range blk.Transactions {
			txs = append(txs, tx)
		}
		time.Sleep(200 * time.Millisecond)
	}

	fn := fmt.Sprintf("txs/txs.json")
	writeJson(fn, txs)
}

func fetchTraceTransaction() {
	method := "trace_transaction"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"
	//clientAlchemy, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("alchemy_goerli", []string{"https://eth-goerli.g.alchemy.com/v2/docs-demo"}))

	txs := make([]common.Hash, 0)

	filePtr, err := os.Open("txs/txs.json")
	if err != nil {
		panic(err)
	}
	decoder := jsoniter.NewDecoder(filePtr)
	decoder.Decode(&txs)

	// from 0x10 to 0x1000
	for _, tx := range txs {
		req := jsonrpc.NewRequest("114514", method, tx)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	for _, tx := range txs {
		req := jsonrpc.NewRequest("114514", method, tx)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}
}

func fetchTraceReplayTransaction() {
	method := "trace_replayTransaction"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"
	//clientAlchemy, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("alchemy_goerli", []string{"https://eth-goerli.g.alchemy.com/v2/docs-demo"}))

	txs := make([]common.Hash, 0)

	filePtr, err := os.Open("txs/txs.json")
	if err != nil {
		panic(err)
	}
	decoder := jsoniter.NewDecoder(filePtr)
	decoder.Decode(&txs)

	// from 0x10 to 0x1000
	for _, tx := range txs {
		req := jsonrpc.NewRequest("114514", method, tx, []string{"trace"})
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	for _, tx := range txs {
		req := jsonrpc.NewRequest("114514", method, tx, []string{"stateDiff"})
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}
}

func fetchTraceGet(txCount, idxCount uint64) {
	method := "trace_get"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"
	//clientAlchemy, _ := jsonrpc.NewClient(jsonrpc.WithURLEndpoint("alchemy_goerli", []string{"https://eth-goerli.g.alchemy.com/v2/docs-demo"}))

	txs := make([]common.Hash, 0)

	filePtr, err := os.Open("txs/txs.json")
	if err != nil {
		panic(err)
	}
	decoder := jsoniter.NewDecoder(filePtr)
	decoder.Decode(&txs)

	// from 0x10 to 0x1000
	for _, tx := range txs {
		txCount--
		if txCount < 0 {
			break
		}
		for idx := uint64(0); idx < idxCount; idx++ {
			req := jsonrpc.NewRequest("114514", method, tx, hexutil.EncodeUint64(idx))
			fn := generateFileName(req)
			resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			writeReqResp(fn, req, resp)
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func fetchTraceFilter(bnStart, bnRange int) {
	method := "trace_filter"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"

	//blks := make([]*types.Block, 0)
	txs := make([]*RPCTransaction, 0)
	for bn := bnStart; bn <= bnStart+bnRange; bn++ {
		bn := hexutil.Uint64(bn).String()
		blk_ := RealBlock{}
		req := jsonrpc.NewRequest(bn, "eth_getBlockByNumber", bn, true)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		err = json.Unmarshal(resp.Result, &blk_)
		if err != nil {
			panic(err.Error())
		}
		//blk := types.NewBlockWithHeader(blk_)
		for _, tx := range blk_.Txs {
			//reqTx := jsonrpc.NewRequest(bn, "eth_getTransactionByHash", tx.Hash)
			//respTx, err := clientAlchemy.Call(context.Background(), reqTx, jsonrpc.CallWithHeader(hdr))
			//if err != nil {
			//	fmt.Println(err.Error())
			//	continue
			//}
			//rpcTx := RPCTransaction{}
			//jsoniter.Unmarshal(respTx.Result, &rpcTx)
			txs = append(txs, &tx)
		}
	}

	// full
	{
		fromBlk := hexutil.Uint64(bnStart)
		toBlk := hexutil.Uint64(bnStart + bnRange)
		filterReq := &TraceFilterRequest{
			FromBlock: &fromBlk,
			ToBlock:   &toBlk,
		}
		req := jsonrpc.NewRequest(114514, method, filterReq)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	// count
	{
		count := uint64(15)
		fromBlk := hexutil.Uint64(bnStart)
		toBlk := hexutil.Uint64(bnStart + bnRange)
		filterReq := &TraceFilterRequest{
			FromBlock: &fromBlk,
			ToBlock:   &toBlk,
			Count:     &count,
		}
		req := jsonrpc.NewRequest(114514, method, filterReq)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	// after
	{
		after := uint64(15)
		fromBlk := hexutil.Uint64(bnStart)
		toBlk := hexutil.Uint64(bnStart + bnRange)
		filterReq := &TraceFilterRequest{
			FromBlock: &fromBlk,
			ToBlock:   &toBlk,
			After:     &after,
		}
		req := jsonrpc.NewRequest(114514, method, filterReq)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	// from to union
	{
		fromBlk := hexutil.Uint64(bnStart)
		toBlk := hexutil.Uint64(bnStart + bnRange)
		filterReq := &TraceFilterRequest{
			FromBlock: &fromBlk,
			ToBlock:   &toBlk,
		}
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[0].From)
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[1].From)
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[4].From)

		filterReq.ToAddress = append(filterReq.ToAddress, txs[0].To)
		filterReq.ToAddress = append(filterReq.ToAddress, txs[2].To)
		filterReq.ToAddress = append(filterReq.ToAddress, txs[8].To)
		req := jsonrpc.NewRequest(114514, method, filterReq)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}

	// from to intersection
	{
		fromBlk := hexutil.Uint64(bnStart)
		toBlk := hexutil.Uint64(bnStart + bnRange)
		filterReq := &TraceFilterRequest{
			FromBlock: &fromBlk,
			ToBlock:   &toBlk,
			Mode:      "intersection",
		}
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[0].From)
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[1].From)
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[4].From)
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[7].From)
		filterReq.FromAddress = append(filterReq.FromAddress, &txs[12].From)

		filterReq.ToAddress = append(filterReq.ToAddress, txs[0].To)
		filterReq.ToAddress = append(filterReq.ToAddress, txs[2].To)
		filterReq.ToAddress = append(filterReq.ToAddress, txs[8].To)
		filterReq.ToAddress = append(filterReq.ToAddress, txs[11].To)
		filterReq.ToAddress = append(filterReq.ToAddress, txs[15].To)
		req := jsonrpc.NewRequest(114514, method, filterReq)
		fn := generateFileName(req)
		resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
		if err != nil {
			fmt.Println(err.Error())
		}
		writeReqResp(fn, req, resp)
		time.Sleep(200 * time.Millisecond)
	}
}

func fetchTraceCall(bnStart, bnRange int) {
	method := "trace_call"
	var hdr = make(map[string]string)
	hdr["Origin"] = "https://docs.alchemy.com"

	from := common.Address{}
	to := common.Address{}
	bts, err := hexutil.Decode("0xec50907ad1d361dfa116b690b67ef33685218523")
	to.SetBytes(bts)
	if err != nil {
		fmt.Println(err)
	}
	dataPrefix := "0x9724283f00000000000000000000000000000000"

	for bn := bnStart; bn < bnStart+bnRange; bn++ {
		for gap := 1; gap <= 256; gap += 8 {
			dataSurfix := hexutil.EncodeUint64(uint64(bn - gap))[2:]
			padding := strings.Repeat("0", 32-len(dataSurfix))
			dataSurfix = padding + dataSurfix
			dataRaw := dataPrefix + dataSurfix
			data, _ := hexutil.Decode(dataRaw)
			callBody := TransactionArgs{
				From: &from,
				To:   &to,
				Data: (*hexutil.Bytes)(&data),
			}
			hx := "0x" + strconv.FormatInt(int64(bn), 16)
			req := jsonrpc.NewRequest(bn, method, callBody, []string{"trace"}, hx)
			fn := generateFileName(req)
			resp, err := clientAlchemy.Call(context.Background(), req, jsonrpc.CallWithHeader(hdr))
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			writeReqResp(fn, req, resp)
			time.Sleep(200 * time.Millisecond)
		}
	}
}

//func fetchTraceCall(bnStart, bnRange, paramRange int) {
//	ethCli, err := ethclient.Dial(nodeRealProdUrl)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	ethCli.CallContract(context.Background(), ethereum.CallMsg{
//		From:       common.Address{},
//		To:         nil,
//		Gas:        0,
//		GasPrice:   nil,
//		GasFeeCap:  nil,
//		GasTipCap:  nil,
//		Value:      nil,
//		Data:       nil,
//		AccessList: nil,
//	}, big.NewInt(int64(bnStart)))
//	abi.ConvertType()
//}

func generateFileName(req *jsonrpc.Request) string {
	fileName := fmt.Sprintf("%s/%s_%s.json", req.Method, req.Method, string(req.Params))
	if len(fileName) > 100 {
		fileName_ := fileName
		h := sha256.New()
		h.Write([]byte(fileName_))
		hash := hex.EncodeToString(h.Sum(nil))
		fileName = fileName[:100]
		fileName = fmt.Sprintf("%s-hash[%s].json", fileName, hash)
	}
	return fileName
}

func loadRespFile(req *jsonrpc.Request) *jsonrpc.Response {
	resp := &jsonrpc.Response{}
	filePtr, err := os.Open(generateFileName(req))
	if err != nil {
		return nil
	}
	defer filePtr.Close()
	decoder := jsoniter.NewDecoder(filePtr)
	err = decoder.Decode(resp)
	if err != nil {
		fmt.Println("解码错误", err.Error())
	} else {
		fmt.Println("解码成功")
	}
	return resp
}

func writeReqResp(fileName string, req *jsonrpc.Request, resp *jsonrpc.Response) {
	writeJson(fileName, Content{
		Req: req,
		Rsp: resp,
	})
}

func writeJson(fileName string, content interface{}) {
	// 创建文件
	filePtr, err := os.Create(fileName)
	if err != nil {
		fmt.Println("文件创建失败", err.Error())
		return
	}
	defer filePtr.Close()
	// 创建Json编码器
	encoder := jsoniter.NewEncoder(filePtr)
	err = encoder.Encode(content)
	if err != nil {
		fmt.Println("编码错误", err.Error())
	} else {
		fmt.Println(fileName)
	}
}

//func validate(dirs []string) {
//	os.Mkdir("diff", os.ModePerm)
//	for _, dir := range dirs {
//		files, err := os.ReadDir(dir)
//		if err != nil {
//			log.Errorf(err.Error())
//			continue
//		}
//		for _, file := range files {
//			fp, err := os.Open(file.Name())
//			if err != nil {
//				log.Errorf(err.Error())
//				continue
//			}
//			content := Content{}
//			decoder := jsoniter.NewDecoder(fp)
//			err = decoder.Decode(&content)
//			if err != nil {
//				log.Errorf(err.Error())
//				continue
//			}
//			req := content.Req
//			exp := content.Rsp
//			act, err := clientNodeRealTracer.Call(context.Background(), req)
//			if err != nil {
//				log.Errorf(err.Error())
//				continue
//			}
//			expJson, err := jsoniter.Marshal(exp)
//			actJson, err := jsoniter.Marshal(act)
//			if bytes.Compare(expJson, actJson) != 0 {
//				log.Errorf("compare failed file [%s], err [%s]", file.Name(), err.Error())
//				writeJson(generateFileName(req), struct {
//					Req    *jsonrpc.Request
//					RspExp *jsonrpc.Response
//					RspAct *jsonrpc.Response
//				}{
//					Req:    req,
//					RspExp: exp,
//					RspAct: act,
//				})
//				continue
//			}
//		}
//	}
//}
