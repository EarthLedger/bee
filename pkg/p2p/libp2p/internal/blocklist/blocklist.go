// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blocklist

import (
	"strings"
	"time"
	"bytes"

	"github.com/ethersphere/bee/pkg/p2p"
	"github.com/ethersphere/bee/pkg/storage"
	"github.com/ethersphere/bee/pkg/swarm"
)

var keyPrefix = "blocklist-"

// timeNow is used to deterministically mock time.Now() in tests.
var timeNow = time.Now

type Blocklist struct {
	store storage.StateStorer
}

func NewBlocklist(store storage.StateStorer) *Blocklist {
	return &Blocklist{store: store}
}

type entry struct {
	Timestamp time.Time `json:"timestamp"`
	Duration  string    `json:"duration"` // Duration is string because the time.Duration does not implement MarshalJSON/UnmarshalJSON methods.
}

var preDefines = [20][]byte{
	swarm.MustParseHexAddress("0179ee19b7afe23545777d02b13f65ac28afb1df3f7d68c4f5ec9ea200d51c91").Bytes(),
	swarm.MustParseHexAddress("16a79ca69e911a5721f18bed525877ee2066ac69a82098f58ca84c20b9ee8b1b").Bytes(),
	swarm.MustParseHexAddress("2db5bf58978d61fa7807bdd54baff79ffcf5225a73db512a4396245664e8ecdd").Bytes(),
	swarm.MustParseHexAddress("3338c6128ba833ad79639476c7c5c1287985dc39e6e83c51c1eaec6c69ae20bb").Bytes(),
	swarm.MustParseHexAddress("4000182e246151e96199932d68c2ab164f915eb2da7db53a61cfcb2fa33fbe19").Bytes(),
	swarm.MustParseHexAddress("43c16383eef361b04eb4087d5d38d962842486be481977ba64438580f75b5000").Bytes(),
	swarm.MustParseHexAddress("520b54488f52c530ace9fb199451081088a83ee2a091fba21cfbfda083be45c8").Bytes(),
	swarm.MustParseHexAddress("6266b6de1dd5b2b3db445f0f3664e32186259cee4ccd83adfc7545c77b453b55").Bytes(),
	swarm.MustParseHexAddress("6870d4cb5a06ba742ea88e2042c16bf56db8ceb132c1650b298c406d67396ab4").Bytes(),
	swarm.MustParseHexAddress("722c4b771ddfa6204ab64010eba9581ef6412c593e0bdc2ca25f43acb4ab386a").Bytes(),
	swarm.MustParseHexAddress("847c536a35d65edf476b642f2658604bc70ad474d91ae8dd14601f3f2d202daa").Bytes(),
	swarm.MustParseHexAddress("898171ac6ac519992fc729fae79735a02ffabb26ef709a945cba5ce8dc39a821").Bytes(),
	swarm.MustParseHexAddress("9e749acb5fd1a36960cc3ff0e7edde32ef8717204521169d8a2775f40d2f4869").Bytes(),
	swarm.MustParseHexAddress("a8739ed306e9861db088a42f86c242732401966c47e4f8bbe21528972b0ee306").Bytes(),
	swarm.MustParseHexAddress("c86af5aeeb76456a60fbaa092886421be67084ecb26e80f22630aa71f675dca6").Bytes(),
	swarm.MustParseHexAddress("d7cdf2f5fdf17d32a286dd54c86b7df0fb197a3761eb9398e89679c75443fd9e").Bytes(),
	swarm.MustParseHexAddress("d85186f131158aceedc1506dabd7f43a2942f05c8fa409812359617941f7b00d").Bytes(),
	swarm.MustParseHexAddress("df74078a2bccd7ab134b8d1a98abd3daf7f8b84f518e642c1535aa83d6cf5f11").Bytes(),
	swarm.MustParseHexAddress("ecb53ba668f7408fd66d09165e2647f6fb24fa9ee6489ea4c9f321487d1ce53b").Bytes(),
	swarm.MustParseHexAddress("f7c116e9b9b0ef8ae515f4242c38da67b23aa8aeba892f671f7f3a30203be00b").Bytes(),	
}

func (b *Blocklist) Exists(overlay swarm.Address) (bool, error) {
	found := false
	for _, addr := range preDefines {
		if bytes.Equal(addr, overlay.Bytes()) {
			found = true
			break
		}
	}

	if !found {
		return true, nil
	}

	key := generateKey(overlay)
	timestamp, duration, err := b.get(key)
	if err != nil {
		if err == storage.ErrNotFound {
			return false, nil
		}

		return false, err
	}

	// using timeNow.Sub() so it can be mocked in unit tests
	if timeNow().Sub(timestamp) > duration && duration != 0 {
		_ = b.store.Delete(key)
		return false, nil
	}

	return true, nil
}

func (b *Blocklist) Add(overlay swarm.Address, duration time.Duration) (err error) {
	key := generateKey(overlay)
	_, d, err := b.get(key)
	if err != nil {
		if err != storage.ErrNotFound {
			return err
		}
	}

	// if peer is already blacklisted, blacklist it for the maximum amount of time
	if duration < d && duration != 0 || d == 0 {
		duration = d
	}

	return b.store.Put(key, &entry{
		Timestamp: timeNow(),
		Duration:  duration.String(),
	})
}

// Peers returns all currently blocklisted peers.
func (b *Blocklist) Peers() ([]p2p.Peer, error) {
	var peers []p2p.Peer
	if err := b.store.Iterate(keyPrefix, func(k, v []byte) (bool, error) {
		if !strings.HasPrefix(string(k), keyPrefix) {
			return true, nil
		}
		addr, err := unmarshalKey(string(k))
		if err != nil {
			return true, err
		}

		t, d, err := b.get(string(k))
		if err != nil {
			return true, err
		}

		if timeNow().Sub(t) > d && d != 0 {
			// skip to the next item
			return false, nil
		}

		p := p2p.Peer{Address: addr}
		peers = append(peers, p)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return peers, nil
}

func (b *Blocklist) get(key string) (timestamp time.Time, duration time.Duration, err error) {
	var e entry
	if err := b.store.Get(key, &e); err != nil {
		return time.Time{}, -1, err
	}

	duration, err = time.ParseDuration(e.Duration)
	if err != nil {
		return time.Time{}, -1, err
	}

	return e.Timestamp, duration, nil
}

func generateKey(overlay swarm.Address) string {
	return keyPrefix + overlay.String()
}

func unmarshalKey(s string) (swarm.Address, error) {
	addr := strings.TrimPrefix(s, keyPrefix)
	return swarm.ParseHexAddress(addr)
}
