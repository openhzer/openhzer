## 关于

`Hzer` 是基于 [Gin](https://github.com/gin-gonic/gin) 进行模块化封装的API开发框架，封装了常用的功能，使用简单，致力于进行快速的业务研发。

供参考学习，线上使用请谨慎！

## 内容

* [x] RESTful API 返回值规范
* [ ] 接口限流
* [ ] 接口跨域
* [ ] 指标记录
* [ ] 接口文档生成
* [ ] 数据库组件
* [ ] 完整的JWT鉴权
* [ ] 定时任务
* [ ] Admin后台
* [ ] 日志记录
* [ ] panic 异常通知

## 使用组件

[Gin](https://github.com/gin-gonic/gin) HTTP服务核心

[Gcache](https://github.com/8treenet/gcache) Gorm缓存驱动

[Light Year Example](https://gitee.com/yinqi/Light-Year-Admin-Template-v4) Admin模板

## 目录说明

```
├─ configs  //配置文件
├─ docs     //文档
├─ internal //服务/业务实现
├─ logs     //日志
├─ pkg      //引用外部包
├─ scripts  //构建脚本相关
├─ static   //静态资源
│   └─ assets    //网页资源
└─ main.go  //代码入口
```

## 注意

此项目正处于开发阶段

