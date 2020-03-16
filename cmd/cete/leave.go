// Copyright (c) 2020 Minoru Osuka
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/mosuka/cete/kvs"
	pbkvs "github.com/mosuka/cete/protobuf/kvs"
	"github.com/urfave/cli"
)

func execLeave(c *cli.Context) error {
	grpcAddr := c.String("grpc-addr")

	id := c.Args().Get(0)
	if id == "" {
		err := errors.New("id argument must be set")
		return err
	}

	req := &pbkvs.LeaveRequest{
		Id: id,
	}

	client, err := kvs.NewGRPCClient(grpcAddr)
	if err != nil {
		return err
	}
	defer func() {
		err := client.Close()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	}()

	err = client.Leave(req)
	if err != nil {
		return err
	}

	return nil
}
