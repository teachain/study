## Composer

Composer 是 PHP 的一个依赖管理工具。它允许你申明项目所依赖的代码库，它会在你的项目中为你安装他们。

Composer 不是一个包管理器。是的，它涉及 "packages" 和 "libraries"，但它在每个项目的基础上进行管理，在你项目的某个目录中（例如 vendor ）进行安装。默认情况下它不会在全局安装任何东西。因此，这仅仅是一个依赖管理。

##### 安装

* for Linux / Unix / OSX

  ```
  curl -sS https://getcomposer.org/installer | sudo php
  ```

​      然后当前文件夹下应该会有一个composer.phar文件，然后把它放到全局路径里或者在PATH里加入它的位置

      ```
mv composer.phar /usr/local/bin/composer

#检测版本
composer --version
      ```

* windows

  ```
  https://getcomposer.org/  下载安装即可
  ```

在项目个根目录创建 `composer.json` 文件

```
composer install
```

打开 `autoload_classmap.php` 你会发现你想引入的文件已经加入进来.

如果你项目有增加新的类和依赖内容要加载，你需要重新生成一次自动加载文件，通过 `composer dump-autoload` 命令来实现

在项目的起始文件 `index.php`

```
<?php

// 这里引入了自动加载文件
require 'vendor/autoload.php';

```

### 关于composer.lock

看一下我们的composer.lock,这是个运行了composer install或者update之后就存在的文件。
composer.json是用来配置我们要引入的库，那这个composer.lock是用来干嘛的呢？答案是，用来精确指定版本。
比如我们一开始手动添加版本的时候，添加的是1.2.*,这样我们在拉取的时候，可能会拉去1.2.1或者1.2.4，但加入他官方的版本升级了，之后使用这个项目的人，install的时候拉去的可能就是1.2.5了，说不定有一些不兼容的问题。于是compsoer便设计了一个lock文件，lock文件里会有完整的包信息和版本号。
关于compser.lock的具体情景：

1. composer install 的时候会首先检查lock文件是否存在，如果存在，就按照composer.lock文件指定的版本，取拉取对应的版本，如果不存在，就拉取默认的版本，并且创建composer.lock文件
2. 如果在composer.json文件中修改了版本号，一定要执行composer update，这样，composer.lock文件就会被更新。
3. 团队开发的时候，一定要把composer.lock文件上传到仓库里，这样才能保证团队里的每一个成员的版本都是一致的。

### composer常用命令汇总

1. `composer install`
2. `composer update`
3. `composer require 库名 [版本]` 引入新库或者修改原来库版本
4. `composer selfupdate` composer自身更新
5. `composer self-update --rollback` 更新之后回退到上次的版本
6. `composer create-project` 从现有的包里创建一个新的项目 如：`composer create-project topthink/think=5.0.* tp5 --prefer-dist` 就可以直接下载一个thinkphp5.0的代码
7. `composer search 包名` 比如 `composer search monolog`就可以搜索composer的官方网站 https://packagist.org/ 上的有的包
8. `composer config` composer配置命令，比如全局设置composer的国内镜像 `composer config -g repo.packagist composer https://packagist.phpcomposer.com`