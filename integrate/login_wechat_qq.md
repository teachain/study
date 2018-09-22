##登陆接入##

1、设置需要用户授权的权限

```
// 设置拉起QQ时候需要用户授权的项
WGPlatform.WGSetPermission(WGQZonePermissions.eOPEN_ALL); 
```
<font color="red">该接口在WGPlatform.Initialized接口之后调用。</font>

2、处理授权登陆

调用以下接口拉起登陆界面

```
   WGPlatform::GetInstance()->WGLogin(ePlatform_QQ); 
```

参数可以为ePlatform_QQ或ePlatform_Weixin

在onLoginNotify接口中获取回调参数数据

3、自动登录

使用

```
WGPlatform::GetInstance()->WGLogin(EPlatform.ePlatform_None)
```

替换老接口

```
WGPlatform::GetInstance()->WGLoginWithLocalInfo()
```