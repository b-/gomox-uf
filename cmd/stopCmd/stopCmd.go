package stopCmd

import (
	"context"
	"fmt"

	"github.com/b-/gomox-uf/internal"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "stop",
	Usage:  "Stop a virtual machine",
	Action: stopVm,
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:     "vmid",
			Usage:    "`VMID` to stop",
			Required: true,
			Aliases:  []string{"v"},
			Action: func(c *cli.Context, vmid uint64) error {
				if vmid < 100 || vmid > 999999999 {
					return fmt.Errorf("VM vmid %d out of range", vmid)
				}
				return nil
			},
		},
		&cli.BoolFlag{
			Name:  "idempotent",
			Usage: "Don't return error if VM is already in requested state",
			Value: false,
		},
	},
}

func stopVm(c *cli.Context) error {
	requestedState := internal.RequestableState(proxmox.StatusVirtualMachineStopped)
	client := internal.InstantiateClient(
		internal.GetPveUrl(c),
		proxmox.Credentials{
			Username: c.String("pveuser"),
			Password: c.String("pvepassword"),
			Realm:    c.String("pverealm"),
		},
	)
	vmid := c.Uint64("vmid")

	vm, err := internal.GetVirtualMachineByVMID(vmid, client, c.Context)
	if err != nil {
		return err
	}

	if vm.IsStopped() {
		msg := fmt.Sprintf("VM %d already in requested state (%s)", vm.VMID, vm.Status)
		switch c.Bool("idempotent") {
		case true:
			logrus.Warn(msg)
			return nil
		case false:
			return fmt.Errorf(msg)
		}
	}
	task, err := internal.RequestState(
		internal.StateRequestParams{RequestedState: requestedState, Vm: vm},
		context.Background(),
	)
	if err != nil {
		return err
	}

	err = task.Ping(context.Background())
	if err != nil {
		return err
	}

	logrus.Info(fmt.Sprintf("state requested! %#v", task))
	return nil
}