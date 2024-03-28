package common

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	armNetwork "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
)

func TestNsgSubnetAssociation(t *testing.T, ctx types.TestContext) {

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID is not set in the environment variables ")
	}

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to get credentials: %e\n", err)
	}

	options := arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzurePublic,
		},
	}

	// Create network security group client
	nsgClient, err := armNetwork.NewSecurityGroupsClient(subscriptionID, credential, &options)
	if err != nil {
		t.Fatalf("Error getting NSG client: %v", err)
	}

	subnetsClient, err := armNetwork.NewSubnetsClient(subscriptionID, credential, &options)
	if err != nil {
		t.Fatalf("Error getting subnets client: %v", err)
	}

	t.Run("IsNsgSubnetAssociated", func(t *testing.T) {
		resourceGroupName := terraform.Output(t, ctx.TerratestTerraformOptions(), "resource_group_name")
		nsgName := terraform.Output(t, ctx.TerratestTerraformOptions(), "name")
		vnetNames := terraform.OutputMap(t, ctx.TerratestTerraformOptions(), "vnet_names")
		subnetNames := terraform.OutputMap(t, ctx.TerratestTerraformOptions(), "vnet_subnets")

		nsg, err := nsgClient.Get(context.Background(), resourceGroupName, nsgName, nil)
		if err != nil {
			t.Fatalf("Error getting nsg: %v", err)
		}
		if nsg.Name == nil {
			t.Fatalf("nsg does not exist")
		}

		for _, vnetName := range vnetNames {
			for _, subnetName := range subnetNames {
				inputSubnetName := strings.Trim(getSubstring(subnetName), "[]")

				subnet, err := subnetsClient.Get(context.Background(), resourceGroupName, vnetName, inputSubnetName, nil)
				if err != nil {
					t.Fatalf("Error getting subnet: %v", err)
				}
				if subnet.Name == nil {
					t.Fatalf("Subnet does not exist")
				}
				subnetNsg := subnet.Properties.NetworkSecurityGroup
				assert.NotEmpty(t, subnetNsg, "Subnet does not have a nsg associated.")
			}
		}
	})
}

func getSubstring(input string) string {
	parts := strings.Split(input, "/")
	return parts[len(parts)-1]
}
