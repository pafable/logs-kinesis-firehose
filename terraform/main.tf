terraform {
  required_version = ">= 0.12.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = var.region
}

data "aws_ami" "amzl2" {
  most_recent = true

  filter {
    name   = "name"
    values = ["${var.ami-name}-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  owners = [var.aws_id]
}

resource "aws_instance" "web" {
  ami                    = data.aws_ami.amzl2.id
  instance_type          = "t2.micro"
  key_name               = var.key
  iam_instance_profile   = var.ec2_role
  vpc_security_group_ids = [var.aws_sg]

  provisioner "local-exec" {
    command = "sleep 60; ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -u ec2-user --private-key ~/.ssh/${var.key}.pem -i '${aws_instance.web.public_ip},' ../ansible-playbooks/001-nginx.yml"
  }

  tags = {
    Name = "${var.ami-name}-ec2"
  }
}