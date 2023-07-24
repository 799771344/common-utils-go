package db

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisFramework struct {
	Client *redis.Client
}

func NewRedisFramework(host, port string, db int) (*RedisFramework, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // 这里可以设置Redis的密码
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisFramework{Client: client}, nil
}

func (rf *RedisFramework) Get(key string) (string, error) {
	return rf.Client.Get(key).Result()
}

func (rf *RedisFramework) Set(key string, value interface{}, expiration time.Duration) error {
	return rf.Client.Set(key, value, expiration).Err()
}

func (rf *RedisFramework) Delete(keys ...string) error {
	return rf.Client.Del(keys...).Err()
}

func (rf *RedisFramework) Exists(key string) (int64, error) {
	return rf.Client.Exists(key).Result()
}

func (rf *RedisFramework) Expire(key string, expiration time.Duration) error {
	return rf.Client.Expire(key, expiration).Err()
}

func main() {
	redisFramework, err := NewRedisFramework("localhost", "6379", 0)
	if err != nil {
		panic(err)
	}

	redisFramework.Set("foo", "bar", time.Minute)

	value, err := redisFramework.Get("foo")
	if err != nil {
		panic(err)
	}
	fmt.Println(value)

	if exists, err := redisFramework.Exists("foo"); err != nil {
		panic(err)
	} else if exists {
		fmt.Println("foo exists")
	} else {
		fmt.Println("foo does not exist")
	}

	redisFramework.Delete("foo")
}
