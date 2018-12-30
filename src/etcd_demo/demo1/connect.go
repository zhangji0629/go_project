package main

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints:   []string{"localhost:2379", "localhost:4001"},
			DialTimeout: 5 * time.Second,
		},
	)

	if err != nil {
		fmt.Println("etcd conn err||", err)
		return
	}

	fmt.Println("etcd conn succ")
	defer cli.Close()
}
