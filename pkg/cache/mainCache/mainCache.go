package mainCache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

// MainRedisConn 主要使用的redis
var MainRedisConn *redis.Pool

const (
	// 表示连接池空闲连接列表的长度限制，空闲列表是一个栈式的结构，先进后出
	maxIdle = 20
	// 连接池的最大数据库连接数。设为0表示无限制。
	maxActive = 100
	// 空闲连接的超时设置，一旦超时，将会从空闲列表中摘除，该超时时间时间应该小于服务端的连接超时设置
	idleTimeout = 180 * time.Second
)

func Init() error {
	MainRedisConn = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")))
			if err != nil {
				return nil, err
			}
			// 主redis无密码 注释了这块
			if os.Getenv("REDIS_PASSWORD") != "" {
				if _, err := c.Do("AUTH", os.Getenv("REDIS_PASSWORD")); err != nil {
					c.Close()
					return nil, err
				}
			}
			// 设定默认数据库[但是放入连接池里面不能重置DataBase]
			if os.Getenv("REDIS_DATABASE") != "" {
				if _, err := c.Do("SELECT", os.Getenv("REDIS_DATABASE")); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

// GetStaticRedisCon 从连接池获取一个redis链接
func GetStaticRedisCon() (redis.Conn, error) {
	conn := MainRedisConn.Get()
	_, err := conn.Do("SELECT", 6)
	if err != nil {
		return conn, err
	}
	return conn, err
}


