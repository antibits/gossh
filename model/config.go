package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
)

const CONFIG_FILE = "./hosts.conf"

type Host struct {
	User    string
	Host    string
	Timeout int
}

func (host *Host) String() string {
	return fmt.Sprintf("%s@%s[%ds]", host.User, host.Host, host.Timeout)
}

type Conf struct {
	md5   string  `json:"-"`
	Hosts []*Host `json:"hosts"`
}

type HostNotFound struct {
}

func (*HostNotFound) Error() string {
	return "host not found"
}

func (cfg *Conf) GetHost(keyword string) ([]Host, error) {
	mHosts := make([]Host, 0)
	for _, host := range cfg.Hosts {
		if strings.HasSuffix(keyword, ";") {
			if strings.HasSuffix(host.Host, keyword[:len(keyword)-1]) {
				mHosts = append(mHosts, *host)
			}
		} else if strings.Contains(host.Host, keyword) {
			mHosts = append(mHosts, *host)
		}
	}
	if len(mHosts) == 0 {
		return nil, &HostNotFound{}
	}
	if len(mHosts) > 1 && strings.HasSuffix(keyword, ";") {
		_mHosts := make([]Host, 0)
		for _, host := range mHosts {
			if host.Host == keyword[:len(keyword)-1] {
				_mHosts = append(_mHosts, host)
				break
			}
		}
		if len(_mHosts) > 0 {
			mHosts = _mHosts
		}
	}
	return mHosts, nil
}

func (cfg *Conf) DeleteHosts(hosts []Host) {
	for _, host := range hosts {
		for i, cfHost := range cfg.Hosts {
			if cfHost.Host == host.Host {
				if i == 0 {
					if len(cfg.Hosts) > 1 {
						cfg.Hosts = cfg.Hosts[1:]
					} else {
						cfg.Hosts = []*Host{}
					}
				} else if i == len(cfg.Hosts)-1 {
					cfg.Hosts = cfg.Hosts[:i]
				} else {
					cfg.Hosts = append(cfg.Hosts[:i], cfg.Hosts[i+1:]...)
				}
				break
			}
		}
	}
}

func (cfg *Conf) NewHost(user, host string, timeoutSec int) {
	for _, h := range cfg.Hosts {
		if h.Host == host {
			h.User = user
			h.Timeout = timeoutSec
			return
		}
	}
	cfg.Hosts = append(cfg.Hosts, &Host{
		User:    user,
		Host:    host,
		Timeout: timeoutSec,
	})
}

func (cfg *Conf) Save() {
	// flush config to config file
	cfgData, _ := json.Marshal(cfg)
	if len(cfgData) > 0 {
		_md5 := md5.Sum(cfgData)
		if cfg.md5 == hex.EncodeToString(_md5[:]) {
			return
		}
		if err := ioutil.WriteFile(CONFIG_FILE, cfgData, os.ModePerm); err != nil {
			fmt.Println(err)
		}
		cfg.md5 = hex.EncodeToString(_md5[:])
	}
}

func (cfg *Conf) Load() {
	if cfgData, err := ioutil.ReadFile(CONFIG_FILE); err != nil {
		fmt.Println(err)
	} else {
		config := &Conf{}
		if err = json.Unmarshal(cfgData, &config); err != nil {
			fmt.Println(err)
		}
		_md5 := md5.Sum(cfgData)
		if cfg.md5 != hex.EncodeToString(_md5[:]) {
			cfg.md5 = hex.EncodeToString(_md5[:])
			cfg.Hosts = config.Hosts
		}
	}
}

var Config = &Conf{}

func init() {
	fmt.Println("loading configurations...")
	Config.Load()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Kill, os.Interrupt)
		select {
		case <-c:
			Config.Save()
		}
	}()
}
