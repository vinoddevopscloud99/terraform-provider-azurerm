package tests

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMNotificationHubNamespace_free(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNotificationHubNamespace_free(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubNamespaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

// Remove in 2.0
func TestAccAzureRMNotificationHubNamespace_freeClassic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNotificationHubNamespace_freeClassic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubNamespaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

// Remove in 2.0
func TestAccAzureRMNotificationHubNamespace_freeNotDefined(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_notification_hub_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMNotificationHubNamespace_freeNotDefined(data),
				ExpectError: regexp.MustCompile("either 'sku_name' or 'sku' must be defined in the configuration file"),
			},
		},
	})
}

func TestAccAzureRMNotificationHubNamespace_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_notification_hub_namespace", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNotificationHubNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNotificationHubNamespace_free(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNotificationHubNamespaceExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMNotificationHubNamespace_requiresImport),
		},
	})
}

func testCheckAzureRMNotificationHubNamespaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).NotificationHubs.NamespacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		namespaceName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, namespaceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on notificationNamespacesClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Notification Hub Namespace does not exist: %s", namespaceName)
		}

		return nil
	}
}

func testCheckAzureRMNotificationHubNamespaceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).NotificationHubs.NamespacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_notification_hub_namespace" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		namespaceName := rs.Primary.Attributes["name"]
		resp, err := client.Get(ctx, resourceGroup, namespaceName)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Notification Hub Namespace still exists:%s", *resp.Name)
		}
	}

	return nil
}

func testAccAzureRMNotificationHubNamespace_free(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "acctestnhn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  namespace_type      = "NotificationHub"

  sku_name = "Free"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMNotificationHubNamespace_freeClassic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "acctestnhn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  namespace_type      = "NotificationHub"

  sku {
    name = "Free"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMNotificationHubNamespace_freeNotDefined(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "acctestnhn-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  namespace_type      = "NotificationHub"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMNotificationHubNamespace_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMNotificationHubNamespace_free(data)
	return fmt.Sprintf(`
%s

resource "azurerm_notification_hub_namespace" "import" {
  name                = "${azurerm_notification_hub_namespace.test.name}"
  resource_group_name = "${azurerm_notification_hub_namespace.test.resource_group_name}"
  location            = "${azurerm_notification_hub_namespace.test.location}"
  namespace_type      = "${azurerm_notification_hub_namespace.test.namespace_type}"

  sku_name = "Free"
}
`, template)
}
