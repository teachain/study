##微信接入（Android）##

* 首先在开发者应用登记页面进行应用的创建和设置，然后提交审核，<font color="red">只有审核通过的应用才能进行开发</font>
* 下载微信终端开发工具包，其中的libammsdk.jar是必须的，该jar包用于实现与微信的通信。



支付(Midas) 独立于MSDK存在


offerId支付时使用，安卓的offerId为手机QQ的appid


Pf 支付需要使用到的字段, 用于数据分析使用, pf的组成为: 唤起平台_账号体系-注册渠道-操作系统-安装渠道-账号体系-appid-openid.

pfKey 支付使用

在MSDKLibrary/jni/CommonFiles/WGPlatform.h中包含了所有的接口说明(java与C++的接口是对应的)


<font color="red"> 为了防止游戏用测试环境上线, SDK内检测到游戏使用测试环境或者开发环境时, 会Toast出类似: “You are using http://msdktest.qq.com”这样的提示, 游戏切换成正式环境域名以后此提示自动消失.</font>

注意：<font color="red">在 assets/msdkconfig.ini 中配置访问的需要访问的MSDK后台环境</font>



游戏在调用WGLogin后可以开始一个倒计时, 倒计时完毕如果没有收到回调则算作超时, 让用户回到登录界面。倒计时推荐时间为30s

在没有授权的情况下，在登录按钮的点击事件的处理函数中调用WGLogin完成授权登录

处理自动登录WGLoginWithLocalInfo

此接口用于已经登录过的游戏, 在用户再次进入游戏时使用, 游戏启动时先调用此接口, 此接
口会尝试到后台验证票据并通过OnLoginNotify将结果回调给游戏。

在主Activity的onCreate里面MSDK初始化以后调用WGLoginWithLocalInfo完成游戏被拉起时的自动登录。（也就是如果已经授权过了，则有可能成功，否则，它是失败的。）

































