package common

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
)

func TestComposableNsgSubnetAssociation(t *testing.T, ctx types.TestContext) {

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID is not set in the environment variables ")
	}

	t.Run("IsNsgSubnetAssociated", func(t *testing.T) {
		assocIds := terraform.OutputMap(t, ctx.TerratestTerraformOptions(), "id")
		assert.NotEmpty(t, assocIds, "Expected NSG subnet association outputs")
		for subnet, id := range assocIds {
			assert.NotEmpty(t, id, "Association ID should not be empty for subnet %s", subnet)
		}
	})
}
