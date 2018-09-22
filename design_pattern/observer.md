##观察者模式##

观察者模式是一种对象行为型模式。

定义：
定义对象之间的一种一对多依赖关系，使得每当一个对象状态发生改变时，其相关依赖对象皆得到通知并被自动更新。

观察者的别名

* 发布-订阅模式（Publish/Subscribe）
* 模型-视图模式(Model/View)
* 源-监听器模式(Source/Listener)
* 从属者模式(Dependents)

示例代码：

```
package com.damin.design;

import java.util.ArrayList;

import com.damin.design.Observer;
//被观察目标
public abstract class Subject {
	// 定义一个观察者集合用于存储所有观察者对象
	protected ArrayList<Observer> observers=new ArrayList<Observer>();
	// 注册方法，用于向观察者集合中增加一个观察者
	public void attach(Observer observer) {
		observers.add(observer);
	}
	// 注销方法，用于在观察者集合中删除一个观察者
	public void detach(Observer observer) {
		observers.remove(observer);
	}
	// 声明抽象通知方法
	public abstract void notifyObserver();
}
```

```
package com.damin.design;
import com.damin.design.Subject;
///被观察具体目标
public class ConcreteSubject extends Subject{
	@Override
	public void notifyObserver() {
		for(Object obs:observers) {  
            ((Observer)obs).update();  
        } 
	}

}
```

```
package com.damin.design;
//观察者
public interface Observer {
    public void update();
}

```

```
package com.damin.design;

import com.damin.design.Observer;
//具体观察者
public class ConcreteObserver implements Observer{
	@Override
	public void update() {
		System.out.println("ConcreteObserver update......");
	}

}
```