package main

import (
	"context"
	"fmt"
	"hazelcast_basics/internal/common"
	"sync"
	"time"

	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/logger"
)

func distributedMapDemo() {
	client, ctx := createNewHazelcastInstance()
	defer client.Shutdown(ctx)
	mp, err := client.GetMap(ctx, "my-map")
	common.PanicIfErr(err)

	for i := 0; i < 1000; i++ {
		mp.Set(ctx, i, i)
	}

}

func distributedMapWithNoLocksDemo() {
	client, ctx := createNewHazelcastInstance()
	defer client.Shutdown(ctx)

	mp, err := client.GetMap(ctx, "my-map")
	common.PanicIfErr(err)

	for i := 0; i < 1000; i++ {
		value, _ := mp.Get(ctx, 1)
		valueInt := value.(int64) + 1
		time.Sleep(1 * time.Millisecond)

		mp.Set(ctx, 1, valueInt)
	}

	result, _ := mp.Get(ctx, 1)
	fmt.Println("Finished. Result:", result)
}

func distributedMapWithPessimisticLocksDemo() {
	client, ctx := createNewHazelcastInstance()
	defer client.Shutdown(ctx)

	mp, err := client.GetMap(ctx, "my-map")
	common.PanicIfErr(err)

	for i := 0; i < 1000; i++ {
		mp.Lock(ctx, 1)
		value, _ := mp.Get(ctx, 1)
		valueInt := value.(int64) + 1
		time.Sleep(1 * time.Millisecond)

		mp.Set(ctx, 1, valueInt)
		mp.Unlock(ctx, 1)
	}

	result, _ := mp.Get(ctx, 1)
	fmt.Println("Finished. Result:", result)
}

func distributedMapWithOptimisticLocksDemo() {
	client, ctx := createNewHazelcastInstance()
	defer client.Shutdown(ctx)

	mp, err := client.GetMap(ctx, "my-map")
	common.PanicIfErr(err)

	for i := 0; i < 1000; i++ {
		for {
			value, _ := mp.Get(ctx, 1)
			valueInt := value.(int64) + 1
			time.Sleep(1 * time.Millisecond)

			res, err := mp.ReplaceIfSame(ctx, 1, valueInt-1, valueInt)
			if res && err == nil {
				break
			}
		}

	}

	result, _ := mp.Get(ctx, 1)
	fmt.Println("Finished. Result:", result)
}

func createNewHazelcastInstance() (*hazelcast.Client, context.Context) {
	ctx := context.TODO()
	config := hazelcast.NewConfig()
	config.Logger.Level = logger.OffLevel
	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	common.PanicIfErr(err)

	return client, ctx
}

func main() {
	// --- No locks --
	fmt.Println("--- No locks ---")
	distributedMapDemo() // Reset map

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			distributedMapWithNoLocksDemo()
		}()
	}

	wg.Wait()
	fmt.Println()

	// --- Pessimistic lock --
	fmt.Println("--- Pessimistic lock ---")
	distributedMapDemo() // Reset map

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			distributedMapWithPessimisticLocksDemo()
		}()
	}

	wg.Wait()
	fmt.Println()

	// --- Optimistic lock --
	fmt.Println("--- Optimistic lock ---")
	distributedMapDemo() // Reset map

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			distributedMapWithOptimisticLocksDemo()
		}()
	}

	wg.Wait()

}
