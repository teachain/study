##composer的使用##
<font color="red">composer.phar 是一个可执行文件。</font>

在mac osx下，我们输入以下命令，从输出来看，发现还没有安装composer

```
yangdamins-MacBook-Pro:~ yangdamin$ php --version
PHP 5.5.36 (cli) (built: May 29 2016 01:07:06) 
Copyright (c) 1997-2015 The PHP Group
Zend Engine v2.5.0, Copyright (c) 1998-2015 Zend Technologies
yangdamins-MacBook-Pro:~ yangdamin$ php composer.phar
Could not open input file: composer.phar
yangdamins-MacBook-Pro:~ yangdamin$ 

```

我们将composer安装到指定的/usr/local/bin目录下,因为权限问题，我们不能通过 -- --install-dir=/usr/local/bin指定安装目录，所以我们的选择是先将composer下载到<font color="red">当前的用户目录下</font>，然后再通过mv来实现全局安装。以下是我执行的命令

```
curl -sS https://getcomposer.org/installer | php #这是下载

sudo mv composer.phar /usr/local/bin/composer #移动并重命名

```

这时候我们在terminal中输入composer,会得到以下的输出,如果你没有将composer进行重命名，保持composer.phar的话，那么就需要输入php composer.phar

```
yangdamins-MacBook-Pro:~ yangdamin$ composer

   ______
  / ____/___  ____ ___  ____  ____  ________  _____
 / /   / __ \/ __ `__ \/ __ \/ __ \/ ___/ _ \/ ___/
/ /___/ /_/ / / / / / / /_/ / /_/ (__  )  __/ /
\____/\____/_/ /_/ /_/ .___/\____/____/\___/_/
                    /_/
Composer version 1.2.0 2016-07-19 01:28:52

Usage:
  command [options] [arguments]

Options:
  -h, --help                     Display this help message
  -q, --quiet                    Do not output any message
  -V, --version                  Display this application version
      --ansi                     Force ANSI output
      --no-ansi                  Disable ANSI output
  -n, --no-interaction           Do not ask any interactive question
      --profile                  Display timing and memory usage information
      --no-plugins               Whether to disable plugins.
  -d, --working-dir=WORKING-DIR  If specified, use the given directory as working directory.
  -v|vv|vvv, --verbose           Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug

Available commands:
  about           Short information about Composer
  archive         Create an archive of this composer package
  browse          Opens the package's repository URL or homepage in your browser.
  clear-cache     Clears composer's internal package cache.
  clearcache      Clears composer's internal package cache.
  config          Set config options
  create-project  Create new project from a package into given directory.
  depends         Shows which packages cause the given package to be installed
  diagnose        Diagnoses the system to identify common errors.
  dump-autoload   Dumps the autoloader
  dumpautoload    Dumps the autoloader
  exec            Execute a vendored binary/script
  global          Allows running commands in the global composer dir ($COMPOSER_HOME).
  help            Displays help for a command
  home            Opens the package's repository URL or homepage in your browser.
  info            Show information about packages
  init            Creates a basic composer.json file in current directory.
  install         Installs the project dependencies from the composer.lock file if present, or falls back on the composer.json.
  licenses        Show information about licenses of dependencies
  list            Lists commands
  outdated        Shows a list of installed packages that have updates available, including their latest version.
  prohibits       Shows which packages prevent the given package from being installed
  remove          Removes a package from the require or require-dev
  require         Adds required packages to your composer.json and installs them
  run-script      Run the scripts defined in composer.json.
  search          Search for packages
  self-update     Updates composer.phar to the latest version.
  selfupdate      Updates composer.phar to the latest version.
  show            Show information about packages
  status          Show a list of locally modified packages
  suggests        Show package suggestions
  update          Updates your dependencies to the latest version according to composer.json, and updates the composer.lock file.
  validate        Validates a composer.json and composer.lock
  why             Shows which packages cause the given package to be installed
  why-not         Shows which packages prevent the given package from being installed
 
```
 
 
 
 得到类似以上的输出，表明我们的机器上composer已经可以全局使用了。
 
##使用composer的关键##

<font color="red">定义一个名叫composer.json文件</font>

composer的命令严重依赖该文件来进行依赖库的安装。

<font color="red">以下提到的composer在没有重命名的前提下都是指php composer.phar</font>

要解决和下载依赖(composer的主要作用)

执行以下命令(前提是已经进入到指定的目录，并且在指定的目录下存在composer.json文件)

```
composer install

```





