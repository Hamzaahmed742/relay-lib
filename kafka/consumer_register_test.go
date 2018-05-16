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

package kafka_test

import (
	"fmt"
	"github.com/Loopring/relay-lib/kafka"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	address := "127.0.0.1:9092"
	register := &kafka.ConsumerRegister{}
	register.Initialize(address)
	err := register.RegisterTopicAndHandler("test", "group1", TestData{}, func(data interface{}) error {
		dataValue := data.(*TestData)
		fmt.Printf("Msg : %s, Timestamp : %s \n", dataValue.Msg, dataValue.Timestamp)
		return nil
	})
	if err != nil {
		fmt.Errorf("Failed register")
		println(err)
	}
	time.Sleep(1000 * time.Second)

	defer func() {
		register.Close()
	}()
}
