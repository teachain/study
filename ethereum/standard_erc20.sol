pragma solidity ^0.5.8;

//基于ERC20编写的一个代币合约
contract StandardERC20{

    uint256 private _totalSupply;

    string private _name;                   //名称，例如"My test token"

    uint8 private _decimals;               //返回token使用的小数点后几位。比如如果设置为3，就是支持0.001表示.

    string private _symbol;               //token简称,like MTT

    mapping (address => uint256) private _balances; //账户与余额
    
    //这个结构很关键
    mapping (address => mapping (address => uint256)) private  _allowed; //允许的金额


    constructor(uint256 initialAmount, string memory tokenName, uint8 decimalUnits, string memory tokenSymbol) public {
        // 设置初始总量
        _totalSupply = initialAmount * 10 ** uint256(decimalUnits);         
        // 初始token数量给予消息发送者，因为是构造函数，所以这里也是合约的创建者
        _balances[msg.sender] = _totalSupply;  
        _name = tokenName;  
        _decimals = decimalUnits;  
        _symbol =tokenSymbol;
    }

    //返回ERC20代币的名字
    function name() public view returns(string memory){
        return  _name;
    }
    //返回代币的简称，例如：MTT
    function symbol() public view returns(string memory){
        return _symbol;
    }
    //返回token使用的小数点后几位。比如如果设置为3，就是支持0.001表示。
    function decimals() public view returns (uint8){
        return _decimals;
    }
    //返回token的总供应量
    function totalSupply() public view returns (uint){
        return _totalSupply;
    }
    //返回某个地址(账户)的账户余额
    function balanceOf(address _owner) public view returns (uint){
        return _balances[_owner];
    }
    //从代币合约的调用者地址上转移_value的数量token到的地址_to，并且必须触发Transfer事件。
    function transfer(address to, uint value) public returns (bool){
        //默认totalSupply 不会超过最大值 (2^256 - 1).
        //如果随着时间的推移将会有新的token生成，则可以用下面这句避免溢出的异常

        //转账者的代币余额必须大于将要转的数量，并且转的数量必须大于0
        require(_balances[msg.sender] >=value && _balances[to] + value > _balances[to]);
        
        //不能是空地址
        require(to != address(0x0));

        //从消息发送者账户中减去token数量value
        _balances[msg.sender] -= value;
        
        //往接收账户增加token数量value
        _balances[to] += value;

        //触发转币交易事件
        emit Transfer(msg.sender, to, value);

    }
    //从地址from发送数量为value的token到地址to,必须触发Transfer事件。
    //transferFrom方法用于允许合同代理某人转移token。条件是from账户必须经过了approve。
    function transferFrom(address from, address to, uint value) public returns (bool){
        //转账者的代币余额必须大于将要转的数量，并且转的数量必须大于0
        require(_balances[from] >= value && _balances[to] + value > _balances[to]);

        require(_allowed[from][msg.sender]>value);
        //不能是空地址
        require(to != address(0x0));
        //从消息发送者账户中减去token数量value
        _balances[from] -= value;
        //往接收账户增加token数量value
        _balances[to] += value;
        //消息发送者可以从账户from中转出的数量减少value(授权数量)
        _allowed[from][msg.sender] -= value;
        //触发转币交易事件
        emit Transfer(from, to, value);
    }
    //批准，授权
    //意思是调用者允许spender动用多少的额度
    //例如msg.sender是老板，spender是会计，
    //老板调用该方法给会计批准了1万块钱
    //会计就可以从这一万块的额度给员工发5k薪水（调用transferFrom（老板，员工，5000））
    function approve(address spender, uint value) public returns (bool){
        _allowed[msg.sender][spender] = value;
        emit Approval(msg.sender,spender, value);
        return true;
    }
    //返回_spender仍然被允许从_owner提取的金额。
    function allowance(address owner, address spender) public view returns (uint){
        return _allowed[owner][spender];//允许_spender从_owner中转出的token数
    }

   // approve是授权第三方（比如某个服务合约）从发送者账户转移代币，然后通过 transferFrom() 函数来执行具体的转移操作。

   // 账户A有1000个ETH，想允许B账户随意调用他的100个ETH，过程如下：

   // A账户按照以下形式调用approve函数approve(B,100)

   // B账户想用这100个ETH中的10个ETH给C账户，调用transferFrom(A, C, 10)

   // 调用allowance(A, B)可以查看B账户还能够调用A账户多少个token

    //在代币被转移时触发
    event Transfer(address indexed from, address indexed to, uint value);
    //在调用approve方法时触发
    event Approval(address indexed owner, address indexed spender, uint value);
}