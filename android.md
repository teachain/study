##代码片段##

```
//跳转到微信

Intent intent = new Intent(Intent.ACTION_MAIN);
ComponentName cmp = new ComponentName("com.tencent.mm","com.tencent.mm.ui.LauncherUI");

intent.addCategory(Intent.CATEGORY_LAUNCHER);
intent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
intent.setComponent(cmp);
startActivity(intent);
```