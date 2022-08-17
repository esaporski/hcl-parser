/* Providers */
provider "aws" {
  region = var.aws_region
}

data "aws_caller_identity" "current" {}

/* Locals */
locals {
  subnet_id          = "subnet-0a0000000a0a0a000"
  security_group_ids = ["sg-0a000000aa000a00a"]
}

/* AMI */
data "aws_ami" "amazon_linux_ami" {
  most_recent = true

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["137112412989"]
}

/* EC2 */
resource "aws_instance" "nginx" {
  ami           = data.aws_ami.amazon_linux.id
  instance_type = var.instance_type

  iam_instance_profile = "AmazonSSMRoleForInstancesQuickSetup"

  subnet_id              = local.subnet_id
  vpc_security_group_ids = local.security_group_ids

  user_data = templatefile("${path.module}/assets/user_data.sh", {
    redirect_url = "http://pudim.com.br"
  })

  tags = {
    Name  = var.app_name
    Owner = var.owner
  }
}
