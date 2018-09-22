/* contract*/
contract MyToken {
	
	/* 设置一个数组存储每个账户的代币信息 */
	mapping (address => uint256) public balanceOf;
	
	/* name 代币名称 */
	string public name;
	
	/* symbol 代币图标 */
	string public symbol;
	
	/* decimals 代币小数点位数 */
	uint8 public decimals;
	
	/* event事件，它的作用是提醒客户端发生了这个事件，你会注意到钱包有时候会在右下角弹出信息 */
	event Transfer(address indexed from, address indexed to, uint256 value);
  
   /*合约创建者调用*/
   /* 接收用户输入，实现代币的初始化 */
	construct(uint256 initialSupply, string tokenName, uint8 decimalUnits, string tokenSymbol) {
		
		balanceOf[msg.sender] = initialSupply; 
	
		name = tokenName; 

        symbol = tokenSymbol; 

        decimals = decimalUnits; 

    }
	/* 代币交易的函数 */
	function transfer(address _to, uint256 _value) public {
		/* 检查发送方有没有足够的代币 */
		if (balanceOf[msg.sender] < _value || balanceOf[_to] + _value < balanceOf[_to])
		   throw;
		/* 交易过程，发送方减去代币*/
		balanceOf[msg.sender] -= _value;

        /*接收方增加代币*/
		balanceOf[_to] += _value;

		/* 提醒客户端发生了交易事件 */
		Transfer(msg.sender, _to, _value);
	}
}