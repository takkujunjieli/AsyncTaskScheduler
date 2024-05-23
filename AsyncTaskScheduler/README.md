# asyncflow

## 项目介绍
flowsvr：任务流服务，对外提供任务处理/查询接口
worker：处理某种/多种任务，其中集成了tasksdk提供自动调度，部署于客户端

## 编译&启动

提前按《第一期定调》文档创建表及插入测试数据。

flowsvr编译：在asyncflow/flowsvr下执行make命令

flowsvr执行前需要本地安装MySQL和Redis，并在asyncflow/flowsvr/src/config/config-test.toml 配置

flowsvr执行：编译后的文件在asyncflow/flowsvr/bin/下面，执行命令./flowsvr test


worker编译：在asyncflow/worker下执行make命令
worker执行：编译后的文件在asyncflow/worker/bin/下面，执行命令./worker

## 创建一条测试任务
asyncflow/worker/src/main_test.go文件，TestCreateTask，使用Goland直接运行即可。
