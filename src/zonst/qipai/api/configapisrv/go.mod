module zonst/qipai/api/configapisrv

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/DHowett/go-plist v0.0.0-20201203080718-1454fab16a06 // indirect
	github.com/andrianbdn/iospng v0.0.0-20180730113000-dccef1992541 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/color v1.12.0 // indirect
	github.com/fwhezfwhez/cmap v1.2.1
	github.com/fwhezfwhez/errorx v1.1.0
	github.com/garyburd/redigo v1.6.2
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.2
	github.com/go-xweb/log v0.0.0-20140701090824-270d183ad77e
	github.com/helloteemo/ashe v0.0.1
	github.com/jinzhu/gorm v1.9.16
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.2
	github.com/phinexdaz/ipapk v0.0.0-20180706142810-91f3861dffcf
	github.com/shogo82148/androidbinary v1.0.2
	github.com/stretchr/testify v1.6.1
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac
	howett.net/plist v0.0.0-20201203080718-1454fab16a06
	zonst/logging v0.0.0
	zonst/qipai/gin-middlewares/secretauth v0.0.0
	zonst/qipai/gin-middlewares/tokenauth v0.0.0
	zonst/qipai/libcos v0.0.0
)

replace (
	github.com/DHowett/go-plist v0.0.0-20201203080718-1454fab16a06 => howett.net/plist v0.0.0-20201203080718-1454fab16a06
	howett.net/plist v0.0.0-20201203080718-1454fab16a06 => github.com/DHowett/go-plist v0.0.0-20201203080718-1454fab16a06
	zonst/logging => ../../../logging
	zonst/qipai/gin-middlewares/secretauth => ../../gin-middlewares/secretauth
	zonst/qipai/gin-middlewares/tokenauth => ../../../../zonst/qipai/gin-middlewares/tokenauth
	zonst/qipai/libcos => ../../libcos
)
