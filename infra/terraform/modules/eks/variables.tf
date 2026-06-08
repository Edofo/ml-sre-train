variable "env" {
  type        = string
  description = "The environment to deploy the EKS cluster to"

  validation {
    condition     = contains(["dev", "prod"], var.env)
    error_message = "Invalid environment. Must be one of: dev, prod."
  }
}

variable "vpc_id" {
  type        = string
  description = "The ID of the VPC to deploy the EKS cluster to"

  validation {
    condition     = var.vpc_id != ""
    error_message = "VPC ID is required."
  }
}

variable "subnet_ids" {
  type        = list(string)
  description = "The IDs of the subnets to deploy the EKS cluster to"

  validation {
    condition     = length(var.subnet_ids) >= 2
    error_message = "Must have at least two subnets."
  }
}
