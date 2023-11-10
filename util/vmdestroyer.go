package util

import (
	"context"
	"fmt"

	"github.com/b-/gomox/tasks"
	"github.com/luthermonson/go-proxmox"
	"github.com/sirupsen/logrus"
)

func DestroyVm(vm *proxmox.VirtualMachine, ctx context.Context) (proxmox.Task, error) {
	task, err := vm.Delete(ctx)
	if err != nil {
		return *task, err
	}
	err = task.Ping(context.Background())
	if err != nil {
		return *task, err
	}
	logrus.Debugf("deletion requested! %#v", task)
	return *task, nil
}

func DestroyVmWithForce(vm *proxmox.VirtualMachine, ctx context.Context) (proxmox.Task, error) {
	logrus.Trace(
		"DestroyVmWithForce(\n",
		fmt.Sprintf("    vm: %#v\n", vm), // todo: learn structured logging
		")",
	)
	if vm.IsRunning() {
		logrus.Warnf(
			"The VM %d was running!\n"+
				"Stopping before destroying.", vm.VMID,
		)
		task, err := vm.Stop(ctx)
		if err != nil {
			return *task, err
		}
		err = tasks.QuietWaitTask(
			*task,
			tasks.DefaultPollInterval,
			ctx,
		)
		if err != nil {
			return *task, err
		}
	}
	task, err := DestroyVm(vm, ctx)
	if err != nil {
		return task, err
	}
	logrus.Info(fmt.Sprintf("deletion requested! %#v", task))
	return task, err
}
