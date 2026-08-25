package model

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/env"
	"github.com/modelbus/one-api-pro/common/helper"
	appLogger "github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/common/random"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)
var DB *gorm.DB
var LOG_DB *gorm.DB

func CreateRootAccountIfNeed() error {
	var user User
	//if user.Status != util.UserStatusEnabled {
	if err := DB.First(&user).Error; err != nil {
		appLogger.SysLog("no user exists, creating a root user for you: username is root, password is 123456")
		hashedPassword, err := common.Password2Hash("123456")
		if err != nil {
			return err
		}
		accessToken := random.GetUUID()
		if config.InitialRootAccessToken != "" {
			accessToken = config.InitialRootAccessToken
		}
		rootUser := User{
			Username:    "root",
			Password:    hashedPassword,
			Role:        RoleRootUser,
			Status:      UserStatusEnabled,
			DisplayName: "Root User",
			AccessToken: accessToken,
			Quota:       500000000000000,
		}
		DB.Create(&rootUser)
		if config.InitialRootToken != "" {
			appLogger.SysLog("creating initial root token as requested")
			token := Token{
				Id:             1,
				UserId:         rootUser.Id,
				Key:            config.InitialRootToken,
				Status:         TokenStatusEnabled,
				Name:           "Initial Root Token",
				CreatedTime:    helper.GetTimestamp(),
				AccessedTime:   helper.GetTimestamp(),
				ExpiredTime:    -1,
				RemainQuota:    500000000000000,
				UnlimitedQuota: true,
			}
			DB.Create(&token)
		}
	}
	return nil
}

func chooseDB(envName string) (*gorm.DB, error) {
	dsn := os.Getenv(envName)

	switch {
	case strings.HasPrefix(dsn, "postgres://"):
		// Use PostgreSQL
		return openPostgreSQL(dsn)
	case dsn != "":
		// Use MySQL
		return openMySQL(dsn)
	default:
		// Use SQLite
		return openSQLite()
	}
}

func openPostgreSQL(dsn string) (*gorm.DB, error) {
	appLogger.SysLog("using PostgreSQL as database")
	common.UsingPostgreSQL = true
	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		PrepareStmt: true, // precompile SQL
		Logger: newGormLogger(),
	})
}

func openMySQL(dsn string) (*gorm.DB, error) {
	appLogger.SysLog("using MySQL as database")
	common.UsingMySQL = true
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true, // precompile SQL
		Logger: newGormLogger(),
	})
}

// gormLogWriter 包装 os.Stdout 实现 gormlogger.Writer 接口
type gormLogWriter struct{}

// Printf 实现 gormlogger.Writer 接口的 Printf 方法
func (gormLogWriter) Printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
}

// newGormLogger 创建一个自定义的 GORM logger
// 关键设置：IgnoreRecordNotFoundError = true
// 避免在 First 查询找不到记录时打印误导性的红色 ERROR 日志
func newGormLogger() gormlogger.Interface {
	return gormlogger.New(
		gormLogWriter{},
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormlogger.Warn, // 只记录警告和错误，不记录 info
			IgnoreRecordNotFoundError: true,            // 忽略 record not found
			Colorful:                  false,
		},
	)
}

func openSQLite() (*gorm.DB, error) {
	appLogger.SysLog("SQL_DSN not set, using SQLite as database")
	common.UsingSQLite = true
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(%d)", common.SQLitePath, common.SQLiteBusyTimeout)
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		PrepareStmt: true, // precompile SQL
		Logger: newGormLogger(),
	})
}

