# Sending Logs to Kinesis Firehose using Kinesis Agent

## Project Objective
The objective of this project is to create AMI using Packer, then use Terraform to provision EC2 instances, and finally Ansible to configure the instances instead of writing shell scripts.

**1. Edit `secrets.json`.** 

In the packer folder edit `secrets.json`, give your AMI a name and AWS region to use. I'm using `us-west-1` region in this example.

_NOTE:_ The timestamp (epoch format) will be appended to the name of the AMI.

EX: `my-ami-1605800314` 

**2. Build AMI image.**

Packer is configured to use the latest Amazon Linux 2 AMI (HVM and ebs backed) as the base for my AMI.
```
cd packer
packer build -var-file secrets.json ami.json
```

**3. Edit `secrets.tfvars`.** 

Inside of the terraform folder, edit `secrets.tfvars`. I'm using `us-west-1` as the region in this example. Change this depending on the region you're deploying to.
Supply the following:
- AWS ID
- SSH Key name
- Desired security to use
- AMI name (use the same name you specified in `secrets.json`)

**4. Build the EC2 instance and Kinesis Firehose Stream.**
```
cd terraform
terraform init
terraform plan -var-file secrets.tfvars -out apply.tfplan
terraform apply apply.tfplan
