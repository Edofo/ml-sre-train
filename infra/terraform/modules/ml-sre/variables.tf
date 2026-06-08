variable "env" {
  type        = string
  description = "The environment to deploy the VPC to"

  validation {
    condition     = contains(["dev", "prod"], var.env)
    error_message = "Invalid environment. Must be one of: dev, prod."
  }
}
