
source "amazon-ebs" "website" {
  ami_name = "ubuntu-20.04-website-{{timestamp}}"
  ami_description = "Website image based on ubuntu 20.04 LTS"

  region = "us-east-1"
  
  instance_type = "t2.small"
  vpc_filter {
    filters = {
      "tag-key": "Builder" // For these, make sure to create a vpc and subnet and tag it with this. Or tag the default vpc
    }
  }
  subnet_filter {
    filters = {
      "tag-key": "Builder"
    }
  }
  temporary_key_pair_name = "ubuntu-packer-{{timestamp}}"
  ssh_username            = "ubuntu"

  source_ami_filter {
    filters = {
      name                = "ubuntu-minimal/images/hvm-ssd/ubuntu-focal-20.04-amd64-minimal*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
}

build {
  sources = ["source.amazon-ebs.website"]

  provisioner "file" {
    source      = "./scripts/" // This copies the files, not the folder itself
    destination = "/home/ubuntu/"
  }

  provisioner "shell" {
    inline = [
      "/home/ubuntu/provision.sh"
    ]
  }
}