package redisengine

import (
	"fmt"
	"math/rand"
	"time"

	redis "github.com/go-redis/redis"
	"github.com/morfien101/go-metrics-auth/config"
	"github.com/silverstagtech/randomstring"
)

var (
	defaultRedisOptions = redis.Options{
		MaxRetries: 3,
	}
)

// RedisEngine is used to handle requests to Redis
type RedisEngine struct {
	config config.RedisConfig
	client *redis.Client
}

// New returns a RedisEngine. When you are ready you are expected to call
// Start to make use of it, then Stop to tear down connections.
func New(c config.RedisConfig) *RedisEngine {
	return &RedisEngine{
		config: c,
	}
}

// Start will connect the redis engine to the redis server.
// It will also check that it can ping the server. It will return an error
// if the Ping fails.
func (re *RedisEngine) Start() error {
	redisOptions := defaultRedisOptions
	redisOptions.Addr = re.config.RedisHost + ":" + re.config.RedisPort
	c := redis.NewClient(&redisOptions)
	if err := c.Ping().Err(); err != nil {
		return err
	}
	re.client = c

	return nil
}

// Stop will shutdown the Redis Engine
// It will return a channel that receives an error
// if anything goes wrong. It will get a nil if its all good.
func (re *RedisEngine) Stop() <-chan error {
	c := make(chan error, 1)
	go func() {
		c <- re.client.Close()
		close(c)
	}()
	return c
}

// CreateCredentials will create a new set of credentials in redis and return them to the
// caller. The returned creds will be in a []string that will have the username, password
// in postion 0 and 1 respectively.
// As a side note this will with 128 bytes into redis.
func (re *RedisEngine) CreateCredentials() ([]string, error) {
	user, _ := randomstring.Generate(4, 4, 4, 4, 64)
	password, _ := randomstring.Generate(4, 4, 4, 4, 64)
	ok, err := re.client.SetNX(user, password, 15*time.Minute).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Failed to create auth")
	}
	return []string{user, password}, nil
}

// GetEndpoint will return a random endpoint configured in the configuration file.
func (re *RedisEngine) GetEndpoint() (string, error) {
	if len(re.config.AvailableEndpoints) < 1 {
		return "", fmt.Errorf("No Endpoint configured")
	}
	return re.config.AvailableEndpoints[rand.Intn(len(re.config.AvailableEndpoints))], nil
}
