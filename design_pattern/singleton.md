##单例模式##

单例模式是一种对象创建型模式。

定义：
确保某一个类只有一个实例，而且自行实例化并向整个系统提供这个实例，这个类称为单例类，它提供全局访问的方法。

在单例类的内部实现只生成一个实例，同时它提供一个静态的getInstance()工厂方法，让客户可以访问它的唯一实例，为了防止在外部对其实例化，将其构造方法设计为私有，在单例类内部定义了一个Singleton类型的静态对象，作为外部共享的唯一实例。

示例代码(java):

```
package com.damin.design;
//单例模式有三个要点
public class Singleton {
	//一是某个类只能有一个实例
	private static Singleton instance=null;
	//二是它必须自行创建这个实例,为了防止外部创建，故设置为私有。
	private Singleton(){
		
	}
	//三是它必须自行向整个系统提供这个实例。
	public static Singleton getInstance(){
		if(instance==null){
			instance=new Singleton();
		}
		return instance;
	}
	public void display(){
		System.out.println("Singleton display......");
	}
}
```