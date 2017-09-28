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

func (self *VIF) GetMAC() (string, error) {
	result := APIResult{}
	err := self.Client.APICall(&result, "VIF.get_MAC", self.Ref)
	if err != nil {
		return "", err
	}
	return result.Value.(string), nil
}

func (self *VIF) GetDevice() (string, error) {
	result := APIResult{}
	err := self.Client.APICall(&result, "VIF.get_device", self.Ref)
	if err != nil {
		return "", err
	}
	return result.Value.(string), nil
}

func (self *VIF) GetIPv4Addresses() (addresses []string, err error) {
	result := APIResult{}
	err = self.Client.APICall(&result, "VIF.get_ipv4_addresses", self.Ref)
	if err != nil {
		return nil, err
	}
	for _, elem := range result.Value.([]interface{}) {
		address := elem.(string)
		addresses = append(addresses, address)
	}
	return addresses, nil
}
