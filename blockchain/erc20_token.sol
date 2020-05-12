pragma solidity ^0.5.9;
//https://eips.ethereum.org/EIPS/eip-20
library SafeMathLib {
  function minus(uint256 a, uint256 b) pure public returns (uint256) {
    require(b <= a);
    return a - b;
  }
  function plus(uint256 a, uint256 b) pure public returns (uint256) {
    uint256 c = a + b;
    require(c>=a && c>=b);
    return c;
  }
}
contract StandardToken {
    using SafeMathLib for uint256;
    string private tokenName;
    string private tokenSymbol;
    uint8 private tokenDecimals = 18;
    mapping (address => uint256) private balances;
    mapping (address => mapping (address => uint256)) private allowed;
    uint256 private tokenInitialSupply;
    uint256 private tokenTotalSupply;

    event Transfer(address indexed _from, address indexed _to, uint256 _value);

    event Approval(address indexed _owner, address indexed _spender, uint256 _value);
    //initial 发行量
    //name    币的名称 比如莱特币
    //symbol  币的符号 LTC
    //decimals 几位小数 18
    constructor(uint256 initial, string memory name, string memory symbol,uint8 decimals) public {
        tokenInitialSupply=initial;
        tokenName=name;
        tokenSymbol=symbol;
        tokenDecimals=decimals;
        tokenTotalSupply = tokenInitialSupply * 10 ** uint256(tokenDecimals);
        balances[msg.sender] = tokenTotalSupply;
    }
   function name() public view returns (string memory){
      return tokenName;
   }
   function symbol() public view returns (string memory){
      return tokenSymbol;
   }
   function decimals() public view returns (uint8){
      return tokenDecimals;
   }
   //返回总发行量
   function totalSupply() public view returns (uint256){
      return tokenTotalSupply;
   }
   //返回某地址的余额
   function balanceOf(address _owner) public view returns (uint256){
      return balances[_owner];
   }
   //转账
   function transfer(address _to, uint256 _value) public returns (bool){
      balances[msg.sender] = balances[msg.sender].minus(_value);
      balances[_to] = balances[_to].plus(_value);
      emit Transfer(msg.sender, _to, _value);
   }
   //从_from地址转账给_to地址
   //msg.sender需要_from事先调用approve授权给msg.sender
   function transferFrom(address _from, address _to, uint256 _value) public returns (bool){
      uint256 _allowance = allowed[_from][msg.sender];
      balances[_to] = balances[_to].plus(_value);
      balances[_from] = balances[_from].minus(_value);
      allowed[_from][msg.sender] = _allowance.minus(_value);
      emit Transfer(_from, _to, _value);
      return true;
   }
   //msg.sender授权给_spender多少额度
   function approve(address _spender, uint256 _value) public returns (bool){
      //授权的额度是任意的，但至于到时能转账多少，则取决于msg.sender有多少额度
      allowed[msg.sender][_spender] = _value;
      emit Approval(msg.sender, _spender, _value);
      return true;
   }
   function allowance(address _owner, address _spender) public view returns (uint256){
      return allowed[_owner][_spender];
   }
 }