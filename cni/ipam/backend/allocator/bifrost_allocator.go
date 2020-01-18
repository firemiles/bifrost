package allocator

import (
	"context"
	"fmt"
	"github.com/containernetworking/cni/pkg/types/current"
	networkv1 "github.com/firemiles/bifrost/controller/api/v1"
	"k8s.io/apimachinery/pkg/types"
	"math/big"
	"net"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BifrostIPAllocator
type BifrostIPAllocator struct {
	network string
	k8sClient client.Client
}

// NewBifrostIPAllocator
func NewBifrostIPAllocator(network string, k8sClient client.Client) *BifrostIPAllocator {
	return &BifrostIPAllocator{
		network: network,
		k8sClient: k8sClient,
	}
}

// Get allocates an IP
func (a *BifrostIPAllocator) Get(id string, ifname string, requestedIP net.IP) (*current.IPConfig, error) {
	var gw net.IP
	var subnet *networkv1.Subnet
	var ipBlockList networkv1.IPBlockList

	if requestedIP != nil {
		if err := canonicalizeIP(&requestedIP); err != nil {
			return nil, err
		}
	}

	err := a.k8sClient.List(context.Background(), &ipBlockList)
	if err != nil {
		return nil, err
	}
	var ipBlock *networkv1.IPBlock
	var cidr *net.IPNet
	for _, ipb := range ipBlockList.Items {
		if ipb.Status.Phase != networkv1.PhaseAvailable {
			continue
		}
		if ipb.Spec.Network == a.network {
			for _, subnetcidr := range ipb.Spec.SubnetSlice {
				_, ipnet, err := net.ParseCIDR(subnetcidr.CIDR)
				if err != nil {
					continue
				}
				if requestedIP == nil {
					for index, value := range ipb.Status.Unallocated {
						if value == "" {
							i := ipToInt(ipnet.IP)
							requestedIP = intToIP(i.Add(i, big.NewInt(int64(index))))
							break
						}
					}
				}
				if ipnet.Contains(requestedIP) {
					ipBlock = &ipb
					cidr = ipnet
					var s networkv1.Subnet
					err := a.k8sClient.Get(context.Background(), types.NamespacedName{Namespace: ipBlock.Namespace, Name: subnetcidr.Name}, &s)
					if err != nil {
						return nil, err
					}
					subnet = &s
					gw = net.ParseIP(subnet.Spec.GateWay)
					goto FOUND
				}
			}
		}
	}

FOUND:
	if ipBlock == nil {
		return nil, fmt.Errorf("not found ip %s in network %s", requestedIP, a.network)
	}

	i := ipToInt(requestedIP)
	delta := int(i.Sub(i, ipToInt(cidr.IP)).Int64())
	if delta < len(ipBlock.Status.Unallocated) && ipBlock.Status.Unallocated[delta] == "" {
		ipBlock.Status.Unallocated[delta] = id
		if ipBlock.Status.Allocations == nil {
			ipBlock.Status.Allocations = make(map[string]int)
		}
		ipBlock.Status.Allocations[id] = delta
		err := a.k8sClient.Status().Update(context.Background(), ipBlock)
		if err != nil {
			return nil, err
		}
		version := "4"
		if requestedIP.To4() == nil {
			version = "6"
		}
		return &current.IPConfig{
			Version: version,
			Address: net.IPNet{
				IP:   requestedIP,
				Mask: cidr.Mask,
			},
			Gateway: gw,
		}, nil
	}
	return nil, fmt.Errorf("nof found free ip in network %s", a.network)
}

// Release clears all IPs allocated for the container with given ID
func (a *BifrostIPAllocator) Release(id string, ifname string) error {
	var ipBlockList networkv1.IPBlockList
	err := a.k8sClient.List(context.Background(), &ipBlockList)
	if err != nil {
		return err
	}
	var ipBlock *networkv1.IPBlock
	for _, ipb := range ipBlockList.Items {
		if ipb.Status.Phase != networkv1.PhaseAvailable {
			continue
		}
		if ipb.Spec.Network == a.network {
			if ipb.Status.Allocations == nil {
				continue
			}
			if _, ok := ipb.Status.Allocations[id]; ok {
				ipBlock = &ipb
				break
			}
		}
	}
	if ipBlock == nil {
		return nil
	}
	ipBlock.Status.Unallocated[ipBlock.Status.Allocations[id]] = ""
	delete(ipBlock.Status.Allocations, id)
	err = a.k8sClient.Status().Update(context.Background(), ipBlock)
	return err
}

func ipToInt(ip net.IP) *big.Int {
	if v := ip.To4(); v != nil {
		return big.NewInt(0).SetBytes(v)
	}
	return big.NewInt(0).SetBytes(ip.To16())
}

func intToIP(i *big.Int) net.IP {
	return net.IP(i.Bytes())
}

// canonicalizeIP makes sure a provided ip is in standard form
func canonicalizeIP(ip *net.IP) error {
	if ip.To4() != nil {
		*ip = ip.To4()
		return nil
	} else if ip.To16() != nil {
		*ip = ip.To16()
		return nil
	}
	return fmt.Errorf("IP %s not v4 nor v6", *ip)
}

