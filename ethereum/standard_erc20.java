
// This file is an automatically generated Java binding. Do not modify as any
// change will likely be lost upon the next re-generation!

package contracts;

import org.ethereum.geth.*;
import org.ethereum.geth.internal.*;


	public class StandardERC20 {
		// ABI is the input ABI used to generate the binding from.
		public final static String ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"initialAmount\",\"type\":\"uint256\"},{\"name\":\"tokenName\",\"type\":\"string\"},{\"name\":\"decimalUnits\",\"type\":\"uint8\"},{\"name\":\"tokenSymbol\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]";

		
			// BYTECODE is the compiled bytecode used for deploying new contracts.
			public final static byte[] BYTECODE = "0x608060405234801561001057600080fd5b506040516107e63803806107e68339818101604052608081101561003357600080fd5b81516020830180519193928301929164010000000081111561005457600080fd5b8201602081018481111561006757600080fd5b815164010000000081118282018710171561008157600080fd5b505060208201516040909201805191949293916401000000008111156100a657600080fd5b820160208101848111156100b957600080fd5b81516401000000008111828201871017156100d357600080fd5b505060ff8516600a0a8702600081815533815260046020908152604090912091909155865191945061010c935060019250860190610139565b506002805460ff191660ff8416179055805161012f906003906020840190610139565b50505050506101d4565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061017a57805160ff19168380011785556101a7565b828001600101855582156101a7579182015b828111156101a757825182559160200191906001019061018c565b506101b39291506101b7565b5090565b6101d191905b808211156101b357600081556001016101bd565b90565b610603806101e36000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063313ce56711610066578063313ce567146101a557806370a08231146101c357806395d89b41146101e9578063a9059cbb146101f1578063dd62ed3e1461021d57610093565b806306fdde0314610098578063095ea7b31461011557806318160ddd1461015557806323b872dd1461016f575b600080fd5b6100a061024b565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100da5781810151838201526020016100c2565b50505050905090810190601f1680156101075780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101416004803603604081101561012b57600080fd5b506001600160a01b0381351690602001356102e0565b604080519115158252519081900360200190f35b61015d610346565b60408051918252519081900360200190f35b6101416004803603606081101561018557600080fd5b506001600160a01b0381358116916020810135909116906040013561034c565b6101ad61045f565b6040805160ff9092168252519081900360200190f35b61015d600480360360208110156101d957600080fd5b50356001600160a01b0316610468565b6100a0610483565b6101416004803603604081101561020757600080fd5b506001600160a01b0381351690602001356104e4565b61015d6004803603604081101561023357600080fd5b506001600160a01b03813581169160200135166105a3565b60018054604080516020601f600260001961010087891615020190951694909404938401819004810282018101909252828152606093909290918301828280156102d65780601f106102ab576101008083540402835291602001916102d6565b820191906000526020600020905b8154815290600101906020018083116102b957829003601f168201915b5050505050905090565b3360008181526005602090815260408083206001600160a01b038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b60005490565b6001600160a01b038316600090815260046020526040812054821180159061038d57506001600160a01b038316600090815260046020526040902054828101115b61039657600080fd5b6001600160a01b038416600090815260056020908152604080832033845290915290205482106103c557600080fd5b6001600160a01b0383166103d857600080fd5b6001600160a01b0380851660008181526004602090815260408083208054889003905593871680835284832080548801905583835260058252848320338452825291849020805487900390558351868152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a39392505050565b60025460ff1690565b6001600160a01b031660009081526004602052604090205490565b60038054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156102d65780601f106102ab576101008083540402835291602001916102d6565b33600090815260046020526040812054821180159061051c57506001600160a01b038316600090815260046020526040902054828101115b61052557600080fd5b6001600160a01b03831661053857600080fd5b336000818152600460209081526040808320805487900390556001600160a01b03871680845292819020805487019055805186815290519293927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a392915050565b6001600160a01b0391821660009081526005602090815260408083209390941682529190915220549056fea265627a7a723058207480d79d86f498ab4625d43023557135bc6286d247326bcbe7193dbe58475dcb64736f6c63430005090032".getBytes();

			// deploy deploys a new Simplechain contract, binding an instance of StandardERC20 to it.
			public static StandardERC20 deploy(TransactOpts auth, SimplechainClient client, BigInt initialAmount, String tokenName, byte decimalUnits, String tokenSymbol) throws Exception {
				Interfaces args = Geth.newInterfaces(4);
				
				  args.set(0, Geth.newInterface()); args.get(0).setBigInt(initialAmount);
				
				  args.set(1, Geth.newInterface()); args.get(1).setString(tokenName);
				
				  args.set(2, Geth.newInterface()); args.get(2).setUint8(decimalUnits);
				
				  args.set(3, Geth.newInterface()); args.get(3).setString(tokenSymbol);
				
				return new StandardERC20(Geth.deployContract(auth, ABI, BYTECODE, client, args));
			}

			// Internal constructor used by contract deployment.
			private StandardERC20(BoundContract deployment) {
				this.Address  = deployment.getAddress();
				this.Deployer = deployment.getDeployer();
				this.Contract = deployment;
			}
		

		// Simplechain address where this contract is located at.
		public final Address Address;

		// Simplechain transaction in which this contract was deployed (if known!).
		public final Transaction Deployer;

		// Contract instance bound to a blockchain address.
		private final BoundContract Contract;

		// Creates a new instance of StandardERC20, bound to a specific deployed contract.
		public StandardERC20(Address address, SimplechainClient client) throws Exception {
			this(Geth.bindContract(address, ABI, client));
		}

		
			

			// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
			//
			// Solidity: function allowance(owner address, spender address) constant returns(uint256)
			public BigInt Allowance(CallOpts opts, Address owner, Address spender) throws Exception {
				Interfaces args = Geth.newInterfaces(2);
				args.set(0, Geth.newInterface()); args.get(0).setAddress(owner);
				args.set(1, Geth.newInterface()); args.get(1).setAddress(spender);
				

				Interfaces results = Geth.newInterfaces(1);
				Interface result0 = Geth.newInterface(); result0.setDefaultBigInt(); results.set(0, result0);
				

				if (opts == null) {
					opts = Geth.newCallOpts();
				}
				this.Contract.call(opts, results, "allowance", args);
				return results.get(0).getBigInt();
				
			}
		
			

			// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
			//
			// Solidity: function balanceOf(_owner address) constant returns(uint256)
			public BigInt BalanceOf(CallOpts opts, Address _owner) throws Exception {
				Interfaces args = Geth.newInterfaces(1);
				args.set(0, Geth.newInterface()); args.get(0).setAddress(_owner);
				

				Interfaces results = Geth.newInterfaces(1);
				Interface result0 = Geth.newInterface(); result0.setDefaultBigInt(); results.set(0, result0);
				

				if (opts == null) {
					opts = Geth.newCallOpts();
				}
				this.Contract.call(opts, results, "balanceOf", args);
				return results.get(0).getBigInt();
				
			}
		
			

			// Decimals is a free data retrieval call binding the contract method 0x313ce567.
			//
			// Solidity: function decimals() constant returns(uint8)
			public byte Decimals(CallOpts opts) throws Exception {
				Interfaces args = Geth.newInterfaces(0);
				

				Interfaces results = Geth.newInterfaces(1);
				Interface result0 = Geth.newInterface(); result0.setDefaultUint8(); results.set(0, result0);
				

				if (opts == null) {
					opts = Geth.newCallOpts();
				}
				this.Contract.call(opts, results, "decimals", args);
				return results.get(0).getUint8();
				
			}
		
			

			// Name is a free data retrieval call binding the contract method 0x06fdde03.
			//
			// Solidity: function name() constant returns(string)
			public String Name(CallOpts opts) throws Exception {
				Interfaces args = Geth.newInterfaces(0);
				

				Interfaces results = Geth.newInterfaces(1);
				Interface result0 = Geth.newInterface(); result0.setDefaultString(); results.set(0, result0);
				

				if (opts == null) {
					opts = Geth.newCallOpts();
				}
				this.Contract.call(opts, results, "name", args);
				return results.get(0).getString();
				
			}
		
			

			// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
			//
			// Solidity: function symbol() constant returns(string)
			public String Symbol(CallOpts opts) throws Exception {
				Interfaces args = Geth.newInterfaces(0);
				

				Interfaces results = Geth.newInterfaces(1);
				Interface result0 = Geth.newInterface(); result0.setDefaultString(); results.set(0, result0);
				

				if (opts == null) {
					opts = Geth.newCallOpts();
				}
				this.Contract.call(opts, results, "symbol", args);
				return results.get(0).getString();
				
			}
		
			

			// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
			//
			// Solidity: function totalSupply() constant returns(uint256)
			public BigInt TotalSupply(CallOpts opts) throws Exception {
				Interfaces args = Geth.newInterfaces(0);
				

				Interfaces results = Geth.newInterfaces(1);
				Interface result0 = Geth.newInterface(); result0.setDefaultBigInt(); results.set(0, result0);
				

				if (opts == null) {
					opts = Geth.newCallOpts();
				}
				this.Contract.call(opts, results, "totalSupply", args);
				return results.get(0).getBigInt();
				
			}
		

		
			// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
			//
			// Solidity: function approve(spender address, value uint256) returns(bool)
			public Transaction Approve(TransactOpts opts, Address spender, BigInt value) throws Exception {
				Interfaces args = Geth.newInterfaces(2);
				args.set(0, Geth.newInterface()); args.get(0).setAddress(spender);
				args.set(1, Geth.newInterface()); args.get(1).setBigInt(value);
				

				return this.Contract.transact(opts, "approve"	, args);
			}
		
			// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
			//
			// Solidity: function transfer(to address, value uint256) returns(bool)
			public Transaction Transfer(TransactOpts opts, Address to, BigInt value) throws Exception {
				Interfaces args = Geth.newInterfaces(2);
				args.set(0, Geth.newInterface()); args.get(0).setAddress(to);
				args.set(1, Geth.newInterface()); args.get(1).setBigInt(value);
				

				return this.Contract.transact(opts, "transfer"	, args);
			}
		
			// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
			//
			// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
			public Transaction TransferFrom(TransactOpts opts, Address from, Address to, BigInt value) throws Exception {
				Interfaces args = Geth.newInterfaces(3);
				args.set(0, Geth.newInterface()); args.get(0).setAddress(from);
				args.set(1, Geth.newInterface()); args.get(1).setAddress(to);
				args.set(2, Geth.newInterface()); args.get(2).setBigInt(value);
				

				return this.Contract.transact(opts, "transferFrom"	, args);
			}
		
	}

