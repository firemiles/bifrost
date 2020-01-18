package main

import (
	"encoding/json"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"

	"github.com/firemiles/bifrost/cni/ipam/backend/allocator"
	networkv1 "github.com/firemiles/bifrost/controller/api/v1"
)
var (
	scheme   = runtime.NewScheme()
	_ = networkv1.AddToScheme(scheme)
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
}

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, bv.BuildString("host-local-bifrost"))
}

func loadNetConf(bytes []byte) (*types.NetConf, string, error) {
	n := &types.NetConf{}
	if err := json.Unmarshal(bytes, n); err != nil {
		return nil, "", fmt.Errorf("failed to load netconf: %v", err)
	}
	return n, n.CNIVersion, nil
}

func cmdCheck(args *skel.CmdArgs) error {
	return nil
}

func cmdAdd(args *skel.CmdArgs) error {
	ipamConf, confVersion, err := allocator.LoadIPAMConfig(args.StdinData, args.Args)
	if err != nil {
		return err
	}
	if ipamConf.KubeConfigPath == "" {
		return fmt.Errorf("kubeConfigPath is empty. stdin: %s, args: %s", args.StdinData, args.Args)
	}
	os.Setenv("KUBECONFIG", ipamConf.KubeConfigPath)
	restConfig, err := config.GetConfig()
	if err != nil {
		return err
	}
	k8sclient, err := client.New(restConfig, client.Options{Scheme: scheme})
	if err != nil {
		return err
	}

	result := &current.Result{}

	if ipamConf.ResolvConf != "" {
		dns, err := parseResolvConf(ipamConf.ResolvConf)
		if err != nil {
			return err
		}
		result.DNS = *dns
	}

	allocator := allocator.NewBifrostIPAllocator(ipamConf.Name, k8sclient)

	ipconfig, err := allocator.Get(args.ContainerID, args.IfName, nil)
	if err != nil {
		return fmt.Errorf("reserve ip for network %s failed. %v", ipamConf.Name, err)
	}
	result.IPs = append(result.IPs, ipconfig)
	result.Routes = ipamConf.Routes

	return types.PrintResult(result, confVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	ipamConf, _, err := allocator.LoadIPAMConfig(args.StdinData, args.Args)
	if err != nil {
		return err
	}
	if ipamConf.KubeConfigPath == "" {
		return fmt.Errorf("kubeConfigPath is empty")
	}
	os.Setenv("KUBECONFIG", ipamConf.KubeConfigPath)
	restConfig, err := config.GetConfig()
	if err != nil {
		return err
	}
	k8sclient, err := client.New(restConfig, client.Options{Scheme: scheme})
	if err != nil {
		return err
	}
	allocator := allocator.NewBifrostIPAllocator(ipamConf.Name, k8sclient)

	err = allocator.Release(args.ContainerID, args.IfName)
	if err != nil {
		return err
	}
	return nil
}