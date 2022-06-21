package initialize

import (
	"fmt"
	"github.com/atmoxao/sql-tools/pkg/global"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/logrusadapter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func init() {
	log.Info("config init")
	pflag.String("config-path", "./conf/", "config file path")
	pflag.String("config-name", "config-local", "config file name")
	pflag.String("config-type", "yaml", "config file type")

	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(fmt.Errorf("fatal  error BindPFlags: %w", err))
	}

	ConfigPath := viper.GetString("config-path")
	ConfigName := viper.GetString("config-name")
	ConfigType := viper.GetString("config-type")

	viper.SetConfigName(ConfigName) // name of config file (without extension)

	viper.SetConfigType(ConfigType) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(ConfigPath) // optionally look for config in the working directory
	viper.AutomaticEnv()
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&global.Conf)
	if err != nil {
		panic(fmt.Errorf("fatal error Unmarshal config : %w", err))
	}
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	log.Info("init logs")
}

func init() {
	log.Info("init mysql")
	err := initMySQL()
	if err != nil {
		log.Fatalln(err)
	}
}

func initMySQL() (err error) {
	dsn := "root:123456@tcp(vm202002:3306)/demo"
	global.Conn, err = sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		return
	}

	logger := log.New()
	logger.Level = log.DebugLevel // miminum level
	logger.SetReportCaller(true)
	logger.Formatter = &log.TextFormatter{} // logrus automatically add time field
	global.Conn.DB = sqldblogger.OpenDriver(
		dsn,
		global.Conn.Driver(),
		logrusadapter.New(logger),
		// optional config...
	)

	global.Conn.SetMaxOpenConns(200)
	global.Conn.SetMaxIdleConns(10)

	err = global.Conn.Ping()

	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		return
	}

	// migration

	err = migration()

	return err
}

func migration() error {

	// Run migrations
	driver, err := mysql.WithInstance(global.Conn.DB, &mysql.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "./databases/migrate"), // file://path/to/directory
		"mysql", driver)

	if err != nil {
		log.Errorf("migration failed... %v", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
		return err
	}

	log.Println("Database migrated")
	return nil
}
