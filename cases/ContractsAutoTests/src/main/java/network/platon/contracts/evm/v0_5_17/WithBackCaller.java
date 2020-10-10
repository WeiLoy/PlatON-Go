package network.platon.contracts.evm.v0_5_17;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.Utf8String;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class WithBackCaller extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610d72806100206000396000f3fe608060405234801561001057600080fd5b50600436106100615760003560e01c80621e257c146100665780630687590a146100aa57806308c2938b14610185578063400f6a601461020857806356ea18ab14610256578063de583cfa1461033b575b600080fd5b6100a86004803603602081101561007c57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610359565b005b610183600480360360408110156100c057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001906401000000008111156100fd57600080fd5b82018360208201111561010f57600080fd5b8035906020019184600183028401116401000000008311171561013157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506104ae565b005b61018d61078a565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101cd5780820151818401526020810190506101b2565b50505050905090810190601f1680156101fa5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102546004803603604081101561021e57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061082c565b005b6103396004803603606081101561026c57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001906401000000008111156102a957600080fd5b8201836020820111156102bb57600080fd5b803590602001918460018302840111640100000000831117156102dd57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290803590602001909291905050506109b0565b005b610343610c8f565b6040518082815260200191505060405180910390f35b8073ffffffffffffffffffffffffffffffffffffffff166055603c604051602401808360ff1681526020018260ff168152602001925050506040516020818303038152906040527f771602f7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b60208310610441578051825260208201915060208101905060208303925061041e565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d80600081146104a3576040519150601f19603f3d011682016040523d82523d6000602084013e6104a8565b606091505b50505050565b600060608373ffffffffffffffffffffffffffffffffffffffff1683604051602401808060200180602001838103835260058152602001807f68656c6c6f000000000000000000000000000000000000000000000000000000815250602001838103825284818151815260200191508051906020019080838360005b8381101561054557808201518184015260208101905061052a565b50505050905090810190601f1680156105725780820380516001836020036101000a031916815260200191505b5093505050506040516020818303038152906040527fae49cd9c000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b602083106106285780518252602082019150602081019050602083039250610605565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461068a576040519150601f19603f3d011682016040523d82523d6000602084013e61068f565b606091505b50915091508161069e57600080fd5b8080602001905160208110156106b357600080fd5b81019080805160405193929190846401000000008211156106d357600080fd5b838201915060208201858111156106e957600080fd5b825186600182028301116401000000008211171561070657600080fd5b8083526020830192505050908051906020019080838360005b8381101561073a57808201518184015260208101905061071f565b50505050905090810190601f1680156107675780820380516001836020036101000a031916815260200191505b5060405250505060019080519060200190610783929190610c98565b5050505050565b606060018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108225780601f106107f757610100808354040283529160200191610822565b820191906000526020600020905b81548152906001019060200180831161080557829003601f168201915b5050505050905090565b600060608373ffffffffffffffffffffffffffffffffffffffff1683604051602401808281526020019150506040516020818303038152906040527f68875570000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b6020831061090857805182526020820191506020810190506020830392506108e5565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461096a576040519150601f19603f3d011682016040523d82523d6000602084013e61096f565b606091505b50915091508161097e57600080fd5b80806020019051602081101561099357600080fd5b810190808051906020019092919050505060008190555050505050565b600060608473ffffffffffffffffffffffffffffffffffffffff168385604051602401808060200180602001838103835260088152602001807f68656c6c6f676173000000000000000000000000000000000000000000000000815250602001838103825284818151815260200191508051906020019080838360005b83811015610a48578082015181840152602081019050610a2d565b50505050905090810190601f168015610a755780820380516001836020036101000a031916815260200191505b5093505050506040516020818303038152906040527fae49cd9c000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b60208310610b2b5780518252602082019150602081019050602083039250610b08565b6001836020036101000a03801982511681845116808217855250505050505090500191505060006040518083038160008787f1925050503d8060008114610b8e576040519150601f19603f3d011682016040523d82523d6000602084013e610b93565b606091505b509150915081610ba257600080fd5b808060200190516020811015610bb757600080fd5b8101908080516040519392919084640100000000821115610bd757600080fd5b83820191506020820185811115610bed57600080fd5b8251866001820283011164010000000082111715610c0a57600080fd5b8083526020830192505050908051906020019080838360005b83811015610c3e578082015181840152602081019050610c23565b50505050905090810190601f168015610c6b5780820380516001836020036101000a031916815260200191505b5060405250505060019080519060200190610c87929190610c98565b505050505050565b60008054905090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610cd957805160ff1916838001178555610d07565b82800160010185558215610d07579182015b82811115610d06578251825591602001919060010190610ceb565b5b509050610d149190610d18565b5090565b610d3a91905b80821115610d36576000816000905550600101610d1e565b5090565b9056fea265627a7a723158204ecfaec5af137ad9deb6f7811b8e158f2faee417e0ef72821bd5bd506ca8daaf64736f6c63430005110032";

    public static final String FUNC_CALLADDLTEST = "callAddlTest";

    public static final String FUNC_CALLDOUBLELTEST = "callDoublelTest";

    public static final String FUNC_CALLGETNAMETEST = "callgetNameTest";

    public static final String FUNC_CALLGETNAMETESTWITHGAS = "callgetNameTestWithGas";

    public static final String FUNC_GETSTRINGRESULT = "getStringResult";

    public static final String FUNC_GETUINTRESULT = "getuintResult";

    protected WithBackCaller(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected WithBackCaller(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> callAddlTest(String other) {
        final Function function = new Function(
                FUNC_CALLADDLTEST, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(other)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callDoublelTest(String other, BigInteger a) {
        final Function function = new Function(
                FUNC_CALLDOUBLELTEST, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(other), 
                new Uint256(a)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callgetNameTest(String other, String name) {
        final Function function = new Function(
                FUNC_CALLGETNAMETEST, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(other), 
                new Utf8String(name)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callgetNameTestWithGas(String other, String name, BigInteger gasValue) {
        final Function function = new Function(
                FUNC_CALLGETNAMETESTWITHGAS, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(other), 
                new Utf8String(name),
                new Uint256(gasValue)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> getStringResult() {
        final Function function = new Function(FUNC_GETSTRINGRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getuintResult() {
        final Function function = new Function(FUNC_GETUINTRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<WithBackCaller> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(WithBackCaller.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<WithBackCaller> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(WithBackCaller.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static WithBackCaller load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new WithBackCaller(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static WithBackCaller load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new WithBackCaller(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
