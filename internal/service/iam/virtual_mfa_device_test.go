package iam_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfiam "github.com/hashicorp/terraform-provider-aws/internal/service/iam"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccVirtualMfaDevice_basic(t *testing.T) {
	var conf iam.VirtualMFADevice
	resourceName := "aws_iam_virtual_mfa_device.test"

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iam.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckVirtualMfaDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMfaDeviceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVirtualMfaDeviceExists(resourceName, &conf),
					acctest.CheckResourceAttrGlobalARN(resourceName, "arn", "iam", fmt.Sprintf("mfa/%s", rName)),
					resource.TestCheckResourceAttrSet(resourceName, "base_32_string_seed"),
					resource.TestCheckResourceAttrSet(resourceName, "qr_code_png"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"path", "virtual_mfa_device_name", "base_32_string_seed", "qr_code_png"},
			},
		},
	})
}

func TestAccVirtualMfaDevice_tags(t *testing.T) {
	var conf iam.VirtualMFADevice
	resourceName := "aws_iam_virtual_mfa_device.test"

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iam.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckVirtualMfaDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMfaDeviceConfigTags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVirtualMfaDeviceExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"path", "virtual_mfa_device_name", "base_32_string_seed", "qr_code_png"},
			},
			{
				Config: testAccVirtualMfaDeviceConfigTags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVirtualMfaDeviceExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccVirtualMfaDeviceConfigTags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVirtualMfaDeviceExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccVirtualMfaDevice_disappears(t *testing.T) {
	var conf iam.VirtualMFADevice
	resourceName := "aws_iam_virtual_mfa_device.test"

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iam.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckVirtualMfaDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMfaDeviceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVirtualMfaDeviceExists(resourceName, &conf),
					acctest.CheckResourceDisappears(acctest.Provider, tfiam.ResourceVirtualMfaDevice(), resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfiam.ResourceVirtualMfaDevice(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckVirtualMfaDeviceDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).IAMConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_iam_virtual_mfa_device" {
			continue
		}

		output, err := tfiam.FindVirtualMfaDevice(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if output != nil {
			return fmt.Errorf("IAM Virtual MFA Device (%s) still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckVirtualMfaDeviceExists(n string, res *iam.VirtualMFADevice) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Virtual MFA Device name is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).IAMConn

		output, err := tfiam.FindVirtualMfaDevice(conn, rs.Primary.ID)
		if err != nil {
			return err
		}

		*res = *output

		return nil
	}
}

func testAccVirtualMfaDeviceConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_virtual_mfa_device" "test" {
  virtual_mfa_device_name = %[1]q
}
`, rName)
}

func testAccVirtualMfaDeviceConfigTags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_iam_virtual_mfa_device" "test" {
  virtual_mfa_device_name = %[1]q

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccVirtualMfaDeviceConfigTags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_iam_virtual_mfa_device" "test" {
  virtual_mfa_device_name = %[1]q

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
