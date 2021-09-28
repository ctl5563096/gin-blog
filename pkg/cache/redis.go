package cache

import (
	"fmt"
	"gin-blog/pkg/cache/mainCache"
	"gin-blog/pkg/util"
	"log"
	"time"
)


/* 区分两种使用场景：
notice:暂时我也不懂 只能简单用一下
- 高频调用的场景，需要尽量压榨redis的性能：
调高MaxIdle的大小，该数目小于maxActive，由于作为一个缓冲区一样的存在，扩大缓冲区自然没有问题
调高MaxActive，考虑到服务端的支持上限，尽量调高
IdleTimeout由于是高频使用场景，设置短一点也无所谓，需要注意的一点是MaxIdle设置的长了，队列中的过期连接可能会增多，这个时候IdleTimeout也要相应变化
- 低频调用的场景，调用量远未达到redis的负载，稳定性为重：
MaxIdle可以设置的小一些
IdleTimeout相应地设置小一些
MaxActive随意，够用就好，容易检测到异常 */
const (
	// 表示连接池空闲连接列表的长度限制，空闲列表是一个栈式的结构，先进后出
	maxIdle = 30
	// 连接池的最大数据库连接数。设为0表示无限制。
	maxActive = 1000
	// 空闲连接的超时设置，一旦超时，将会从空闲列表中摘除，该超时时间时间应该小于服务端的连接超时设置
	idleTimeout = 200 * time.Millisecond
)


//	初始化函数
func Init()  {
	fmt.Println(">开始初始化缓存连接池...")
	if err := mainCache.Init(); err != nil {
		util.WriteLog("redis_connect_error",4,"链接主redis失败,失败原因" + err.Error())
		log.Fatal(err)
	}
	fmt.Println(">>>初始化缓存连接池完成")
}

