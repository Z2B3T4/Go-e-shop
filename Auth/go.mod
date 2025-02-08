module Auth

go 1.23.2

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang-jwt/jwt/v4 v4.5.1
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.35.1
	project1/Common v1.23.2
)

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.18 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/nacos-group/nacos-sdk-go v1.1.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	gopkg.in/ini.v1 v1.42.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/gorm v1.25.12 // indirect
)

require (
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	project1/Auth v1.23.2
	project1/User v1.23.2
)

replace project1/User => ../User

replace project1/Common => ../Common

replace project1/Auth => ../Auth
