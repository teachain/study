
<font color="red">苹果审核app时，仍然在沙盒环境下测试。</font>
苹果验证返回

* 正常情况都是返回status:0,
* {"status"21002,"exception":"java.lang.IllegalArgumentException"}这个应该是提交的数据不是合法base64的，或是说并不是接口要求的json数据。也就是有可能是因为你从客户端提交给你的服务器的时候没有进行urlencode,导致接收到的数据与客户端发送的不一致。

* 如果只返回{"status":21002}，应该是内购破解。


##in app purchase客户端接入##

1、 在工程中引入 storekit.framework 和 #import <StoreKit/StoreKit.h>

获取"产品付费数量等于0这个问题"的原因，原因是你在 ItunesConnect 里的 “Contracts, Tax, and Banking ”没有完成设置账户信息。


#if defined (__cplusplus)
extern "C"
{
#endif
    
    //这里写你的两个函数函数
    
#if defined (__cplusplus)
}
#endif