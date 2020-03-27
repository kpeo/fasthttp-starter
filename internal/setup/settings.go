package setup

import (
	"database/sql"
	"github.com/fasthttp/router"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

type Settings struct {
	File      string
	Type      string
	Name      string
	Prefix    string
	Version   string
	BuildTime string
	Commit    string
	Config    *viper.Viper
	Logger    *zap.Logger
	Db        *sql.DB
	Router    *router.Router
}

func NewSettings(app *Settings) (*viper.Viper, error) {
	v := viper.New()
	v.SetEnvPrefix(app.Prefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if len(app.File) > 0 {
		v.SetConfigType(app.SafeType())
		v.SetConfigFile(app.File)
		err := v.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (a Settings) GetLogger() *zap.Logger {
	return a.Logger
}

func (a Settings) GetConfig() *viper.Viper {
	return a.Config
}

func (a Settings) GetDb() *sql.DB {
	return a.Db
}

func (a Settings) SafeType() string {
	switch t := a.Type; t {
	case "toml", "yml", "yaml":
	default:
		return "yml"
	}
	return a.Type
}
