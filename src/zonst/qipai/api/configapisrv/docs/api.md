 <center><h1>配置服务 API 接口</h1></center>

# 1. web 端请求

## 1.1. 下载页面配置 API 接口

### 1.1.1. 获取指定平台的下载页面是否已关闭

地址： /v1/download/page/list/is/closed

方法： POST

接口参数传递方式：form 或 json

参数：

| 名称         | 类型    | 是否必填 | 备注         |
| :----------- | :------ | :------- | :----------- |
| game\_id_list | []int32 | 是       | 平台 ID 列表 |

返回值：

| 名称   | 类型         | 备注                         |
| :----- | :----------- | :--------------------------- |
| errno  | string       | 0 表示无异常，其它表示有异常 |
| errmsg | string       | 异常信息                     |
| data   | []GameStatus | 参考 GameStatus 字段说明     |

GameStatus 字段说明

| 名称    | 类型  | 备注               |
| :------ | :---- | :----------------- |
| game_id | int32 | 平台 ID            |
| close   | bool  | 是否关闭了下载页面 |

```
示例：

```

## 1.2. 游戏配置接口

### 1.2.1. 获取游戏平台列表

地址： /v1/game/forbid/list

方法： POST

接口参数传递方式：json

参数：

| 名称         | 类型    | 是否必填 | 备注         |
| :----------- | :------ | :------- | :----------- |
| game\_id_list | []int32 | 是       | 平台 ID 列表 |

返回值：

| 名称   | 类型   | 备注                         |
| :----- | :----- | :--------------------------- |
| errno  | int32  | 0 表示无异常，其它表示有异常 |
| errmsg | string | 异常信息                     |
| data   | []Data | 数据列表，参考 Data 说明     |

Data 说明

| 名称    | 类型  | 备注                               |
| :------ | :---- | :--------------------------------- |
| game_id | int32 | 游戏平台 ID                        |
| forbid  | int32 | 1 表示该平台已被禁止，0 表示不禁止 |

```
示例：

```

### 1.2.2. 开启/关闭游戏平台

地址： /v1/game/forbid/update

方法： POST

接口参数传递方式：json

参数：

| 名称     | 类型   | 是否必填 | 备注                     |
| :------- | :----- | :------- | :----------------------- |
| game_id  | int32  | 是       | 平台 ID                  |
| forbid   | int32  | 是       | 1 表示禁止，0 表示不禁止 |
| user_id  | int32  | 是       | 用户 ID，一般为工号      |
| username | string | 是       | 用户名                   |

返回值：

| 名称   | 类型   | 备注                         |
| :----- | :----- | :--------------------------- |
| errno  | int32  | 0 表示无异常，其它表示有异常 |
| errmsg | string | 异常信息                     |

### 1.2.3. 开启/关闭游戏下载包

地址： /v1/game/version/status/update

方法： POST

接口参数传递方式：json

参数：

| 名称     | 类型   | 是否必填 | 备注                     |
| :------- | :----- | :------- | :----------------------- |
| open\_package_id  | int32  | 否       | 要打开的包的自增ID，当前状态是关闭  两个参数必须传一个|
| close\_package_id   | int32  | 否       | 要关闭的包的自增ID,当前状态是打开，只剩一个打开的包无法关闭 两个参数必须传一个 |

返回值：

| 名称   | 类型   | 备注                         |
| :----- | :----- | :--------------------------- |
| errno  | int32  | 0 表示无异常，其它表示有异常 |
| errmsg | string | 异常信息                     |



### 1.2.4. 获取当前平台的下载包列表

地址： /v1/game/version/list

方法： POST

接口参数传递方式：json

参数：

| 名称     | 类型   | 是否必填 | 备注                     |
| :------- | :----- | :------- | :----------------------- |
| game_id  | int32  | 是       | 当前平台ID |
| package_type   | string  | 否       | 包的类型 取值 ios 或 android|

返回值：

| 名称   | 类型   | 备注                         |
| :----- | :----- | :--------------------------- |
| errno  | int32  | 0 表示无异常，其它表示有异常 |
| errmsg | string | 异常信息                     |
|data|[]OnGameVersionListResponce|数据列表|

结构体OnGameVersionListResponce

