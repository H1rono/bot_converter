package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"git.trap.jp/toki/bot_converter/migrate"
	"git.trap.jp/toki/bot_converter/router"
	"git.trap.jp/toki/bot_converter/service"
)

// Config describes server config.
type Config struct {
	// Port number to listen on.
	Port int `mapstructure:"port" yaml:"port"`

	// Origin is the origin URL of the bot. e.g. http://q.trap.jp
	Origin string `mapstructure:"origin" yaml:"origin"`

	// Traq describes traq bot settings.
	Traq struct {
		// VerificationToken is verification token for verifying http bot events.
		VerificationToken string `mapstructure:"verificationToken" yaml:"verificationToken"`
		// AccessToken is access token for accessing traq API.
		AccessToken string `mapstructure:"accessToken" yaml:"accessToken"`
		// UserID is the user UUID of the bot.
		UserID string `mapstructure:"userID" yaml:"userID"`
		// Prefix is the bot command prefix.
		Prefix string `mapstructure:"prefix" yaml:"prefix"`
	} `mapstructure:"traq" yaml:"traq"`

	// MariaDB describes db settings.
	MariaDB struct {
		// Port is MariaDB port.
		Port int `mapstructure:"port" yaml:"port"`
		// Hostname is MariaDB host.
		Hostname string `mapstructure:"hostname" yaml:"hostname"`
		// Username is MariaDB user.
		Username string `mapstructure:"username" yaml:"username"`
		// Password is password for the above user.
		Password string `mapstructure:"password" yaml:"password"`
		// Database is database name.
		Database string `mapstructure:"database" yaml:"database"`
	} `mapstructure:"mariadb" yaml:"mariadb"`
}

var c Config

func init() {
	viper.SetDefault("port", 3000)
	viper.SetDefault("origin", "")
	viper.SetDefault("traq.verificationToken", "")
	viper.SetDefault("traq.accessToken", "")
	viper.SetDefault("traq.userID", uuid.Nil)
	viper.SetDefault("traq.prefix", "/")
	viper.SetDefault("mariadb.port", 3306)
	viper.SetDefault("mariadb.hostname", "localhost")
	viper.SetDefault("mariadb.username", "root")
	viper.SetDefault("mariadb.password", "password")
	viper.SetDefault("mariadb.database", "poll")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("an error occurred while unmarshalling config: %s", err)
	}
}

// initDB initializes DB connection and executes migration.
func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true&loc=Local",
			c.MariaDB.Username,
			c.MariaDB.Password,
			c.MariaDB.Hostname,
			c.MariaDB.Port,
			c.MariaDB.Database,
		),
		DefaultStringSize:       256,
		DontSupportRenameIndex:  true,
		DontSupportRenameColumn: true,
	}))
	if err != nil {
		return nil, err
	}

	if err := migrate.Migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

// provideRouterConfig provides router.Config.
func provideRouterConfig() router.Config {
	return router.Config{
		AccessToken: c.Traq.AccessToken,
	}
}

// provideBotConfig provides service.Config.
func provideBotConfig() service.Config {
	return service.Config{
		VerificationToken: c.Traq.VerificationToken,
		AccessToken:       c.Traq.AccessToken,
		BotID:             uuid.Must(uuid.FromString(c.Traq.UserID)),
		Prefix:            c.Traq.Prefix,
		Origin:            c.Origin,
	}
}
