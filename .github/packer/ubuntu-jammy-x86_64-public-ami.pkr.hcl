packer {
  required_plugins {
    amazon = {
      source  = "github.com/hashicorp/amazon"
      version = "~> 1"
    }
    ansible = {
      source  = "github.com/hashicorp/ansible"
      version = "~> 1"
    }
  }
}

variable "skip_create_ami" {
  type    = string
  default = "${env("SKIP_CREATE_AMI")}"
}

variable "tag" {
  type    = string
  default = "${env("TAG")}"
}

variable "version" {
  type    = string
  default = "jammy-22.04"
}

data "amazon-ami" "autogenerated_1" {
  filters = {
    architecture        = "x86_64"
    name                = "ubuntu/images/*ubuntu-${var.version}-*-server-*"
    root-device-type    = "ebs"
    virtualization-type = "hvm"
  }
  most_recent = true
  owners      = ["099720109477"]
  region      = "us-east-1"
}

locals {
  skip_create_ami = var.skip_create_ami == "True"
  timestamp = regex_replace(timestamp(), "[- TZ:]", "") 
  clean_name = regex_replace(timestamp(), "[^a-zA-Z0-9-]", "-")
}

source "amazon-ebs" "autogenerated_1" {
  ami_groups      = ["all"]
  ami_name        = "public-metal-ubuntu-${var.version}-${var.tag}-${local.timestamp}"
  instance_type   = "c5.large"
  region          = "us-east-1"
  skip_create_ami = local.skip_create_ami
  source_ami      = "${data.amazon-ami.autogenerated_1.id}"
  ssh_username    = "ubuntu"
  tags = {
    Base_AMI_Name = "{{ .SourceAMIName }}"
    Name          = "public-metal-ubuntu-${var.version}-${var.tag}-${local.clean_name}"
    Release       = "${var.version}"
  }
}

build {
  sources = ["source.amazon-ebs.autogenerated_1"]

  provisioner "shell" {
    inline = ["while [ ! -f /var/lib/cloud/instance/boot-finished ]; do echo 'Waiting for cloud-init...'; sleep 1; done", "wait_apt=$(ps aux | grep apt | wc -l)", "while [ \"$wait_apt\" -gt \"1\" ]; do echo \"waiting for apt to be ready....\"; wait_apt=$(ps aux | grep apt | wc -l); sleep 5; done", "sudo apt-get -y update", "sudo apt-get install -y python3-boto3 golang"]
  }

  provisioner "ansible" {
    extra_arguments = ["-e", "component=public-ami build=packer os_release=jammy tag=${var.tag}"]
    playbook_file   = ".github/packer/create_public_ami.yml"
    roles_path      = ".github/packer/roles/"
    use_proxy       = false
  }

  provisioner "shell" {
    execute_command = "sudo bash -x {{ .Path }}"
    script          = ".github/packer/clean-public-ami.sh"
  }

}
