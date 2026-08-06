package common

import (
	"context"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	armnetwork "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
)

func TestComposableNsgSubnetAssociation(t *testing.T, ctx types.TestContext) {

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID is not set in the environment variables ")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to get credentials: %e\n", err)
	}

	clientFactory, err := armnetwork.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		t.Fatalf("Unable to get clientFactory: %e\n", err)
	}

	subnetsClient := clientFactory.NewSubnetsClient()
	nsgClient := clientFactory.NewSecurityGroupsClient()

	resourceGroupName := terraform.Output(t, ctx.TerratestTerraformOptions(), "resource_group_name")
	nsgName := terraform.Output(t, ctx.TerratestTerraformOptions(), "name")
	subnetIDs := terraform.OutputMap(t, ctx.TerratestTerraformOptions(), "subnet_ids")

	t.Run("IsNsgSubnetAssociated", func(t *testing.T) {

		nsg, err := nsgClient.Get(context.Background(), resourceGroupName, nsgName, nil)
		if err != nil {
			t.Fatalf("Error getting NSG: %v", err)
		}
		if nsg.ID == nil {
			t.Fatalf("NSG ID is nil")
		}
		expectedNsgID := *nsg.ID

		for _, subnetID := range subnetIDs {
			parsedSubnetID, err := arm.ParseResourceID(subnetID)
			if err != nil {
				t.Fatalf("Error parsing subnet ID %q: %v", subnetID, err)
			}

			subnet, err := subnetsClient.Get(
				context.Background(),
				parsedSubnetID.ResourceGroupName,
				parsedSubnetID.Parent.Name,
				parsedSubnetID.Name,
				nil,
			)
			if err != nil {
				t.Fatalf("Error getting subnet: %v", err)
			}
			if subnet.Name == nil {
				t.Fatalf("Subnet does not exist")
			}
			assert.NotNil(t, subnet.Properties.NetworkSecurityGroup, "Subnet does not have an NSG associated.")
			assert.NotNil(t, subnet.Properties.NetworkSecurityGroup.ID, "Subnet NSG ID is nil.")
			assert.Equal(t, expectedNsgID, *subnet.Properties.NetworkSecurityGroup.ID, "Subnet is not associated with the expected NSG.")
		}
	})
}
