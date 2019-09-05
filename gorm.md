#### 定义模型

gorm.Model是一个包含了ID,CreaatedAt,UpdateAt,DeleteAt四个字段的GoLang结构体，你可以将它嵌入到你自己的模型中，当然你也可以完全使用自己的模型。(也就是是否嵌入gorm.Model是可选的)

```
// gorm.Model 定义
type Model struct {
  ID        uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt *time.Time
}

// Inject fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` into model `User`
// 将 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`字段注入到`User`模型中
type User struct {
  gorm.Model
  Name string
}
```



#### 表名

* 表名默认就是结构体名称的复数。（可以通过调用db.SingularTable(true)让其不加s）

* 可以通过定义func TableName()string方法来定义表名。

  ```
  type Product struct {
  	gorm.Model
  	Code string
  	Price uint
  	Name string
  }
  
  func (this*Product) TableName()string{
  	return "product"
  }
  ```



#### 列名

列名由字段名称进行下划线分割来生成

## 连接数据库

想要连接数据库，你需要先导入对应数据库的驱动，比如：

```
import _ "github.com/go-sql-driver/mysql"
```

GORM 已经包装了一些驱动，以便更容易的记住导入路径，所以你可以这样导入 mysql 的驱动：

```
import _ "github.com/jinzhu/gorm/dialects/mysql"
// import _ "github.com/jinzhu/gorm/dialects/postgres"
// import _ "github.com/jinzhu/gorm/dialects/sqlite"
// import _ "github.com/jinzhu/gorm/dialects/mssql"
```

## 查询

```
//通过主键查询第一条记录
db.First(&user)
//// SELECT * FROM users ORDER BY id LIMIT 1;

// 随机取一条记录
db.Take(&user)
//// SELECT * FROM users LIMIT 1;

// 通过主键查询最后一条记录
db.Last(&user)
//// SELECT * FROM users ORDER BY id DESC LIMIT 1;

// 拿到所有的记录
db.Find(&users)
//// SELECT * FROM users;

// 查询指定的某条记录(只可在主键为整数型时使用)
db.First(&user, 10)
//// SELECT * FROM users WHERE id = 10;
```