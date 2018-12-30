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
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 10 * time.Second,
		},
	)

	if err != nil {
		fmt.Println("etcd conn err||", err)
		return
	}

	fmt.Println("etcd conn succ")
	defer cli.Close()

	/*resp, err := cli.Put(context.Background(), "/logagent/conf1/", "1111111111111")*/
	/*fmt.Println(resp, err)*/

	fmt.Println(cli.Get(context.Background(), "/logagent/conf/"))

	/*for {*/
	//rch := cli.Watch(context.Background(), "/logagent/conf1/")
	//for wresp := range rch {
	//for _, ev := range wresp.Events {
	//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
	//}
	//}
	/*}*/
}
