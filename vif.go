package xen

type VIF APIObject

func (self *VIF) Destroy() error {
	result := APIResult{}
	return self.Client.APICall(&result, "VIF.destroy", self.Ref)
}

func (self *VIF) GetNetwork() (*Network, error) {
	network := new(Network)
	result := APIResult{}
	err := self.Client.APICall(&result, "VIF.get_network", self.Ref)
	if err != nil {
		return nil, err
	}
	network.Ref = result.Value.(string)
	network.Client = self.Client
	return network, nil
}
