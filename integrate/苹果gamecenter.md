##gamecenter接入##

关于登录授权处理注意点：

 Games that support multitasking should take special note of this behavior. When your game moves into the background, the player may launch the Game Center app and sign out. Also, another player might sign in before control is returned to your app. Whenever your game moves to the foreground, it may need to disable its Game Center features when there is no longer an authenticated player or it may need to refresh its internal state based on the identity of the new local player.

意思是特别注意处理在游戏再次进入到前台时，需要检测玩家是否已经sign out.



Apple为大家接入GameCenter提供了GameKit.framework，所以需要以下操作

* 在项目的Build Phases的Link Binary With Libraries中添加GameKit.framework
* 在需要使用GameCenter的类中都要导入GameKit.h
* 在.h文件中加入协议 GKGameCenterControllerDelegate

