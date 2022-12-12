package ec2_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccIPAMPoolCIDRsDataSource_basic(t *testing.T) {
	dataSourceName := "data.aws_vpc_ipam_pool_cidrs.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAMPoolCIDRsDataSourceConfig_basicOneCIDRs,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ipam_pool_cidrs.#", "1"),
				),
			},
			{
				Config: testAccIPAMPoolCIDRsDataSourceConfig_basicTwoCIDRs,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ipam_pool_cidrs.#", "2"),
				),
			},
			{
				Config: testAccIPAMPoolCIDRsDataSourceConfig_basicTwoCIDRsFiltered,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "ipam_pool_cidrs.#", "1"),
				),
			},
		},
	})
}

var testAccIPAMPoolCIDRsDataSourceConfig_basicOneCIDRs = acctest.ConfigCompose(testAccIPAMPoolConfig_basic, `
resource "aws_vpc_ipam_pool_cidr" "test" {
  ipam_pool_id = aws_vpc_ipam_pool.test.id
  cidr         = "172.2.0.0/16"
}

data "aws_vpc_ipam_pool_cidrs" "test" {
  ipam_pool_id = aws_vpc_ipam_pool.test.id

  depends_on = [
    aws_vpc_ipam_pool_cidr.test
  ]
}
`)

var testAccIPAMPoolCIDRsDataSourceConfig_basicTwoCIDRs = acctest.ConfigCompose(testAccIPAMPoolConfig_basic, `
resource "aws_vpc_ipam_pool_cidr" "test" {
  ipam_pool_id = aws_vpc_ipam_pool.test.id
  cidr         = "172.2.0.0/16"
}

resource "aws_vpc_ipam_pool_cidr" "testtwo" {
  ipam_pool_id = aws_vpc_ipam_pool.test.id
  cidr         = "10.2.0.0/16"
}

data "aws_vpc_ipam_pool_cidrs" "test" {
  ipam_pool_id = aws_vpc_ipam_pool.test.id

  depends_on = [
    aws_vpc_ipam_pool_cidr.test,
    aws_vpc_ipam_pool_cidr.testtwo,
  ]
}
`)

var testAccIPAMPoolCIDRsDataSourceConfig_basicTwoCIDRsFiltered = acctest.ConfigCompose(testAccIPAMPoolConfig_basic, `
resource "aws_vpc_ipam_pool_cidr" "test" {
  ipam_pool_id = aws_vpc_ipam_pool.test.id
  cidr         = "172.2.0.0/16"
}

resource "aws_vpc_ipam_pool_cidr" "testtwo" {
  ipam_pool_id = aws_vpc_ipam_pool.test.id
  cidr         = "10.2.0.0/16"
}

data "aws_vpc_ipam_pool_cidrs" "test" {
  ipam_pool_id = aws_vpc_ipam_pool.test.id

  filter {
    name   = "cidr"
    values = ["10.*"]
  }

  depends_on = [
    aws_vpc_ipam_pool_cidr.test,
    aws_vpc_ipam_pool_cidr.testtwo,
  ]
}
`)
