package ecs

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
)

type stepAttachKeyPar struct {
}

func (s *stepAttachKeyPar) Run(state multistep.StateBag) multistep.StepAction {
	keyPairName := state.Get("keyPair").(string)
	if keyPairName == "" {
		return multistep.ActionContinue
	}
	ui := state.Get("ui").(packer.Ui)
	client := state.Get("client").(*ecs.Client)
	config := state.Get("config").(Config)
	instance := state.Get("instance").(*ecs.InstanceAttributesType)

	err := client.AttachKeyPair(&ecs.AttachKeyPairArgs{RegionId: common.Region(config.AlicloudRegion), KeyPairName: keyPairName,
		InstanceIds: "[\"" + instance.InstanceId + "\"]"})
	if err != nil {
		err := fmt.Errorf("Error attaching keypair %s to instance %s : %s", keyPairName, instance.InstanceId, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf("Attach keypair %s to instance: %s", keyPairName, instance.InstanceId))

	return multistep.ActionContinue
}

func (s *stepAttachKeyPar) Cleanup(state multistep.StateBag) {
	keyPairName := state.Get("keyPair").(string)
	if keyPairName == "" {
		return
	}
	client := state.Get("client").(*ecs.Client)
	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)
	instance := state.Get("instance").(*ecs.InstanceAttributesType)

	err := client.DetachKeyPair(&ecs.DetachKeyPairArgs{RegionId: common.Region(config.AlicloudRegion), KeyPairName: keyPairName,
		InstanceIds: "[\"" + instance.InstanceId + "\"]"})
	if err != nil {
		err := fmt.Errorf("Error Detaching keypair %s to instance %s : %s", keyPairName, instance.InstanceId, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return
	}

	ui.Message(fmt.Sprintf("Detach keypair %s from instance: %s", keyPairName, instance.InstanceId))

}
