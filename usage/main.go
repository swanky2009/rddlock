package main

import (
	"fmt"

	"github.com/everfore/rddlock"
	"gopkg.in/ezbuy/redis-orm.v1/orm"
	redis "gopkg.in/redis.v5"
)

func main() {
	// rds, err := orm.NewRedisClient("localhost", 32768, "", 0)
	rds, err := orm.NewRedisClusterClient(&redis.ClusterOptions{
		Addrs: []string{"XXX", "XXX", "XXX"},
	})
	if err != nil {
		panic(err)
	}

	lock_key := "lock-key"

	locked, ex := rddlock.Lock(rds, lock_key, 5)
	if locked {
		fmt.Printf("LOCK %s: %+v\n", lock_key, locked)
		unlocked := rddlock.UnLock(rds, lock_key, ex)
		if unlocked {
			fmt.Printf("UNLOCK %s: %+v\n", lock_key, unlocked)
		} else {
			unlocked = rddlock.UnLockUnsafe(rds, lock_key)
			fmt.Printf("UNLOCK-UNSAFE %s: %+v\n", lock_key, unlocked)
		}
	}

	// retry lock

	// 1. lock the key first
	locked, _ = rddlock.Lock(rds, lock_key, 5)
	fmt.Printf("FIRST step, LOCK %s:%+v\n", lock_key, locked)
	// 2. retry to lock the locked key
	locked, _ = rddlock.LockRetry(rds, lock_key, 100, 100)
	fmt.Printf("SECOND step, LOCK-RETRY %s:%+v\n", lock_key, locked)
}

// Output
// LOCK lock-key: true
// UNLOCK lock-key: true
// FIRST step, LOCK lock-key:true
// SECOND step, LOCK-RETRY lock-key:true
