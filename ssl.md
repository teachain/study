### 密钥协商的步骤


1. 客户端连上服务端
2. 服务端发送 CA 证书给客户端
3. 客户端验证该证书的可靠性
4. 客户端从 CA 证书中取出公钥
5. 客户端生成一个随机密钥 k，并用这个公钥加密得到 k'
6. 客户端把 k' 发送给服务端
7. 服务端收到 k' 后用自己的私钥解密得到 k
8. 此时双方都得到了密钥 k，协商完成。

### ◇如何防范偷窥（嗅探）


　　攻击方式1
　　攻击者虽然可以监视网络流量并拿到公钥，但是【无法】通过公钥推算出私钥（这点由 RSA 算法保证）

　　攻击方式2
　　攻击者虽然可以监视网络流量并拿到 k'，但是攻击者没有私钥，【无法解密】 k'，因此也就无法得到 k

### ◇如何防范篡改（假冒身份）


 　　攻击方式1
 　　如果攻击者在第2步篡改数据，伪造了证书，那么客户端在第3步会发现（这点由证书体系保证）

 　　攻击方式2
 　　如果攻击者在第6步篡改数据，伪造了k'，那么服务端收到假的k'之后，解密会失败（这点由 RSA 算法保证）。服务端就知道被攻击了。



