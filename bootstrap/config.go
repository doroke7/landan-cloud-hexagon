package bootstrap

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

//nolint:stylecheck,revive
type Config struct {
	HTTP struct {
		PORT string `mapstructure:"port"`
	} `mapstructure:"http"`
	WEBSOCKET struct {
		PORT string `mapstructure:"port"`
	} `mapstructure:"websocket"`
	SERVICES struct {
		HTTP struct {
			HOST  string `mapstructure:"host"`
			PORT  string `mapstructure:"port"`
			ADMIN struct {
				PUBLIC_KEY     string `mapstructure:"public_key"`
				PRIVATE_KEY    string `mapstructure:"private_key"`
				SIGNATURE      bool   `mapstructure:"signature"`
				AUTHENTICATION bool   `mapstructure:"authentication"`
				SALT           string `mapstructure:"salt"`
				JWT            struct {
					SECRET string `mapstructure:"secret"`
					KEY    string `mapstructure:"key"`
					IV     string `mapstructure:"iv"`
				} `mapstructure:"jwt"`
			} `mapstructure:"admin"`
			APP struct {
				PUBLIC_KEY     string `mapstructure:"public_key"`
				PRIVATE_KEY    string `mapstructure:"private_key"`
				SIGNATURE      bool   `mapstructure:"signature"`
				AUTHENTICATION bool   `mapstructure:"authentication"`
				SALT           string `mapstructure:"salt"`
				JWT            struct {
					SECRET string `mapstructure:"secret"`
					KEY    string `mapstructure:"key"`
					IV     string `mapstructure:"iv"`
				} `mapstructure:"jwt"`
			} `mapstructure:"app"`
		} `mapstructure:"http"`
		FACADE struct {
			HOST  string `mapstructure:"host"`
			PORT  string `mapstructure:"port"`
			ADMIN struct {
				PUBLIC_KEY  string `mapstructure:"public_key"`
				PRIVATE_KEY string `mapstructure:"private_key"`
				SIGNATURE   bool   `mapstructure:"signature"`
				SALT        string `mapstructure:"salt"`
				JWT         struct {
					SECRET string `mapstructure:"secret"`
					KEY    string `mapstructure:"key"`
					IV     string `mapstructure:"iv"`
				} `mapstructure:"jwt"`
			} `mapstructure:"admin"`
			APP struct {
				PUBLIC_KEY     string `mapstructure:"public_key"`
				PRIVATE_KEY    string `mapstructure:"private_key"`
				SIGNATURE      bool   `mapstructure:"signature"`
				AUTHENTICATION bool   `mapstructure:"authentication"`
				SALT           string `mapstructure:"salt"`
				JWT            struct {
					SECRET string `mapstructure:"secret"`
					KEY    string `mapstructure:"key"`
					IV     string `mapstructure:"iv"`
				} `mapstructure:"jwt"`
			} `mapstructure:"app"`
		} `mapstructure:"facade"`
		RESOURCE struct {
			HOST     string `mapstructure:"host"`
			PORT     string `mapstructure:"port"`
			USER     string `mapstructure:"user"`
			PASSWORD string `mapstructure:"password"`
		} `mapstructure:"resource"`
		WEBSOCKET struct {
			HOST  string `mapstructure:"host"`
			PORT  string `mapstructure:"port"`
			ADMIN struct {
				PUBLIC_KEY     string `mapstructure:"public_key"`
				PRIVATE_KEY    string `mapstructure:"private_key"`
				SIGNATURE      bool   `mapstructure:"signature"`
				AUTHENTICATION bool   `mapstructure:"authentication"`
				SALT           string `mapstructure:"salt"`
				JWT            struct {
					SECRET string `mapstructure:"secret"`
					KEY    string `mapstructure:"key"`
					IV     string `mapstructure:"iv"`
				} `mapstructure:"jwt"`
			} `mapstructure:"admin"`
			APP struct {
				PUBLIC_KEY     string `mapstructure:"public_key"`
				PRIVATE_KEY    string `mapstructure:"private_key"`
				SIGNATURE      bool   `mapstructure:"signature"`
				AUTHENTICATION bool   `mapstructure:"authentication"`
				SALT           string `mapstructure:"salt"`
				JWT            struct {
					SECRET string `mapstructure:"secret"`
					KEY    string `mapstructure:"key"`
					IV     string `mapstructure:"iv"`
				} `mapstructure:"jwt"`
			} `mapstructure:"app"`
		} `mapstructure:"websocket"`
		SOURCE struct {
			HOST string `mapstructure:"host"`
			PORT string `mapstructure:"port"`
		} `mapstructure:"source"`
		TCP struct {
			HOST string `mapstructure:"host"`
			PORT string `mapstructure:"port"`
		} `mapstructure:"tcp"`
		UDP struct {
			HOST string `mapstructure:"host"`
			PORT string `mapstructure:"port"`
		} `mapstructure:"udp"`
	} `mapstructure:"services"`
	CLIENTS struct {
		FACADE struct {
			HOSTS   []string `mapstructure:"hosts"`
			PORTS   []string `mapstructure:"ports"`
			TIMEOUT int      `mapstructure:"timeout"`
		} `mapstructure:"facade"`
		RESOURCE struct {
			HOSTS    []string `mapstructure:"hosts"`
			PORTS    []string `mapstructure:"ports"`
			USER     string   `mapstructure:"user"`
			PASSWORD string   `mapstructure:"password"`
			TIMEOUT  int      `mapstructure:"timeout"`
		} `mapstructure:"resource"`
		SOURCE struct {
			HOSTS   []string `mapstructure:"hosts"`
			PORTS   []string `mapstructure:"ports"`
			WEIGHTS []uint   `mapstructure:"weights"`
			TIMEOUT int      `mapstructure:"timeout"`
		} `mapstructure:"source"`
		TCP struct {
			HOSTS   []string `mapstructure:"hosts"`
			PORTS   []string `mapstructure:"ports"`
			TIMEOUT int      `mapstructure:"timeout"`
			POOL    int      `mapstructure:"pool"`
		} `mapstructure:"tcp"`
	} `mapstructure:"clients"`
	DATABASE struct {
		USER                 string `mapstructure:"user"`
		PASSWORD             string `mapstructure:"password"`
		PREFIX               string `mapstructure:"prefix"`
		CHARSET              string `mapstructure:"charset"`
		NAME                 string `mapstructure:"name"`
		MAX_IDLE_CONNECTIONS int    `mapstructure:"max_idle_connections"`
		TIMEOUT              int    `mapstructure:"timeout"`
		READ                 struct {
			HOSTS []string `mapstructure:"hosts"`
			PORTS []string `mapstructure:"ports"`
		} `mapstructure:"read"`
		WRITE struct {
			HOSTS []string `mapstructure:"hosts"`
			PORTS []string `mapstructure:"ports"`
		} `mapstructure:"write"`
	} `mapstructure:"database"`
	DB struct {
		HOST string `mapstructure:"host"` // 映射键名：它告诉解码器，配置文件（或 Map）里的键名如果是 "host"，就对应填入结构体的 HOST 字段
		USER string `mapstructure:"user"`
		PASS string `mapstructure:"pass"` //
	} `mapstructure:"db"`
	MONGODB struct {
		PROTOCOL      string `mapstructure:"protocol"`
		HOST          string `mapstructure:"host"`
		PORT          string `mapstructure:"port"`
		NAME          string `mapstructure:"name"`
		USER          string `mapstructure:"user"`
		PASSWORD      string `mapstructure:"password"`
		AUTH_SOURCE   string `mapstructure:"auth_source"`
		MAX_POOL_SIZE uint64 `mapstructure:"max_pool_size"`
		MIN_POOL_SIZE uint64 `mapstructure:"min_pool_size"`
	} `mapstructure:"mongodb"`
	POSTGRESQL struct {
		HOST                 string `mapstructure:"host"`
		PORT                 string `mapstructure:"port"`
		USER                 string `mapstructure:"user"`
		PASSWORD             string `mapstructure:"password"`
		NAME                 string `mapstructure:"name"`
		SSLMODE              string `mapstructure:"sslmode"`
		MAX_IDLE_CONNECTIONS int    `mapstructure:"max_idle_connections"`
		TIMEOUT              int    `mapstructure:"timeout"`
	} `mapstructure:"postgresql"`
	SQLITE struct {
		PATH                 string `mapstructure:"path"`
		MAX_IDLE_CONNECTIONS int    `mapstructure:"max_idle_connections"`
		TIMEOUT              int    `mapstructure:"timeout"`
	} `mapstructure:"sqlite"`
	TIDB struct {
		USER                 string   `mapstructure:"user"`
		PASSWORD             string   `mapstructure:"password"`
		PREFIX               string   `mapstructure:"prefix"`
		CHARSET              string   `mapstructure:"charset"`
		NAME                 string   `mapstructure:"name"`
		MAX_IDLE_CONNECTIONS int      `mapstructure:"max_idle_connections"`
		TIMEOUT              int      `mapstructure:"timeout"`
		HOSTS                []string `mapstructure:"hosts"`
		PORTS                []string `mapstructure:"ports"`
	} `mapstructure:"tidb"`
	GAUSSDB struct {
		HOST                 string `mapstructure:"host"`
		PORT                 string `mapstructure:"port"`
		USER                 string `mapstructure:"user"`
		PASSWORD             string `mapstructure:"password"`
		NAME                 string `mapstructure:"name"`
		SSLMODE              string `mapstructure:"sslmode"`
		MAX_IDLE_CONNECTIONS int    `mapstructure:"max_idle_connections"`
		TIMEOUT              int    `mapstructure:"timeout"`
	} `mapstructure:"gaussdb"`
	ORACLE struct {
		HOST                 string `mapstructure:"host"`
		PORT                 string `mapstructure:"port"`
		USER                 string `mapstructure:"user"`
		PASSWORD             string `mapstructure:"password"`
		SERVICE              string `mapstructure:"service"`
		MAX_IDLE_CONNECTIONS int    `mapstructure:"max_idle_connections"`
		TIMEOUT              int    `mapstructure:"timeout"`
	} `mapstructure:"oracle"`
	CLICKHOUSE struct {
		HOST                 string `mapstructure:"host"`
		PORT                 string `mapstructure:"port"`
		USER                 string `mapstructure:"user"`
		PASSWORD             string `mapstructure:"password"`
		NAME                 string `mapstructure:"name"`
		DIAL_TIMEOUT         string `mapstructure:"dial_timeout"`
		READ_TIMEOUT         string `mapstructure:"read_timeout"`
		MAX_IDLE_CONNECTIONS int    `mapstructure:"max_idle_connections"`
		TIMEOUT              int    `mapstructure:"timeout"`
	} `mapstructure:"clickhouse"`
	REDIS struct {
		HOST     string   `mapstructure:"host"`
		HOSTS    []string `mapstructure:"hosts"`
		USERNAME string   `mapstructure:"username"`
		PASSWORD string   `mapstructure:"password"` //
		PORT     string   `mapstructure:"port"`
		PORTS    []string `mapstructure:"ports"`
		DB       int      `mapstructure:"db"`
		CLUSTER  bool     `mapstructure:"cluster"`
		TIMEOUT  int      `mapstructure:"timeout"`
	} `mapstructure:"redis"`
	AMQP struct {
		HOST string `mapstructure:"host"`
		PORT string `mapstructure:"port"`
		USER string `mapstructure:"user"`
		PASS string `mapstructure:"pass"` //
	} `mapstructure:"amqp"`
	KAFKA struct {
		BROKERS []string `mapstructure:"brokers"`
		USER    string   `mapstructure:"user"`
		PASS    string   `mapstructure:"pass"`
		VERSION string   `mapstructure:"version"`
	} `mapstructure:"kafka"`
	BEANSTALK struct {
		HOST    string `mapstructure:"host"`
		PORT    string `mapstructure:"port"`
		TIMEOUT int    `mapstructure:"timeout"`
	} `mapstructure:"beanstalk"`
	NATS struct {
		HOST    string `mapstructure:"host"`
		PORT    string `mapstructure:"port"`
		USER    string `mapstructure:"user"`
		PASS    string `mapstructure:"pass"`
		TIMEOUT int    `mapstructure:"timeout"`
	} `mapstructure:"nats"`
	ACTIVEMQ struct {
		HOST    string `mapstructure:"host"`
		PORT    string `mapstructure:"port"`
		USER    string `mapstructure:"user"`
		PASS    string `mapstructure:"pass"`
		TIMEOUT int    `mapstructure:"timeout"`
	} `mapstructure:"activemq"`
	ZEROMQ struct {
		HOST    string `mapstructure:"host"`
		PORT    string `mapstructure:"port"`
		TIMEOUT int    `mapstructure:"timeout"`
	} `mapstructure:"zeromq"`
	PULSAR struct {
		HOST    string `mapstructure:"host"`
		PORT    string `mapstructure:"port"`
		TIMEOUT int    `mapstructure:"timeout"`
	} `mapstructure:"pulsar"`
	REDPANDA struct {
		BROKERS []string `mapstructure:"brokers"`
		USER    string   `mapstructure:"user"`
		PASS    string   `mapstructure:"pass"`
		VERSION string   `mapstructure:"version"`
	} `mapstructure:"redpanda"`
	DEFAULT struct {
		DEBUG bool `mapstructure:"debug"`
	} `mapstructure:"default"`
	PARTITIONS map[string]string `mapstructure:"partitions"`
	LOGGERS    struct {
		DEFAULT struct {
			DIRECTORY   string `mapstructure:"directory"`
			MAX_SIZE    int    `mapstructure:"max_size"`
			MAX_BACKUPS int    `mapstructure:"max_backups"`
			MAX_AGE     int    `mapstructure:"max_age"`

			COMPRESS bool `mapstructure:"compress"`
			STATUS   bool `mapstructure:"status"`
		} `mapstructure:"default"`
		MIDDLEWARE struct {
			DIRECTORY   string `mapstructure:"directory"`
			MAX_SIZE    int    `mapstructure:"max_size"`
			MAX_BACKUPS int    `mapstructure:"max_backups"`
			MAX_AGE     int    `mapstructure:"max_age"`

			COMPRESS bool `mapstructure:"compress"`
			STATUS   bool `mapstructure:"status"`
		} `mapstructure:"middleware"`
		CONTROLLER struct {
			DIRECTORY   string `mapstructure:"directory"`
			MAX_SIZE    int    `mapstructure:"max_size"`
			MAX_BACKUPS int    `mapstructure:"max_backups"`
			MAX_AGE     int    `mapstructure:"max_age"`
			COMPRESS    bool   `mapstructure:"compress"`
			STATUS      bool   `mapstructure:"status"`
		} `mapstructure:"controller"`
		INTERCEPTOR struct {
			DIRECTORY   string `mapstructure:"directory"`
			MAX_SIZE    int    `mapstructure:"max_size"`
			MAX_BACKUPS int    `mapstructure:"max_backups"`
			MAX_AGE     int    `mapstructure:"max_age"`
			COMPRESS    bool   `mapstructure:"compress"`
			STATUS      bool   `mapstructure:"status"`
		} `mapstructure:"interceptor"`
		SDK struct {
			DIRECTORY   string `mapstructure:"directory"`
			MAX_SIZE    int    `mapstructure:"max_size"`
			MAX_BACKUPS int    `mapstructure:"max_backups"`
			MAX_AGE     int    `mapstructure:"max_age"`
			COMPRESS    bool   `mapstructure:"compress"`
			STATUS      bool   `mapstructure:"status"`
		} `mapstructure:"sdk"`
		SERVICE struct {
			DIRECTORY   string `mapstructure:"directory"`
			MAX_SIZE    int    `mapstructure:"max_size"`
			MAX_BACKUPS int    `mapstructure:"max_backups"`
			MAX_AGE     int    `mapstructure:"max_age"`
			COMPRESS    bool   `mapstructure:"compress"`
			STATUS      bool   `mapstructure:"status"`
		} `mapstructure:"service"`
	} `mapstructure:"loggers"`
	TABLE struct {
		ADMIN_USER struct {
			PASSWORD string `mapstructure:"password"`
		} `mapstructure:"admin_user"`
	} `mapstructure:"table"`
}

