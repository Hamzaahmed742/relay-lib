/*

  Copyright 2017 Loopring Project Ltd (Loopring Foundation).

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package broadcast

import (
	"github.com/Loopring/relay-lib/log"
	"github.com/Loopring/relay-lib/utils"
	"sync"
)

var broadcaster *Broadcaster

type Broadcaster struct {
	publishers  []Publisher
	subscribers []Subscriber
}

type Publisher interface {
	Name() string
	PubOrder(hash, orderData string) error
}

type Subscriber interface {
	Name() string
	Next() ([][]byte, error)
}

func PubOrder(hash, orderData string) map[string]error {
	var (
		errs    map[string]error
		errsMtx sync.RWMutex
		wg      sync.WaitGroup
	)
	errsMtx = sync.RWMutex{}
	for _, publisher := range broadcaster.publishers {
		wg.Add(1)
		go func(publisher Publisher) {
			defer func() {
				wg.Add(-1)
			}()
			if err := publisher.PubOrder(hash, orderData); nil != err {
				errsMtx.Lock()
				if nil == errs {
					errs = make(map[string]error)
				}
				errs[publisher.Name()] = err
				errsMtx.Unlock()
			}
		}(publisher)
	}
	wg.Wait()
	return errs
}

func SubOrderNext() (<-chan interface{}, error) {
	in, out := utils.MakeInfinite()
	for _, subscriber := range broadcaster.subscribers {
		go func(subscriber Subscriber) {
			for {
				if ordersData, err := subscriber.Next(); nil == err {
					for _, data := range ordersData {
						in <- data
					}
				} else {
					log.Errorf("occurs err:%s, when subscribing:%s order ", err.Error(), subscriber.Name())
				}
			}
		}(subscriber)
	}
	return out, nil
}

func Initialize(publishers []Publisher, subscribers []Subscriber) {
	broadcaster = &Broadcaster{}
	broadcaster.publishers = publishers
	broadcaster.subscribers = subscribers
}

func IsInit() bool {
	return broadcaster != nil
}
