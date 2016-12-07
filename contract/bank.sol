pragma solidity ^0.4.6;
contract Bank {
    
    uint64 private count;
    address private owner;
    uint256 private totalAmount;
    
    mapping (address => uint256) clients;
    
    event InitEvt(address sender);
    event DepositEvt(address sender, uint256 amount, uint256 totalAmount, uint64 count);
    event WithdrawEvt(address sender, uint256 amount, uint256 totalAmount, uint64 count);
    
    function Bank() {
      owner = msg.sender;
      count = 0;
      totalAmount = 0;
      InitEvt(owner);
    }
    
    function() payable {
        deposit(); 
    }
    
    function deposit () payable {
        if (clients[msg.sender] == 0){
            //new client
            count += 1;
        }
        totalAmount += msg.value;
        clients[msg.sender] += msg.value;
        DepositEvt(msg.sender, msg.value, totalAmount, count);
    }
    
    function withdraw () {
        uint256 balance = clients[msg.sender];
        if (balance > 0){
            if (msg.sender.send(balance)){
                totalAmount-=balance;
                count-=1;
                delete clients[msg.sender];
                WithdrawEvt(msg.sender, balance, totalAmount, count);
            }else{
                throw;
            }
        }
    }
    
    function getOwner() constant returns (address) {
        return owner;
    }
    
    function getCount() constant returns (uint64){
      return count;
    }
    
    function getTotalAmount() constant returns (uint256){
      return totalAmount;
    }

}