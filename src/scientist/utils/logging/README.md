# README
支持日志等级的简单日志库

## 使用

~~~
import "zonst/logging"

logging.Debugf("OnGameSocket: status:%v\n", status)
logging.Infof("OnGameSocket: status:%v\n", status)
logging.Errorf("OnGameSocket: status:%v\n", status)


## TODO
* 支持UDP传输日志