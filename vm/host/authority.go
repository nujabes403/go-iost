package host

import (
	"encoding/json"

	"github.com/iost-official/go-iost/account"
	"github.com/iost-official/go-iost/core/contract"
)

type Authority struct {
	h *Host
}

func (h *Authority) RequireAuth(id, p string) (bool, *contract.Cost) {
	authList := h.h.ctx.Value("auth_list")
	authMap := authList.(map[string]int)

	return h.requireAuth(id, p, authMap)
}

func (h *Authority) ReadAuth(id string) (*account.Account, *contract.Cost) {
	acc, cost := h.h.GlobalMapGet("iost.auth", "account", id)
	if acc == nil {
		return nil, cost
	}
	var a account.Account
	err := json.Unmarshal([]byte(acc.(string)), &a)
	if err != nil {
		panic(err)
	}
	return &a, cost
}

func (h *Authority) requireAuth(id, permission string, auth map[string]int) (bool, *contract.Cost) {
	a, c := h.ReadAuth(id)
	if a == nil {
		return false, c
	}

	p, ok := a.Permissions[permission]
	if !ok {
		p = a.Permissions["active"]
	}

	u := p.Users
	for _, g := range p.Groups {
		u = append(u, g.Users...)
	}

	var weight int
	for _, user := range u {
		if user.IsKeyPair {
			if _, ok := auth[user.ID]; ok {
				weight += user.Weight
				if weight >= p.Threshold {
					return true, c
				}
			}
		} else {
			ok, cost := h.requireAuth(user.ID, user.Permission, auth)
			c.AddAssign(cost)
			if ok {
				weight += user.Weight
				if weight >= p.Threshold {
					return true, c
				}
			}
		}
	}

	return weight >= p.Threshold, c
}
