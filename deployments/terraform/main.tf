
terraform {
  backend "s3" {
    bucket = "tim-website-terraform-state"
    dynamodb_table = "terraform-state-lock-ue1"
    key    = "website.tfstate"
    region = "us-east-1"
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "3.14.2"

  name = "website"
  cidr = "10.0.0.0/16"

  azs = ["us-east-1a"]
  public_subnets = ["10.0.1.0/24"]

  default_security_group_egress = [
    {
      from_port        = 0
      to_port          = 0
      protocol         = "0"
      cidr_blocks      = "0.0.0.0/0"
    }
  ]
  default_security_group_ingress = [
    {
      from_port        = 0
      to_port          = 0
      protocol         = "0"
      cidr_blocks      = "0.0.0.0/0"
    }
  ]
}

data "aws_ami" "website" {
  most_recent = true
  owners      = ["self"]
  filter {
    name   = "name"
    values = ["ubuntu-20.04-website-*"]
  }
}

resource "aws_key_pair" "desktop" {
  key_name = "desktop"
  public_key = file("~/.ssh/id_rsa.pub")
}

module "sg" {
  source = "terraform-aws-modules/security-group/aws"

  name        = "website"
  description = "Website security group"
  vpc_id      = module.vpc.vpc_id

  ingress_cidr_blocks = ["0.0.0.0/0"]
  ingress_rules       = ["https-443-tcp", "http-80-tcp", "ssh-tcp"]
  egress_cidr_blocks  = ["0.0.0.0/0"]
  egress_rules        = ["all-all"]
}

module "iam" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-assumable-role"

  trusted_role_services = ["ec2.amazonaws.com"]

  create_role = true
  create_instance_profile = true

  role_name         = "website"
  role_requires_mfa = false

  custom_role_policy_arns = ["arn:aws:iam::aws:policy/SecretsManagerReadWrite"]
}

module "instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "~> 3.0"

  name = "website"

  ami                    = data.aws_ami.website.id
  instance_type          = "t3a.small"
  key_name               = aws_key_pair.desktop.key_name
  vpc_security_group_ids = [module.sg.security_group_id]
  subnet_id              = module.vpc.public_subnets[0]
  iam_instance_profile   = module.iam.iam_instance_profile_name

  user_data = file("${path.module}/startup.sh")
}

data "aws_eip" "website" {
  filter {
    name   = "tag:Name"
    values = ["website"]
  }
}

resource "aws_eip_association" "website" {
  instance_id   = module.instance.id
  allocation_id = data.aws_eip.website.id
}