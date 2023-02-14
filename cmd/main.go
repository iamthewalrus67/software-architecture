package main

import (
	"context"
	"fmt"
	"hazelcast_basics/internal/common"
	"math/rand"
	"os"
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
		mp.Set(ctx, i, 0)
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

func consumerQueue(id int) {
	client, ctx := createNewHazelcastInstance()
	defer client.Shutdown(ctx)

	queue, err := client.GetQueue(ctx, "queue")
	common.PanicIfErr(err)

	for {
		item, _ := queue.Take(ctx)
		fmt.Printf("consumer %d item: %d\n", id, item)
		if item.(int64) == -1 {
			queue.Put(ctx, -1)
			break
		}
		time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
	}

	fmt.Printf("Consumer %d finished\n", id)
}

func producerQueue() {
	client, ctx := createNewHazelcastInstance()
	defer client.Shutdown(ctx)

	queue, err := client.GetQueue(ctx, "queue")
	common.PanicIfErr(err)
	queue.Clear(ctx)

	for i := 0; i < 100; i++ {
		queue.Put(ctx, i)
		time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
	}

	queue.Put(ctx, -1)
	fmt.Println("Producer finished")
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
	mapDemo, queueDemo := false, false
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-m", "--map":
			mapDemo = true
			break
		case "-q", "--queue":
			queueDemo = true
			break
		default:
			panic(fmt.Sprintln("Unknown argument", arg))
		}
	}

	var wg sync.WaitGroup

	// --- Map demonstration ---
	if mapDemo {
		fmt.Println("Distributed map demonstration")

		// --- No locks --
		fmt.Println("--- No locks ---")
		distributedMapDemo() // Reset map

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

	// --- Bounded queue ---
	if queueDemo {
		if mapDemo {
			fmt.Println() // Additional new line between tests
		}
		fmt.Println("Bounded queue demonstration")

		// Create consumer goroutines
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				consumerQueue(i)
			}(i)
		}

		// Create producer goroutines
		wg.Add(1)
		go func() {
			defer wg.Done()

			producerQueue()
		}()

		wg.Wait()
	}
}