| 名称   | 类型   | 备注                         |
| :----- | :----- | :--------------------------- |
|id|int|游戏包自增ID|
|game_id|int|平台ID|
|game_name|string|游戏名称|
|version|string|版本|
|package_os|string|包类型 ios 或 android |
|package_md5|string|包md5值 |
|package_size|string|包大小 |
|log_time|string|上传时间 |
|remark|string|备注 |
|status|int|状态 1:关闭 2:开启 |

```
示例：

```

### 1.2.5 套餐拷贝

地址： /v1/pay/list/copy

方法： POST

接口参数传递方式：json

参数：

| 名称     | 类型   | 是否必填 | 备注                     |
| :------- | :----- | :------- | :----------------------- |
| id_list  | []int32  | 是       | 要复制的套餐自增ID列表|
| to\_game_id   | int32  | 是       | 复制的目标平台ID |

返回值：

| 名称   | 类型   | 备注                         |
| :----- | :----- | :--------------------------- |
| errno  | int32  | 0 表示无异常，其它表示有异常 |
| errmsg | string | 异常信息                     |

### 1.2.6 套餐添加

地址： /v1/pay/list/add

方法： POST

接口参数传递方式：json

参数：

| 名称     | 类型   | 是否必填 | 备注                     |
| :------- | :----- | :------- | :----------------------- |
| game_id  | int32  | 是       | 平台ID|
| game\_area_id  | int32  | 是       | 子游戏ID 默认-1 卡包:-2 子游戏套餐传对应子游戏ID|
| card_num  | int32  | 是       | 房卡数量|
| status  | int32  | 是       | 套餐状态 0:关闭 1:开启|
| price  | int32  | 是       | 套餐价格 单位:分|
| add_time  | string  | 是       | 套餐生效开始时间|
| end_time  | string  | 是       | 套餐生效结束时间|
| category  | string  | 是       | 套餐类型 user:普通用户套餐 agent:代理套餐 club:亲友圈套餐|
| present\_room_card  | int32  | 否       |套餐赠送房卡数量 |
| represent  | string  | 否      | 套餐描述|
| apple_key   | int32  | 否       | 苹果内购套餐key |
| is_white   | bool  | 否      | 是否是白名单套餐 true:是 默认:否|
| is_black   | bool  | 否       | 是否是黑名单套餐 true:是 默认:否|
|white_list  | []int | 否    | 白名单列表 |
|black_list  | []int | 否    | 黑名单列表 |
| card\_pack_id   | int32  | 否      | 卡包ID 只有卡包套餐才有该字段 |
| is\_senior_agent   | int32  | 否       | 是否是高级代理才能购买的套餐 0:否-默认 1:是 |
| purchase_limit_number   | int32  | 否      | 限购次数, 0代表不限制 |

返回值：

| 名称   | 类型   | 备注                         |
| :----- | :----- | :--------------------------- |
| errno  | int32  | 0 表示无异常，其它表示有异常 |
| errmsg | string | 异常信息                     |


### 1.2.7 套餐修改

地址： /v1/pay/list/update

方法： POST

接口参数传递方式：json

参数：

| 名称     | 类型   | 是否必填 | 备注                     |
| :------- | :----- | :------- | :----------------------- |
| id  | int32         | 是|套餐ID|
| game_id  | int32  | 是       | 平台ID|
| game\_area_id  | int32  | 是       | 子游戏ID 默认-1 卡包:-2 子游戏套餐传对应子游戏ID|
| card_num  | int32  | 是       | 房卡数量|
| status  | int32  | 是       | 套餐状态 0:关闭 1:开启|
| price  | int32  | 是       | 套餐价格 单位:分|
| add_time  | string  | 是       | 套餐生效开始时间 格式:2006-01-02 15:04:05|
| category  | string  | 是       | 套餐类型 user:普通用户套餐 agent:代理套餐 club:亲友圈套餐|
| end_time  | string  | 是       | 套餐生效结束时间 格式:2006-01-02 15:04:05|
| present\_room_card  | int32  | 否       |套餐赠送房卡数量 |
| represent  | string  | 否      | 套餐描述|
| apple_key   | int32  | 否       | 苹果内购套餐key |
| is_white   | bool  | 否      | 是否是白名单套餐 true:是 默认:否|
| is_black   | bool  | 否       | 是否是黑名单套餐 true:是 默认:否|
|white_list  | []int | 否    | 白名单列表 |
|black_list  | []int | 否    | 黑名单列表 |
| card\_pack_id   | int32  | 否      | 卡包ID 只有卡包套餐才有该字段 |
| is\_senior_agent   | int32  | 否       | 是否是高级代理才能购买的套餐 0:否-默认 1:是 |
| purchase_limit_number   | int32  | 否      | 限购次数, 0代表不限制 |