var CONFIG Config

func init() {
	// 加载 .env 到系统环境变量（文件不存在时不报错）
	_ = godotenv.Load()

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()                                   // 读取环境变量
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 嵌套字段用 _ 连接

	// 自动读取 ./config/ 目录下所有 yaml 文件，文件名作为顶层命名空间
	// 例: db.yaml 内的 host → db.host
	aFiles, oErr := filepath.Glob("./config/*.yaml")
	if oErr != nil {
		log.Fatalf("failed to glob config dir: %v", oErr)
	}
	for _, sFile := range aFiles {
		sName := strings.TrimSuffix(filepath.Base(sFile), filepath.Ext(sFile))

		oSub := viper.New()
		oSub.SetConfigFile(sFile)
		if oErr := oSub.ReadInConfig(); oErr != nil {
			log.Fatalf("failed to read config %s: %v", sFile, oErr)
		}
		// 以文件名包一层后合并到主 viper
		viper.MergeConfigMap(map[string]interface{}{
			sName: oSub.AllSettings(),
		})
	}

	// 自动绑定所有 key，让环境变量覆盖嵌套字段生效
	for _, sKey := range viper.AllKeys() {
		_ = viper.BindEnv(sKey)
	}

	// 简单来说，它的作用是：
	// 把 Viper 内部缓存的所有配置数据（Map 格式），一次性“灌入”到你定义的 Go 结构体（Struct）中。
	if err := viper.Unmarshal(&CONFIG); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
}
