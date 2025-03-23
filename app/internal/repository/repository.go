package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"app/pkg/jwt"
	"app/pkg/log"
	"app/pkg/zapgorm2"
)

const ctxTxKey = "TxKey"

type Repository struct {
	db []*gorm.DB
	//rdb    *redis.Client
	logger *log.Logger
}

func NewRepository(
	logger *log.Logger,
	db []*gorm.DB,
	// rdb *redis.Client,
) *Repository {
	return &Repository{
		db: db,
		//rdb:    rdb,
		logger: logger,
	}
}

type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error, indexes ...int) error
}

func NewTransaction(r *Repository) Transaction {
	return r
}

// DB return tx
// If you need to create a Transaction, you must call DB(ctx) and Transaction(ctx,fn)
func (r *Repository) DB(ctx context.Context, indexes ...int) *gorm.DB {
	v := ctx.Value(ctxTxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	if len(indexes) > 0 {
		return r.db[1].WithContext(ctx)
	}
	return r.db[0].WithContext(ctx)
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error, indexes ...int) error {
	if len(indexes) > 0 {
		return r.db[1].WithContext(ctx).Transaction(
			func(tx *gorm.DB) error {
				ctx = context.WithValue(ctx, ctxTxKey, tx)
				return fn(ctx)
			},
		)
	}
	return r.db[0].WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			ctx = context.WithValue(ctx, ctxTxKey, tx)
			return fn(ctx)
		},
	)
}

func NewDB(conf *viper.Viper, l *log.Logger) []*gorm.DB {
	var (
		db       *gorm.DB
		db_admin *gorm.DB
		err      error
	)

	logger := zapgorm2.New(l.Logger)
	driver := conf.GetString("data.db.user.driver")
	dsn := conf.GetString("data.db.user.dsn")
	dsn_admin := conf.GetString("data.db.user.dsn_admin")

	// GORM doc: https://gorm.io/docs/connecting_to_the_database.html
	switch driver {
	case "mysql":
		db, err = gorm.Open(
			mysql.Open(dsn), &gorm.Config{
				Logger: logger,
			},
		)
		db_admin, err = gorm.Open(
			mysql.Open(dsn_admin), &gorm.Config{
				Logger: logger,
			},
		)
	case "postgres":
		db, err = gorm.Open(
			postgres.New(
				postgres.Config{
					DSN:                  dsn,
					PreferSimpleProtocol: true, // disables implicit prepared statement usage
				},
			), &gorm.Config{},
		)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		panic("unknown db driver")
	}
	if err != nil {
		panic(err)
	}
	db = db.Debug()
	db_admin = db_admin.Debug()

	// Connection Pool config
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	sqlDB, err = db_admin.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return []*gorm.DB{db, db_admin}
}

func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     conf.GetString("data.redis.addr"),
			Password: conf.GetString("data.redis.password"),
			DB:       conf.GetInt("data.redis.db"),
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	return rdb
}

func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return v.(*jwt.MyCustomClaims).UserId
}
