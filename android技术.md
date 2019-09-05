mvp 

### retrofit 

Retrofit其实我们可以理解为OkHttp的加强版，它也是一个网络加载框架。底层是使用OKHttp封装的。准确来说,网络请求的工作本质上是OkHttp完成，而 Retrofit 仅负责网络请求接口的封装

### dagger2 

Dagger2是现在非常火的一个依赖注入框架。

依赖注入主要有以下几种方式:

* 构造函数注入
* setter方法注入
* 接口注入
* 依赖注入框架
* 

### rxjava 

RxJava 其实就是提供一套异步编程的 API，这套 API 是基于观察者模式的，而且是链式调用的，所以使用 RxJava 编写的代码的逻辑会非常简洁。

这么来理解这么框架比较好

被观察者，观察者，订阅

用大白话来将就是被观察者发生了某件事，需要观察者做出某个动作，那被观察者就必然要知道观察者，所以是被观察者订阅观察者，这样被观察者发生事情了才能告诉观察者。

所以代码应该是 observable.subscribe(observer);observable为被观察者，observer为观察者。

然后链式调用的话，就一句话，被观察者订阅观察者

```
Observable.create(new ObservableOnSubscribe < Integer > () {
    @Override
    public void subscribe(ObservableEmitter < Integer > e) throws Exception {
        Log.d(TAG, "=========================currentThread name: " + Thread.currentThread().getName());
        e.onNext(1);
        e.onNext(2);
        e.onNext(3);
        e.onComplete();
    }
})
.subscribe(new Observer < Integer > () {
    @Override
    public void onSubscribe(Disposable d) {
        Log.d(TAG, "======================onSubscribe");
    }

    @Override
    public void onNext(Integer integer) {
        Log.d(TAG, "======================onNext " + integer);
    }

    @Override
    public void onError(Throwable e) {
        Log.d(TAG, "======================onError");
    }

    @Override
    public void onComplete() {
        Log.d(TAG, "======================onComplete");
    }
});

```

- ObserveOn

  specify the Scheduler on which an observer will observe this Observable
  指定一个观察者在哪个调度器上观察这个Observable

  也就是指定一个***观察者***在哪个调度器上观察这个被观察者，一般observeOn(AndroidSchedulers.mainThread())：在UI线程中处理返回结果，这什么意思呢，其实是我在主线程上观察这个被观察者，一旦被观察者有数据更改，作为观察者，我是需要更新自身UI

- SubscribeOn

  specify the Scheduler on which an Observable will operate
  指定Observable自身在哪个调度器上执行

  也就是指定***被观察者***在哪个调度器上执行，多次调用SubscribeOn并不会起作用，只有第一次调用时会指定Observable自己在哪个调度器执行。

kotlin，

还有热修复 组件化 插件化