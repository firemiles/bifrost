/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"sync"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	networkv1 "github.com/firemiles/bifrost/controller/api/v1"
)

// IPBlockReconciler reconciles a IPBlock object
type IPBlockReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	ipBlockCache     map[string][]net.IPNet // key: subnet, value: ipblocks
	ipBlockCacheLock sync.Mutex
	once             sync.Once
}

// +kubebuilder:rbac:groups=network.crd.firemiles.top,resources=ipblocks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=network.crd.firemiles.top,resources=ipblocks/status,verbs=get;update;patch

func (r *IPBlockReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("ipblock", req.NamespacedName)

	var ipBlock networkv1.IPBlock
	if err := r.Get(ctx, req.NamespacedName, &ipBlock); err != nil {
		log.Error(err, "unable to fetch IPBlock")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	r.once.Do(func() {
		r.ipBlockCache = make(map[string][]net.IPNet)
		if err := r.initializeCache(ctx); err != nil {
			log.Error(err, "refresh cache failed")
			panic(err)
		}
	})

	if ipBlock.Status.Phase == networkv1.PhaseAvailable {
		log.V(1).Info("ipblock is ready, skip reconcile", "ipblock", ipBlock)
		return ctrl.Result{}, nil
	}

	networkName := ipBlock.Spec.Network
	subnetSlice := ipBlock.Spec.SubnetSlice
	if networkName == "" {
		err := fmt.Errorf("network is empty")
		return ctrl.Result{}, err
	}
	var network networkv1.Network

	if err := r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: networkName}, &network); err != nil {
		log.Error(err, "unable to fetch network", "name", networkName)
		return ctrl.Result{}, nil
	}
	if len(network.Spec.SubnetSlices) == 0 {
		log.Error(nil, "network subnet slice is empty", "network", network)
	}
	if len(subnetSlice) == 0 {
		// auto assign first subnet slice
		for _, subnet := range network.Spec.SubnetSlices[0] {
			subnetSlice = append(subnetSlice, networkv1.SubnetCIDR{Name: subnet})
		}
	}

	var subnets = make(map[string]networkv1.Subnet)
	for _, subnetcidr := range subnetSlice {
		var subnet networkv1.Subnet
		if err := r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: subnetcidr.Name}, &subnet); err != nil {
			log.Error(err, "unable to fetch subnet %s/%s", req.Namespace, subnetcidr.Name)
			return ctrl.Result{}, nil
		}
		subnets[subnet.Name] = subnet
	}
	for i, subnetcidr := range subnetSlice {
		if subnetcidr.CIDR == "" {
			cidr, err := r.AllocateIPBlock(subnetcidr.Name, subnets[subnetcidr.Name].Spec.CIDR, ipBlock.Spec.NetMask)
			if err != nil {
				return ctrl.Result{}, err
			}
			subnetSlice[i].CIDR = cidr
		} else {
			if err := r.AddIPBlockCache(subnetcidr.Name, subnetcidr.CIDR); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	var err error
	// rollback
	defer func() {
		if err != nil {
			for _, subnetcidr := range subnetSlice {
				r.DelIPBlockCache(subnetcidr.Name, subnetcidr.CIDR)
			}
		}
	}()
	ipBlock.Spec.SubnetSlice = subnetSlice
	ipBlock.Status.Phase = networkv1.PhaseAvailable

	ipBlock.Status.Unallocated = make([]string, 1<<uint(32-ipBlock.Spec.NetMask))
	if err = r.Update(ctx, &ipBlock); err != nil {
		return ctrl.Result{}, err
	}
	if err = r.Status().Update(ctx, &ipBlock); err != nil {
		log.Error(err, "update ipblock status failed", "ipblock", ipBlock)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *IPBlockReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkv1.IPBlock{}).
		Complete(r)
}

func (r *IPBlockReconciler) initializeCache(ctx context.Context) error {
	log := r.Log
	var ipBlocks networkv1.IPBlockList
	if err := r.List(ctx, &ipBlocks); err != nil {
		return errors.WithMessage(err, "unable to list ipblocks")
	}
	for _, ipBlock := range ipBlocks.Items {
		if !r.IPBlockIsActive(&ipBlock) {
			continue
		}
		for _, subnetCIDR := range ipBlock.Spec.SubnetSlice {
			if err := r.AddIPBlockCache(subnetCIDR.Name, subnetCIDR.CIDR); err != nil {
				log.Error(err, "cidr add to subnet cache failed, need to fix", "subnet", subnetCIDR)
			}
		}
	}
	return nil
}

func (r *IPBlockReconciler) AddIPBlockCache(subnet string, cidr string) error {
	r.ipBlockCacheLock.Lock()
	defer r.ipBlockCacheLock.Unlock()
	return r.addIPBlockCache(subnet, cidr)
}

func (r *IPBlockReconciler) addIPBlockCache(subnet string, cidr string) error {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		r.Log.Error(err, "parse cidr failed", "cidr", cidr)
		return err
	}
	r.Log.Info("add cidr to cache for subnet", "cidr", cidr, "subnet", subnet)
	if r.ipBlockCache[subnet] == nil {
		r.ipBlockCache[subnet] = append(r.ipBlockCache[subnet], *ipnet)
		return nil
	}
	for _, ipn := range r.ipBlockCache[subnet] {
		if ipnet.Contains(ipn.IP) || ipn.Contains(ipnet.IP) {
			return fmt.Errorf("CIDR %s is conflicts cidr %s in subnet %s, %v", cidr, ipn.String(), subnet, r.ipBlockCache)
		}
	}
	r.ipBlockCache[subnet] = append(r.ipBlockCache[subnet], *ipnet)
	return nil
}

