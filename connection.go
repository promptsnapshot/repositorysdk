package repositorysdk

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// PostgresDatabaseConfig is a struct that holds the configuration details required to establish a connection
// with a PostgreSQL database.
type PostgresDatabaseConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Name        string `mapstructure:"name"`
	SSL         string `mapstructure:"ssl"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`
	MaxOpenConn int    `mapstructure:"max_open_conn"`
}

// GetMaxIdleConn returns the maximum number of idle connections in the connection pool.
// If the value is not set, the default value of 10 is returned.
//
// Returns:
// - int: the maximum number of idle connections.
func (c *PostgresDatabaseConfig) GetMaxIdleConn() int {
	if c.MaxIdleConn == 0 {
		return 10
	}

	return c.MaxIdleConn
}

// GetMaxOpenConn returns the maximum number of open connections in the connection pool.
// If the value is not set, the default value of 10 is returned.
//
// Returns:
// - int: the maximum number of open connections.
func (c *PostgresDatabaseConfig) GetMaxOpenConn() int {
	if c.MaxOpenConn == 0 {
		return 10
	}

	return c.MaxOpenConn
}

// InitPostgresDatabase initializes a connection to a PostgreSQL database using the given configuration details.
//
// Parameters:
// - conf: a pointer to a PostgresDatabaseConfig struct containing the database configuration details.
// - isDebug: a boolean value to enable or disable the GORM logging mode.
//
// Returns:
// - *gorm.DB: a pointer to the GORM database object.
// - error: an error if something goes wrong, otherwise nil.
func InitPostgresDatabase(conf *PostgresDatabaseConfig, isDebug bool) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", conf.Host, conf.Port, conf.User, conf.Password, conf.Name, conf.SSL)

	gormConf := &gorm.Config{}

	if !isDebug {
		gormConf.Logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConf)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(conf.GetMaxIdleConn())
	sqlDB.SetMaxOpenConns(conf.GetMaxOpenConn())

	return db, nil
}

// RedisConfig is a struct that holds the configuration details required to establish a connection
// with a Redis database.
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// InitRedisConnect initializes a connection to a Redis database using the given configuration details.
//
// Parameters:
// - conf: a pointer to a RedisConfig struct containing the database configuration details.
//
// Returns:
// - *redis.Client: a pointer to the Redis client object.
// - error: an error if something goes wrong, otherwise nil.
func InitRedisConnect(conf *RedisConfig) (cache *redis.Client, err error) {
	cache = redis.NewClient(&redis.Options{
		Addr:     conf.Host,
		Password: conf.Password,
		DB:       conf.DB,
	})

	return
}
