#功能
1）实现用户管理，对用户分权限管理所属项目
  管理员账号：创建|删除研发账号、创建项目（包括svn仓库）
  研发账号：用于发布更新项目svn仓库代码
  运维账号：用于发布更新项目、重启服务
2）实现kubernetes平台资源管理
  管理员账号：配置svc服务负载
  研发账号：更新代码和发布，包括docker项目构建、上传镜像仓库、发布kubernetes集群
3）工作流pipeline
  1、创建项目（包括svn仓库地址、docker基础镜像，如php-fpm、kubernetes编排文件yaml）
  目录结构：
  |-build //构建相关目录
  |    |-src //源码目录
  |    |--Dockerfile //docker编译文件，根据选型docker镜像自动生成
  |    |--startup.sh //docker启动脚本，控制启动
  |  
  |-apply //kubernetes编排目录
       |-configmap //配置文件目录
       |-pvc //挂载存储目录
       |-svc //服务负载目录
       |--[项目名称].yaml //项目编译文件
                          //kind种类：Deployment、Statefulset、Daemonset
                          //name名称：项目名称
                          //replicas个数：运行pod个数
                          //nodeSelector调度：idc位置
                          //resources限制：
                          //  cpu
                          //  memory
4）代码结构
|-src
|   |-backend
|   |    |-client //与kubernetes交互client-go
|   |    |-conf //配置文件kube-config位置
|   |    |-api //go-gin接口目录
|   |    |   |-controllers //接口控制层
|   |    |   |-database //数据库配置
|   |    |   |-models //与数据库交互，业务逻辑
|   |    |   |-router //路由定义
|   |    |-middleware //中间件   
|   |    |-pkg //外部代码包 
|   |    |  |-e //消息格式
|   |    |  |-setting //配置
|   |    |  |-util //工具类
|   |    |   
|   |    |-yaml //存放kubernetes生成yaml编排文件
|   |        |-[项目名称]
|   |        |     |-configmap
|   |        |     |-pvc
|   |        |     |-svc
|   |        |     |--[项目].yaml
|   |        |
|   |        |-[项目名称...]
|   |        
|   |-frontend
|        |-vue
|
