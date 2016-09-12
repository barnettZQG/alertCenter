# 安装要求
项目 | 要求
---- | ---
go    | 版本1.6.3及以上
beego | 版本1.7.0及以上

# API说明

API | 说明
---- | ----
/api/teams |获取全部团队信息
/api/teamUsers | 获取团队带用户信息
/api/addTeam | 添加团队
/api/receive | 接收报警信息
/api/v1/alerts | 接收prometheus报警信息
/api/alert/handle/?:ID/?:type | 处理报警
/api/alerts/?:receiver/:pageSize/:page | 获取报警信息带分页
/api/ignoreRule/:user"| 获取忽略规则
/api/addIgnoreRule" | 添加忽略规则
/api/ignoreAlert/:mark" | 忽略某报警