// Copyright 2018 golory Authors @1pb.club. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package golory

import (
	"github.com/go-redis/redis"
)

// RedisCfg wrapped go-redis options and some other golory options
type RedisCfg struct {
	Addr string
	// TODO cluster options
	// TODO logger option?
}

// RedisClient wrapped go-redis RedisClient
type RedisClient struct {
	*redis.Client
}

// init redis connection from redis config
func (cfg *RedisCfg) init() (*RedisClient, error) {
	return &RedisClient{
		redis.NewClient(&redis.Options{
			Addr: cfg.Addr,
		}),
	}, nil
}

// close underlying redis client
func (r *RedisClient) close() error {
	return r.Close()
}
