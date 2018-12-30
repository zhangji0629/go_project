package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints:   []string{"localhost:2379", "localhost:2380"},
			DialTimeout: 10 * time.Second,
		},
	)

	if err != nil {
		fmt.Println("etcd conn err||", err)
		return
	}

	fmt.Println("etcd conn succ")
	defer cli.Close()

	ctx := context.Background()
	_, err = cli.Put(ctx, "/logagent/conf/", "sample_value")
	if err != nil {
		fmt.Println("set err|", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "/logagent/conf/")
	cancel()
	if err != nil {
		fmt.Println("get err||", err)
		return
	}

	fmt.Println("hhh")
	for _, ev := range resp.Kvs {
		fmt.Println(string(ev.Key), string(ev.Value), string(ev.Lease), string(ev.Version))
	}
}
