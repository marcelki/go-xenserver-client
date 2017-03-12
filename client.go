package xen

import (
	"errors"
	"fmt"

	"github.com/nilshell/xmlrpc"
)

type Client struct {
	Session  interface{}
	Host     string
	Url      string
	Username string
	Password string
	RPC      *xmlrpc.Client
}

type APIResult struct {
	Status           string
	Value            interface{}
	ErrorDescription string
}

type APIObject struct {
	Ref    string
	Client *Client
}

func (c *Client) RPCCall(result interface{}, method string, params []interface{}) (err error) {
	p := new(xmlrpc.Params)
	p.Params = params
	return c.RPC.Call(method, *p, result)

}

func (c *Client) Login() (err error) {
	result := xmlrpc.Struct{}

	params := make([]interface{}, 2)
	params[0] = c.Username
	params[1] = c.Password

	err = c.RPCCall(&result, "session.login_with_password", params)
	if err == nil {
		// err might not be set properly, so check the reference
		if result["Value"] == nil {
			return errors.New("Invalid credentials supplied")
		}
	}
	c.Session = result["Value"]
	return err
}

func (c *Client) APICall(result *APIResult, method string, params ...interface{}) (err error) {
	if c.Session == nil {
		return fmt.Errorf("No session. Unable to make call")
	}

	//Make a params slice which will include the session
	p := make([]interface{}, len(params)+1)
	p[0] = c.Session

	if params != nil {
		for idx, element := range params {
			p[idx+1] = element
		}
	}

	res := xmlrpc.Struct{}

	err = c.RPCCall(&res, method, p)

	if err != nil {
		return err
	}

	result.Status = res["Status"].(string)

	if result.Status != "Success" {
		return fmt.Errorf("API Error: %s", res["ErrorDescription"])
	} else {
		result.Value = res["Value"]
	}
	return
}

func (c *Client) GetHosts() (hosts []*Host, err error) {
	hosts = make([]*Host, 0)
	result := APIResult{}
	err = c.APICall(&result, "host.get_all")
	if err != nil {
		return hosts, err
	}
	for _, elem := range result.Value.([]interface{}) {
		host := new(Host)
		host.Ref = elem.(string)
		host.Client = c
		hosts = append(hosts, host)
	}
	return hosts, nil
}

func (c *Client) GetPools() (pools []*Pool, err error) {
	pools = make([]*Pool, 0)
	result := APIResult{}
	err = c.APICall(&result, "pool.get_all")
	if err != nil {
		return pools, err
	}

	for _, elem := range result.Value.([]interface{}) {
		pool := new(Pool)
		pool.Ref = elem.(string)
		pool.Client = c
		pools = append(pools, pool)
	}

	return pools, nil
}

func (c *Client) GetDefaultSR() (sr *SR, err error) {
	pools, err := c.GetPools()

	if err != nil {
		return nil, err
	}

	pool_rec, err := pools[0].GetRecord()

	if err != nil {
		return nil, err
	}

	if pool_rec["default_SR"] == "" {
		return nil, errors.New("No default_SR specified for the pool.")
	}

	sr = new(SR)
	sr.Ref = pool_rec["default_SR"].(string)
	sr.Client = c

	return sr, nil
}

func (c *Client) GetVMByUuid(vm_uuid string) (vm *VM, err error) {
	vm = new(VM)
	result := APIResult{}
	err = c.APICall(&result, "VM.get_by_uuid", vm_uuid)
	if err != nil {
		return nil, err
	}
	vm.Ref = result.Value.(string)
	vm.Client = c
	return
}

func (c *Client) GetHostByUuid(host_uuid string) (host *Host, err error) {
	host = new(Host)
	result := APIResult{}
	err = c.APICall(&result, "host.get_by_uuid", host_uuid)
	if err != nil {
		return nil, err
	}
	host.Ref = result.Value.(string)
	host.Client = c
	return
}

func (c *Client) GetVMByNameLabel(name_label string) (vms []*VM, err error) {
	vms = make([]*VM, 0)
	result := APIResult{}
	err = c.APICall(&result, "VM.get_by_name_label", name_label)
	if err != nil {
		return vms, err
	}

	for _, elem := range result.Value.([]interface{}) {
		vm := new(VM)
		vm.Ref = elem.(string)
		vm.Client = c
		vms = append(vms, vm)
	}

	return vms, nil
}

func (c *Client) GetVMAll() (vms []*VM, err error) {
	vms = make([]*VM, 0)
	result := APIResult{}
	err = c.APICall(&result, "VM.get_all")
	if err != nil {
		return vms, err
	}

	for _, elem := range result.Value.([]interface{}) {
		vm := new(VM)
		vm.Ref = elem.(string)
		vm.Client = c
		vms = append(vms, vm)
	}

	return vms, nil
}