func (r *IPBlockReconciler) DelIPBlockCache(subnet string, cidr string) {
	r.ipBlockCacheLock.Lock()
	defer r.ipBlockCacheLock.Unlock()
	if r.ipBlockCache[subnet] != nil {
		for i, ipnet := range r.ipBlockCache[subnet] {
			if ipnet.String() == cidr {
				r.ipBlockCache[subnet] = append(r.ipBlockCache[subnet][:i], r.ipBlockCache[subnet][i+1:]...)
				return
			}
		}
	}
}

func (r *IPBlockReconciler) AllocateIPBlock(subnetName string, subnetCIDR string, netmask int) (string, error) {
	r.ipBlockCacheLock.Lock()
	defer r.ipBlockCacheLock.Unlock()
	return r.allocateIPBlock(subnetName, subnetCIDR, netmask)
}

func (r *IPBlockReconciler) allocateIPBlock(subnetName string, subnetCIDR string, netmask int) (string, error) {
	var cidr net.IPNet
	ip, subnetipnet, err := net.ParseCIDR(subnetCIDR)
	if err != nil {
		return "", err
	}
	cidr = net.IPNet{
		IP:   ip.Mask(net.CIDRMask(netmask, 32)),
		Mask: net.CIDRMask(netmask, 32),
	}
	if r.ipBlockCache[subnetName] == nil {
		return cidr.String(), r.addIPBlockCache(subnetName, cidr.String())
	}

	next := func(ipnet net.IPNet) net.IPNet {
		ipb := big.NewInt(0).SetBytes([]byte(ipnet.IP))
		ipb.Add(ipb, big.NewInt(0).Lsh(big.NewInt(1), uint(32-netmask)))
		b := ipb.Bytes()
		b = append(make([]byte, len(ipnet.IP)-len(b)), b...)
		return net.IPNet{IP: net.IP(b), Mask: ipnet.Mask}
	}

	ones, _ := subnetipnet.Mask.Size()
	delta := netmask - ones
	if delta < 0 {
		return "", fmt.Errorf("subnet cidr is not match netmask")
	}

	for count := 1 << uint(delta); count > 0; count-- {
		conflict := false
		for _, ipnet := range r.ipBlockCache[subnetName] {
			if ipnet.Contains(cidr.IP) || cidr.Contains(ipnet.IP) {
				conflict = true
				break
			}
		}
		if !conflict {
			return cidr.String(), r.addIPBlockCache(subnetName, cidr.String())
		}
		cidr = next(cidr)
		r.Log.Info("next cidr", "cidr", cidr.String())
	}

	return "", fmt.Errorf("not found cidr")
}

func (r *IPBlockReconciler) IPBlockIsActive(ipBlock *networkv1.IPBlock) bool {
	switch ipBlock.Status.Phase {
	case "", networkv1.PhasePending, networkv1.PhaseFailed:
		return false
	case networkv1.PhaseAvailable, networkv1.PhaseBound, networkv1.PhaseRelease:
		return true
	default:
		r.Log.Error(fmt.Errorf("unknow IPBlock Phase %s", ipBlock.Status.Phase), "")
	}
	return false
}
