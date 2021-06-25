package types

import (
	providertypes "github.com/sentinel-official/hub/x/provider/types"
)

type Provider struct {
	Address     string `json:"address"`
	Name        string `json:"name"`
	Identity    string `json:"identity"`
	Website     string `json:"website"`
	Description string `json:"description"`
}

func NewProviderFromRaw(v *providertypes.Provider) Provider {
	return Provider{
		Address:     v.Address,
		Name:        v.Name,
		Identity:    v.Identity,
		Website:     v.Website,
		Description: v.Description,
	}
}

type Providers []Provider

func NewProvidersFromRaw(v providertypes.Providers) Providers {
	items := make(Providers, 0, len(v))
	for i := 0; i < len(v); i++ {
		items = append(items, NewProviderFromRaw(&v[i]))
	}

	return items
}
