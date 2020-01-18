package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"

	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"

	"google.golang.org/grpc"

	"github.com/firemiles/bifrost/pkg/cni"
	pb "github.com/firemiles/bifrost/rpc"
)

const (
	address = "unix:///run/bifrost.sock"
)

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, bv.BuildString("bifrost"))
}

func loadConf(bytes []byte, cmdCheck bool) (*cni.NetConf, string, error) {
	n := &cni.NetConf{}
	if err := json.Unmarshal(bytes, n); err != nil {
		return nil, "", fmt.Errorf("failed to load netconf: %v", err)
	}

	if cmdCheck {
		return n, n.CNIVersion, nil
	}

	var err error
	// Parse previous result
	if n.NetConf.RawPrevResult != nil {
		if err = version.ParsePrevResult(&n.NetConf); err != nil {
			return nil, "", fmt.Errorf("could not parse prevResult: %v", err)
		}

		_, err = current.NewResultFromResult(n.PrevResult)
		if err != nil {
			return nil, "", fmt.Errorf("could not convert result to current version: %v", err)
		}
	}

	return n, n.CNIVersion, nil
}

func cmdAdd(args *skel.CmdArgs) error {
	n, cniVersion, err := loadConf(args.StdinData, false)
	if err != nil {
		return err
	}

	netConf, err := json.Marshal(n)
	if err != nil {
		return err
	}
	var cniRequest = &pb.CNIRequest{
		Request: string(netConf),
	}
	var result *current.Result

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewBifrostBackendClient(conn)
	r, err := c.ADD(context.Background(), cniRequest)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(r.GetReply()), result)
	if err != nil {
		return err
	}
	return types.PrintResult(result, cniVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	n, _, err := loadConf(args.StdinData, false)
	if err != nil {
		return err
	}

	netConf, err := json.Marshal(n)
	if err != nil {
		return err
	}
	var cniRequest = &pb.CNIRequest{
		Request: string(netConf),
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	c := pb.NewBifrostBackendClient(conn)
	r, err := c.DEL(context.Background(), cniRequest)
	if err != nil {
		return err
	}

	message := r.GetReply()
	if message == "" {
		return nil
	}

	return errors.New(message)
}

func cmdCheck(args *skel.CmdArgs) error {
	n, _, err := loadConf(args.StdinData, false)
	if err != nil {
		return err
	}

	netConf, err := json.Marshal(n)
	if err != nil {
		return err
	}
	var cniRequest = &pb.CNIRequest{
		Request: string(netConf),
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	c := pb.NewBifrostBackendClient(conn)
	r, err := c.CHECK(context.Background(), cniRequest)
	if err != nil {
		return err
	}

	message := r.GetReply()
	if message == "" {
		return nil
	}

	return errors.New(message)
}
