/*

Copyright 2020 The Vouch Proxy Authors.
Use of this source code is governed by The MIT License (MIT) that
can be found in the LICENSE file. Software distributed under The
MIT License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
OR CONDITIONS OF ANY KIND, either express or implied.

*/

package jwtmanager

import (
	"encoding/json"
	"testing"

	"github.com/nholuongut/vouch-proxy/pkg/cfg"
	"github.com/nholuongut/vouch-proxy/pkg/structs"

	"github.com/stretchr/testify/assert"
)

var (
	u1 = structs.User{
		Username: "test@testing.com",
		Name:     "Test Name",
	}
	t1 = structs.PTokens{
		PAccessToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IjRvaXU4In0.eyJzdWIiOiJuZnlmZSIsImF1ZCI6ImltX29pY19jbGllbnQiLCJqdGkiOiJUOU4xUklkRkVzUE45enU3ZWw2eng2IiwiaXNzIjoiaHR0cHM6XC9cL3Nzby5tZXljbG91ZC5uZXQ6OTAzMSIsImlhdCI6MTM5MzczNzA3MSwiZXhwIjoxMzkzNzM3MzcxLCJub25jZSI6ImNiYTU2NjY2LTRiMTItNDU2YS04NDA3LTNkMzAyM2ZhMTAwMiIsImF0X2hhc2giOiJrdHFvZVBhc2praVY5b2Z0X3o5NnJBIn0.g1Jc9DohWFfFG3ppWfvW16ib6YBaONC5VMs8J61i5j5QLieY-mBEeVi1D3vr5IFWCfivY4hZcHtoJHgZk1qCumkAMDymsLGX-IGA7yFU8LOjUdR4IlCPlZxZ_vhqr_0gQ9pCFKDkiOv1LVv5x3YgAdhHhpZhxK6rWxojg2RddzvZ9Xi5u2V1UZ0jukwyG2d4PRzDn7WoRNDGwYOEt4qY7lv_NO2TY2eAklP-xYBWu0b9FBElapnstqbZgAXdndNs-Wqp4gyQG5D0owLzxPErR9MnpQfgNcai-PlWI_UrvoopKNbX0ai2zfkuQ-qh6Xn8zgkiaYDHzq4gzwRfwazaqA",
		PIdToken:     "eyJhbGciOiJSUzI1NiIsImtpZCI6IjRvaXU4In0.eyJzdWIiOiJuZnlmZSIsImF1ZCI6ImltX29pY19jbGllbnQiLCJqdGkiOiJUOU4xUklkRkVzUE45enU3ZWw2eng2IiwiaXNzIjoiaHR0cHM6XC9cL3Nzby5tZXljbG91ZC5uZXQ6OTAzMSIsImlhdCI6MTM5MzczNzA3MSwiZXhwIjoxMzkzNzM3MzcxLCJub25jZSI6ImNiYTU2NjY2LTRiMTItNDU2YS04NDA3LTNkMzAyM2ZhMTAwMiIsImF0X2hhc2giOiJrdHFvZVBhc2praVY5b2Z0X3o5NnJBIn0.g1Jc9DohWFfFG3ppWfvW16ib6YBaONC5VMs8J61i5j5QLieY-mBEeVi1D3vr5IFWCfivY4hZcHtoJHgZk1qCumkAMDymsLGX-IGA7yFU8LOjUdR4IlCPlZxZ_vhqr_0gQ9pCFKDkiOv1LVv5x3YgAdhHhpZhxK6rWxojg2RddzvZ9Xi5u2V1UZ0jukwyG2d4PRzDn7WoRNDGwYOEt4qY7lv_NO2TY2eAklP-xYBWu0b9FBElapnstqbZgAXdndNs-Wqp4gyQG5D0owLzxPErR9MnpQfgNcai-PlWI_UrvoopKNbX0ai2zfkuQ-qh6Xn8zgkiaYDHzq4gzwRfwazaqA",
	}

	lc VouchClaims

	claimjson = `{
		"sub": "f:a95afe53-60ba-4ac6-af15-fab870e72f3d:mrtester",
		"groups": ["Website Users", "Test Group"],
		"given_name": "Mister",
		"family_name": "Tester",
		"email": "mrtester@test.int"
	}`
	customClaims = structs.CustomClaims{}
)

func init() {
	// log.SetLevel(log.DebugLevel)

	cfg.InitForTestPurposes()
	Configure()

	lc = VouchClaims{
		u1.Username,
		Sites,
		customClaims.Claims,
		t1.PAccessToken,
		t1.PIdToken,
		StandardClaims,
	}
	json.Unmarshal([]byte(claimjson), &customClaims.Claims)
}

func TestClaims(t *testing.T) {
	populateSites()
	log.Debugf("jwt config %s %d", string(cfg.Cfg.JWT.Secret), cfg.Cfg.JWT.MaxAge)
	assert.NotEmpty(t, cfg.Cfg.JWT.Secret)
	assert.NotEmpty(t, cfg.Cfg.JWT.MaxAge)

	// now := time.Now()
	// d := time.Duration(ExpiresAtMinutes) * time.Minute
	// log.Infof("lc d %s", d.String())
	// lc.StandardClaims.ExpiresAt = now.Add(time.Duration(ExpiresAtMinutes) * time.Minute).Unix()
	// log.Infof("lc expiresAt %d", now.Unix()-lc.StandardClaims.ExpiresAt)
	uts := CreateUserTokenString(u1, customClaims, t1)
	utsParsed, _ := ParseTokenString(uts)
	log.Infof("utsParsed: %+v", utsParsed)
	log.Infof("Sites: %+v", Sites)
	assert.True(t, SiteInToken(cfg.Cfg.Domains[0], utsParsed))
}
