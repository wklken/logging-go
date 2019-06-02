package hook

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// TODO: 1. settings here
// TODO: 2. pool here

// Add test:
// 1. https://github.com/rogierlommers/logrus-redis-hook/blob/master/logrus_redis.go
// 2. https://github.com/lazyjin/logrus-redis-cluster-hook/blob/master/logrus_redis.go

// HookConfig stores configuration needed to setup the hook
type RedisHookConfig struct {
	Host     string
	Port     int
	DB       int
	Key      string
	Password string
	PoolSize int

	App      string
	Hostname string

	LogFormat string
}

// RedisHook to sends logs to Redis server
type RedisLogHook struct {
	redisClient *redis.Client
	redisKey    string
	logFormat   string

	app      string
	hostname string
}

// redis: https://github.com/TykTechnologies/tyk/blob/master/redis_logrus_hook.go
func (r RedisLogHook) New(name string, settings map[string]string, formatter logrus.Formatter) (logrus.Hook, error) {
	if err := validateRequiredHookSettings(name, settings, []string{"host", "port", "db", "key"}); err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(settings["port"])
	if err != nil {
		return nil, errors.New("port should be integer")
	}
	db, err := strconv.Atoi(settings["db"])
	if err != nil {
		return nil, errors.New("db should be integer")
	}

	hookConfig := RedisHookConfig{
		Host: settings["host"],
		Port: port,
		DB:   db,
		Key:  settings["key"],
	}
	if password, ok := settings["password"]; ok {
		hookConfig.Password = password
	} else {
		hookConfig.Password = ""
	}
	if app, ok := settings["app"]; ok {
		hookConfig.App = app
	}
	if hostname, ok := settings["hostname"]; ok {
		hookConfig.Hostname = hostname
	}
	if logformat, ok := settings["logformat"]; ok {
		hookConfig.LogFormat = logformat
	}
	if poolSize, ok := settings["poolsize"]; ok {
		pl, cErr := strconv.Atoi(poolSize)
		if cErr != nil {
			return nil, errors.New("poolsize should be integer")
		}
		hookConfig.PoolSize = pl
	} else {
		hookConfig.PoolSize = 3
	}

	hook, err := newRedisHook(hookConfig)
	if err != nil {
		return nil, err
	}

	return hook, nil
}

// NewHook creates a hook to be added to an instance of logger
func newRedisHook(config RedisHookConfig) (*RedisLogHook, error) {
	redisClient, err := newRedisClient(config.Host, config.Password, config.Port, config.DB, config.PoolSize)

	if err != nil {
		return nil, err
	}

	return &RedisLogHook{
		redisClient: redisClient,
		redisKey:    config.Key,

		app:       config.App,
		hostname:  config.Hostname,
		logFormat: config.LogFormat,
	}, nil

}

// Fire is called when a log event is fired.
func (r *RedisLogHook) Fire(entry *logrus.Entry) error {
	var msg interface{}

	switch r.logFormat {
	case "logstashv0":
		msg = createV0Message(entry, r.app, r.hostname)
	case "logstashv1":
		msg = createV1Message(entry, r.app, r.hostname)
	case "json":
	default:
		msg = createMessage(entry)
	}

	// Marshal into json message
	js, err := jsoniter.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error creating message for REDIS: %s", err)
	}

	c := r.redisClient

	// send message
	_, err = c.RPush(r.redisKey, js).Result()
	if err != nil {
		return fmt.Errorf("error sending message to REDIS: %s", err)
	}

	return nil
}

// Levels returns the available logging levels.
func (r *RedisLogHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// TODO: add sync.Pool here
func createMessage(entry *logrus.Entry) map[string]interface{} {
	m := make(map[string]interface{})

	m["message"] = entry.Message
	m["level"] = entry.Level.String()
	m["time"] = entry.Time.UTC().Format(time.RFC3339Nano)
	for k, v := range entry.Data {
		m[k] = v
	}
	return m
}

func createV0Message(entry *logrus.Entry, appName, hostname string) map[string]interface{} {
	m := make(map[string]interface{})
	m["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	m["@source_host"] = hostname
	m["@message"] = entry.Message

	fields := make(map[string]interface{})
	fields["level"] = entry.Level.String()
	fields["application"] = appName

	for k, v := range entry.Data {
		fields[k] = v
	}
	m["@fields"] = fields

	return m
}

func createV1Message(entry *logrus.Entry, appName, hostname string) map[string]interface{} {
	m := make(map[string]interface{})
	m["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	m["host"] = hostname
	m["message"] = entry.Message
	m["level"] = entry.Level.String()
	m["application"] = appName
	for k, v := range entry.Data {
		m[k] = v
	}

	return m
}

func newRedisClient(server, password string, port int, db int, poolSize int) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", server, port)

	c := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password,
		DB:          db,
		PoolSize:    poolSize,
		IdleTimeout: 180 * time.Second,
	})

	_, err := c.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to REDIS: %s", err)
	}
	return c, nil
}