func InitDB() {
	var err error
	DB, err = chooseDB("SQL_DSN")
	if err != nil {
		appLogger.FatalLog("failed to initialize database: " + err.Error())
		return
	}

	sqlDB := setDBConns(DB)

	if !config.IsMasterNode {
		return
	}

	if common.UsingMySQL {
		_, _ = sqlDB.Exec("DROP INDEX idx_channels_key ON channels;") // TODO: delete this line when most users have upgraded
	}

	appLogger.SysLog("database migration started")
	if err = migrateDB(); err != nil {
		appLogger.FatalLog("failed to migrate database: " + err.Error())
		return
	}
	appLogger.SysLog("database migrated")
	if err = InitDefaultPrices(); err != nil {
		appLogger.SysError("failed to initialize default prices: " + err.Error())
	}
	InitModelPriceCache()
	InitGroupPriceCache()
}

func migrateDB() error {
	var err error
	if err = DB.AutoMigrate(&Channel{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&Token{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&User{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&Option{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&Redemption{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&Ability{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&Log{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&Channel{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&Plan{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&UserPlan{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&Order{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&SystemSetting{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&PlanUsage{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&ChannelCounter{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&ModelPrice{}); err != nil {
		return err
	}
	if err = DB.AutoMigrate(&GroupPrice{}); err != nil {
		return err
	}
	// Add order_id column to existing user_plans on legacy databases.
	if common.UsingMySQL {
		var colCount int
		DB.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'user_plans' AND COLUMN_NAME = 'order_id'").Scan(&colCount)
		if colCount == 0 {
			DB.Exec("ALTER TABLE user_plans ADD COLUMN order_id int(11) NOT NULL DEFAULT 0 AFTER plan_id")
		}
	} else if common.UsingPostgreSQL {
		var colCount int
		DB.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_name = 'user_plans' AND column_name = 'order_id'").Scan(&colCount)
		if colCount == 0 {
			DB.Exec("ALTER TABLE user_plans ADD COLUMN order_id integer NOT NULL DEFAULT 0")
		}
	}
	// Migrate old tokens column to prompt_tokens if it exists
	if common.UsingMySQL {
		var colCount int
		DB.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'plan_usages' AND COLUMN_NAME = 'tokens'").Scan(&colCount)
		if colCount > 0 {
			DB.Exec("ALTER TABLE plan_usages CHANGE COLUMN tokens prompt_tokens bigint NOT NULL DEFAULT 0")
		}
	} else if common.UsingPostgreSQL {
		var colCount int
		DB.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_name = 'plan_usages' AND column_name = 'tokens'").Scan(&colCount)
		if colCount > 0 {
			DB.Exec("ALTER TABLE plan_usages RENAME COLUMN tokens TO prompt_tokens")
		}
	}
	return nil
}

func InitLogDB() {
	if os.Getenv("LOG_SQL_DSN") == "" {
		LOG_DB = DB
		return
	}

	appLogger.SysLog("using secondary database for table logs")
	var err error
	LOG_DB, err = chooseDB("LOG_SQL_DSN")
	if err != nil {
		appLogger.FatalLog("failed to initialize secondary database: " + err.Error())
		return
	}

	setDBConns(LOG_DB)

	if !config.IsMasterNode {
		return
	}

	appLogger.SysLog("secondary database migration started")
	err = migrateLOGDB()
	if err != nil {
		appLogger.FatalLog("failed to migrate secondary database: " + err.Error())
		return
	}
	appLogger.SysLog("secondary database migrated")
}

func migrateLOGDB() error {
	var err error
	if err = LOG_DB.AutoMigrate(&Log{}); err != nil {
		return err
	}
	return nil
}

func setDBConns(db *gorm.DB) *sql.DB {
	if config.DebugSQLEnabled {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		appLogger.FatalLog("failed to connect database: " + err.Error())
		return nil
	}

	sqlDB.SetMaxIdleConns(env.Int("SQL_MAX_IDLE_CONNS", 100))
	sqlDB.SetMaxOpenConns(env.Int("SQL_MAX_OPEN_CONNS", 1000))
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(env.Int("SQL_MAX_LIFETIME", 60)))
	return sqlDB
}

func closeDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	return err
}

func CloseDB() error {
	if LOG_DB != DB {
		err := closeDB(LOG_DB)
		if err != nil {
			return err
		}
	}
	return closeDB(DB)
}
