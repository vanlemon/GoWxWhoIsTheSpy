module lmf.mortal.com/GoWxWhoIsTheSpy

go 1.15

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/jinzhu/gorm v1.9.16
	github.com/scylladb/go-set v1.0.2
	lmf.mortal.com/GoLimiter v0.0.0-00010101000000-000000000000
	lmf.mortal.com/GoLogs v0.0.0-00010101000000-000000000000
)

replace (
	lmf.mortal.com/GoLimiter => ../GoLimiter
	lmf.mortal.com/GoLogs => ../GoLogs
)
