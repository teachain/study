//基于ERC20编写的一个代币合约
contract StandardERC20{

    uint256 public totalSupply;

    string public name;                   //名称，例如"My test token"

    uint8 public decimals;               //返回token使用的小数点后几位。比如如果设置为3，就是支持0.001表示.

    string public symbol;               //token简称,like MTT

    mapping (address => uint256) balances; //账户与余额
    
    //这个结构很关键
    mapping (address => mapping (address => uint256)) allowed; //允许的金额


    constructor(uint256 _initialAmount, string _tokenName, uint8 _decimalUnits, string _tokenSymbol) public {
        // 设置初始总量
        totalSupply = _initialAmount * 10 ** uint256(_decimalUnits);         
        // 初始token数量给予消息发送者，因为是构造函数，所以这里也是合约的创建者
        balances[msg.sender] = totalSupply;  
 
        name = _tokenName;  

        decimals = _decimalUnits;  

        symbol = _tokenSymbol;
    }

    //返回ERC20代币的名字
	function name() public constant returns(string name){
        return  name;
	}
	//返回代币的简称，例如：MTT
	function symbol() public constant returns(string symbol){
        return symbol;
	}
	//返回token使用的小数点后几位。比如如果设置为3，就是支持0.001表示。
	function decimals() public constant returns (uint8 decimals){
        return decimals;
	}
	//返回token的总供应量
	function totalSupply() public constant returns (uint totalSupply){
        return totalSupply;
	}
	//返回某个地址(账户)的账户余额
	function balanceOf(address _owner) public constant returns (uint balance){
        return balances[_owner]
	}
	//从代币合约的调用者地址上转移_value的数量token到的地址_to，并且必须触发Transfer事件。
	function transfer(address _to, uint _value) public returns (bool success){
        //默认totalSupply 不会超过最大值 (2^256 - 1).
        //如果随着时间的推移将会有新的token生成，则可以用下面这句避免溢出的异常
        require(balances[msg.sender] >= _value && balances[_to] + _value > balances[_to]);
        require(_to != 0x0);
        balances[msg.sender] -= _value;//从消息发送者账户中减去token数量_value
        balances[_to] += _value;//往接收账户增加token数量_value
        Transfer(msg.sender, _to, _value);//触发转币交易事件

	}
	//从地址_from发送数量为_value的token到地址_to,必须触发Transfer事件。
	//transferFrom方法用于允许合同代理某人转移token。条件是from账户必须经过了approve。
	function transferFrom(address _from, address _to, uint _value) public returns (bool success){
		require(balances[_from] >= _value && balances[_to] + _value > balances[_to]);
		require(_to != 0x0);
        balances[_from] -= _value;//从消息发送者账户中减去token数量_value
        balances[_to] += _value;//往接收账户增加token数量_value
        allowed[_from][msg.sender] -= _value;//消息发送者可以从账户_from中转出的数量减少_value
        Transfer(_from, _to, _value);//触发转币交易事件
	}
	//允许_spender多次取回您的帐户，最高达_value金额。 如果再次调用此函数，它将以_value覆盖当前的余量。
	function approve(address _spender, uint _value) public returns (bool success){
        allowed[msg.sender][_spender] = _value;
        Approval(msg.sender, _spender, _value);
        return true
	}
	//返回_spender仍然被允许从_owner提取的金额。
	function allowance(address _owner, address _spender) public constant returns (uint remaining){
        return allowed[_owner][_spender];//允许_spender从_owner中转出的token数
	}

   // approve是授权第三方（比如某个服务合约）从发送者账户转移代币，然后通过 transferFrom() 函数来执行具体的转移操作。

   // 账户A有1000个ETH，想允许B账户随意调用他的100个ETH，过程如下：

   // A账户按照以下形式调用approve函数approve(B,100)

   // B账户想用这100个ETH中的10个ETH给C账户，调用transferFrom(A, C, 10)

   // 调用allowance(A, B)可以查看B账户还能够调用A账户多少个token

	//在代币被转移时触发
	event Transfer(address indexed _from, address indexed _to, uint _value);
	//在调用approve方法时触发
	event Approval(address indexed _owner, address indexed _spender, uint _value);
}