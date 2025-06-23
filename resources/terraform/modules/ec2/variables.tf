variable "ami_id" {}
variable "instance_type" {}
variable "key_name" {}
variable "root_volume_size" {}
variable "prefix" {}
variable "subnet_id" {}
variable "vpc_id" {}
variable "security_group_id" {}
variable "private_key_path" {}
variable "preserve_eip" {
  type    = bool
  default = true
}
variable "rke2_version" {}
variable "cert_manager_version" {}
variable "encryption_secret_key" {}
variable "input_cluster_config" {}
variable "cattle_config" {}
variable "rancher_version" {
}

variable "rancher_password" {
  description = "Bootstrap password for Rancher"
  type        = string
  sensitive   = false
}
variable "rancher_repo_url" {
}

variable "install_rancher" {
}
