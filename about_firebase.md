firebase 推送服务

 https://console.firebase.google.com/

* 在Android Studio中引入messaging的依赖

```
在项目/app/build.gradle中 添加
compile 'com.google.firebase:firebase-messaging:9.6.0'

```

* 在自己的APP中进行处理消息的相关编码

```
编写一个继承自FirebaseMessagingService的service。通过这个service来进行相关消息处理。
通过Firebase Notification推送的消息包含两种类型的消息：「通知」和「数据」
通知：用于在状态栏自动创建通知时，一些必要的参数。例如标题和文字信息等。

数据：推送过来的附加数据，可以自定义。
处理Firebase Notification推送的消息，有两种情况：「APP处于前台」和「APP处于后台」
APP处于前台运行时：这种情况下，主要通过FirebaseMessagingService中的onMessageReceived(RemoteMessage remoteMessage)方法来处理消息。我们只要重写该方法，就可以实现自己想要的一些功能。

APP处于后台：这种情况下，FirebaseMessagingService会自动根据「通知」信息的一些参数，为我们在状态栏中创建一条通知信息。我们在点击这条通知时，会默认打开我们的APP。而「数据」信息需要通过Intent的getExtra来获取。
```

* 在Firebase控制台中通过Notification发送推送消息