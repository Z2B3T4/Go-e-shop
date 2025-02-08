module project1

go 1.23.2

require (
	google.golang.org/grpc v1.67.1
	project1/User v1.23.2
)

require (
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/gorm v1.25.12 // indirect
	project1/Auth v1.23.2 // indirect
	project1/Common v1.23.2 // indirect
)

replace project1/User => ./User

replace (
	project1/Auth => ./Auth
	project1/Common => ./Common
)