返回值：

| 名称   | 类型   | 备注                         |
| :----- | :----- | :--------------------------- |
| errno  | int32  | 0 表示无异常，其它表示有异常 |
| errmsg | string | 异常信息                     |

### 1.2.8 套餐删除

地址： /v1/pay/list/delete

方法： POST

接口参数传递方式：json

参数：

| 名称     | 类型   | 是否必填 | 备注                     |
| :------- | :----- | :------- | :----------------------- |
| id  | int32  | 是       | 套餐ID|

返回值：

| 名称   | 类型   | 备注                         |
| :----- | :----- | :--------------------------- |
| errno  | int32  | 0 表示无异常，其它表示有异常 |
| errmsg | string | 异常信息                     |

### 1.2.8 套餐列表获取

地址： /v1/pay/list/list

方法： POST

接口参数传递方式：json

参数：

| 名称     | 类型   | 是否必填 | 备注                     |
| :------- | :----- | :------- | :----------------------- |
| game_id  | int32  | 是       | 平台ID|
| game\_area_id   | int32  | 是       | 套餐子游戏ID |
| source_type   | int32  | 是       | 套餐类型 user:普通用户套餐 agent:代理套餐 club:亲友圈套餐 不传值代表所有类型|
| is_apple_product   | bool  |否| 是否是苹果内购套餐       | true:是 false:否 |

返回值：

| 名称   | 类型   | 备注                         |
| :----- | :----- | :--------------------------- |
| errno  | int32  | 0 表示无异常，其它表示有异常 |
| errmsg | string | 异常信息                     |
| data | []RoomCardProductList | 套餐列表                    

结构体RoomCardProductList

| 名称   | 类型   | 备注                         |
| :----- | :----- | :--------------------------- |
| id  | int32         | 套餐IDID|
| game_id  | int32         | 平台ID|
| game\_area_id  | int32         | 子游戏ID 默认-1 卡包:-2 子游戏套餐传对应子游戏ID|
| card_num  | int32        | 房卡数量|
| status  | int32      | 套餐状态 0:关闭 1:开启|
| price  | int32       | 套餐价格 单位:分|
| add_time  | string         | 套餐生效开始时间 格式:2006-01-02 15:04:05|
| category  | string         | 套餐类型 user:普通用户套餐 agent:代理套餐 club:亲友圈套餐|
| end_time  | string         | 套餐生效结束时间 格式:2006-01-02 15:04:05|
| present\_room_card  | int32       |套餐赠送房卡数量 |
| represent  | string       | 套餐描述|
| apple_key   | int32      | 苹果内购套餐app_key |
| is_white   | bool       | 是否是白名单套餐 true:是 默认:否|
| is_black   | bool        | 是否是黑名单套餐 true:是 默认:否|
|white_list  | []int   | 白名单列表 |
|black_list  | []int   | 黑名单列表 |
| card\_pack_id   | int32       | 卡包ID 只有卡包套餐才有该字段 |
| is\_senior_agent   | int32      | 是否是高级代理才能购买的套餐 0:否-默认 1:是 |
| purchase_limit_number   | int32  |  限购次数, 0代表不限制 |






# 2. 接口地址

## 2.1. 测试环境

http://123.206.215.185:9901

## 2.2. 正式环境（老环境）

外网：https://configapi2.xq5.com
内网：http://10.105.53.248:9901

## 2.3. 正式环境（新环境）

外网：https://configapi2.213451.com
内网：http://10.99.3.4:9901

## 2.4. 省外环境 (华为)
外网：https://configapi2.xq668.com  (http://configapi2.xq668.com)
