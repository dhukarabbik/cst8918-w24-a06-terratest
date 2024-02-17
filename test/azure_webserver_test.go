package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "a7085dd8-a65b-48a5-841f-f653e242d534"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "dhuk0001",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
}
func TestNICExistsAndConnectedToVM(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "dhuk0001",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	nicID := azure.GetNetworkInterfaceID(t, vmName, resourceGroupName, subscriptionID)
	assert.NotNil(t, nicID, "Network Interface not found for VM")
	assert.True(t, azure.NetworkInterfaceIsAttachedToVM(t, nicID, vmName, resourceGroupName, subscriptionID), "Network Interface is not attached to VM")
}

func TestCorrectUbuntuVersionRunningOnVM(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "dhuk0001",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	correctUbuntuVersion := "22_04-lts-gen2"
	actualUbuntuVersion := azure.GetVirtualMachineOSVersion(t, vmName, resourceGroupName, subscriptionID)
	assert.Equal(t, correctUbuntuVersion, actualUbuntuVersion, "Incorrect Ubuntu version running on the VM")
}
