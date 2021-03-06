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

// Package golory is ALL IN ONE package for go software
// development with best practice usages support
package golory

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
)

var (
	gly                       *golory
	glyLogger                 *LoggerClient
	goloryDefaultLoggerConfig = LoggerCfg{
		Debug: true,
		Level: "info",
		Path:  "./golory.log",
	}
)

// golory struct is used to hold all data.
type golory struct {
	cfg        *goloryConfig
	components *handler
	booted     bool
}

// goloryConfig is used to store golory configurations.
type goloryConfig struct {
	// golory namespace
	Golory struct {
		Debug    bool
		Logger   map[string]LoggerCfg
		Redis    map[string]RedisCfg
		Gorm     map[string]GormCfg
		TDengine map[string]TDengineCfg `toml:"tdengine"`
	}
}

func init() {
	gly = &golory{
		booted:     false,
		cfg:        &goloryConfig{},
		components: newHandler(),
	}
}

// Boot initiate components from configuration file or binary content.
// Toml, Json, Yaml supported.
func Boot(cfg interface{}) error {
	if gly.booted {
		// TODO do clear stuff
		gly.booted = false
		// return errors.New("booted")
	}
	switch cfg.(type) {
	case string:
		if err := parseFile(cfg.(string)); err != nil {
			return err
		}
	case []byte:
		if err := parseBytes(cfg.([]byte)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("exepect string or []byte, but found %v", reflect.TypeOf(cfg))
	}

	// do initiation
	if err := gly.init(); err != nil {
		return fmt.Errorf("init golory failed: %s", err)
	}
	gly.booted = true
	return nil
}

// Shutdown close components
func Shutdown() error {
	if !gly.booted {
		return nil
	}
	return gly.components.closeAll()
}

// Initate golory components from file.
func parseFile(path string) error {
	// read file to []byte
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return parseBytes(b)
}

// Initiate golory components from binary content.
func parseBytes(b []byte) error {
	if err := parseCfg(b); err != nil {
		return err
	}
	return nil
}

// Do parse config.
// It will try several formats one by one.
func parseCfg(b []byte) error {
	// try file formats
	var err error
	if err = toml.Unmarshal(b, &gly.cfg); err == nil {
		return nil
	}
	e := wrap(ErrParseCfg, err)
	if err = yaml.Unmarshal(b, &gly.cfg); err == nil {
		return nil
	}
	e = wrap(e, err)
	if err = json.Unmarshal(b, &gly.cfg); err == nil {
		return nil
	}
	return wrap(e, err)
}

// Init all components
func (g *golory) init() error {
	e := wrap(nil, g.initLogger())
	e = wrap(e, g.initRedis())
	e = wrap(e, g.initGorm())
	e = wrap(e, g.initTDengine())
	return e
}

// Init log component
func (g *golory) initLogger() error {
	if g.cfg.Golory.Logger == nil {
		// empty map
		return nil
	}

	for key, cfg := range g.cfg.Golory.Logger {
		logger, err := cfg.init()
		if err != nil {
			// fast fail
			return wrap(errors.New("init logger error"), err)
		}
		g.components.setLogger(key, logger)
	}
	return nil
}

func (g *golory) initRedis() error {
	if g.cfg.Golory.Redis == nil {
		// empty map
		return nil
	}
	for key, cfg := range g.cfg.Golory.Redis {
		c, err := cfg.init()
		if err != nil {
			return wrap(errors.New("init redis error"), err)
		}
		g.components.setRedis(key, c)
	}
	return nil
}

func (g *golory) initGorm() error {
	if g.cfg.Golory.Gorm == nil {
		return nil
	}
	for key, cfg := range g.cfg.Golory.Gorm {
		c, err := cfg.init()
		if err != nil {
			return wrap(errors.New("init gorm error"), err)
		}
		g.components.setGrom(key, c)
	}
	return nil
}

func (g *golory) initTDengine() error {
	if g.cfg.Golory.TDengine == nil {
		return nil
	}
	for key, cfg := range g.cfg.Golory.TDengine {
		c, err := cfg.init()
		if err != nil {
			return wrap(errors.New("init tdengine error"), err)
		}
		g.components.setTDengine(key, c)
	}
	return nil
}
