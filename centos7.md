



```
curl -sL -o /etc/yum.repos.d/khara-nodejs.repo https://copr.fedoraproject.org/coprs/khara/nodejs/repo/epel-7/khara-nodejs-epel-7.repo

```

```
 yum install -y nodejs nodejs-npm
```

通过npm安装react-tools

```
npm install -g react-tools
```

安装Mesos

```
wget http://www.apache.org/dist/mesos/1.3.0/mesos-1.3.0.tar.gz
tar -zxf mesos-1.3.0.tar.gz
```

http://mesos.apache.org/gettingstarted/