func (c *Client) GetHostByNameLabel(name_label string) (hosts []*Host, err error) {
	hosts = make([]*Host, 0)
	result := APIResult{}
	err = c.APICall(&result, "host.get_by_name_label", name_label)
	if err != nil {
		return hosts, err
	}

	for _, elem := range result.Value.([]interface{}) {
		host := new(Host)
		host.Ref = elem.(string)
		host.Client = c
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func (c *Client) GetSRByNameLabel(name_label string) (srs []*SR, err error) {
	srs = make([]*SR, 0)
	result := APIResult{}
	err = c.APICall(&result, "SR.get_by_name_label", name_label)
	if err != nil {
		return srs, err
	}

	for _, elem := range result.Value.([]interface{}) {
		sr := new(SR)
		sr.Ref = elem.(string)
		sr.Client = c
		srs = append(srs, sr)
	}

	return srs, nil
}

func (c *Client) GetNetworks() (networks []*Network, err error) {
	networks = make([]*Network, 0)
	result := APIResult{}
	err = c.APICall(&result, "network.get_all")
	if err != nil {
		return nil, err
	}

	for _, elem := range result.Value.([]interface{}) {
		network := new(Network)
		network.Ref = elem.(string)
		network.Client = c
		networks = append(networks, network)
	}

	return networks, nil
}

func (c *Client) GetNetworkByUuid(network_uuid string) (network *Network, err error) {
	network = new(Network)
	result := APIResult{}
	err = c.APICall(&result, "network.get_by_uuid", network_uuid)
	if err != nil {
		return nil, err
	}
	network.Ref = result.Value.(string)
	network.Client = c
	return
}

func (c *Client) GetNetworkByNameLabel(name_label string) (networks []*Network, err error) {
	networks = make([]*Network, 0)
	result := APIResult{}
	err = c.APICall(&result, "network.get_by_name_label", name_label)
	if err != nil {
		return networks, err
	}

	for _, elem := range result.Value.([]interface{}) {
		network := new(Network)
		network.Ref = elem.(string)
		network.Client = c
		networks = append(networks, network)
	}

	return networks, nil
}

func (c *Client) GetVdiByNameLabel(name_label string) (vdis []*VDI, err error) {
	vdis = make([]*VDI, 0)
	result := APIResult{}
	err = c.APICall(&result, "VDI.get_by_name_label", name_label)
	if err != nil {
		return vdis, err
	}

	for _, elem := range result.Value.([]interface{}) {
		vdi := new(VDI)
		vdi.Ref = elem.(string)
		vdi.Client = c
		vdis = append(vdis, vdi)
	}

	return vdis, nil
}

func (c *Client) GetSRByUuid(sr_uuid string) (sr *SR, err error) {
	sr = new(SR)
	result := APIResult{}
	err = c.APICall(&result, "SR.get_by_uuid", sr_uuid)
	if err != nil {
		return nil, err
	}
	sr.Ref = result.Value.(string)
	sr.Client = c
	return
}

func (c *Client) GetVdiByUuid(vdi_uuid string) (vdi *VDI, err error) {
	vdi = new(VDI)
	result := APIResult{}
	err = c.APICall(&result, "VDI.get_by_uuid", vdi_uuid)
	if err != nil {
		return nil, err
	}
	vdi.Ref = result.Value.(string)
	vdi.Client = c
	return
}

func (c *Client) GetPIFs() (pifs []*PIF, err error) {
	pifs = make([]*PIF, 0)
	result := APIResult{}
	err = c.APICall(&result, "PIF.get_all")
	if err != nil {
		return pifs, err
	}
	for _, elem := range result.Value.([]interface{}) {
		pif := new(PIF)
		pif.Ref = elem.(string)
		pif.Client = c
		pifs = append(pifs, pif)
	}

	return pifs, nil
}

func (c *Client) CreateTask() (task *Task, err error) {
	result := APIResult{}
	err = c.APICall(&result, "task.create", "packer-task", "Packer task")

	if err != nil {
		return
	}

	task = new(Task)
	task.Ref = result.Value.(string)
	task.Client = c
	return
}

func (c *Client) CreateNetwork(name_label string, name_description string, bridge string) (network *Network, err error) {
	network = new(Network)

	net_rec := make(xmlrpc.Struct)
	net_rec["name_label"] = name_label
	net_rec["name_description"] = name_description
	net_rec["bridge"] = bridge
	net_rec["other_config"] = make(xmlrpc.Struct)

	result := APIResult{}
	err = c.APICall(&result, "network.create", net_rec)
	if err != nil {
		return nil, err
	}
	network.Ref = result.Value.(string)
	network.Client = c

	return network, nil
}

func NewClient(host, username, password string) (c Client) {
	c.Host = host
	c.Url = "http://" + host
	c.Username = username
	c.Password = password
	c.RPC, _ = xmlrpc.NewClient(c.Url, nil)
	return
}
