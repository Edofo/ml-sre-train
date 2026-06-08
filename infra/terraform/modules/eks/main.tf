module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.23"

  name               = "eks-ml-sre-${var.env}"
  kubernetes_version = "1.33"

  endpoint_public_access = true

  enable_cluster_creator_admin_permissions = true

  compute_config = {
    enabled    = true
    node_pools = ["general-purpose"]
  }

  vpc_id     = var.vpc_id
  subnet_ids = var.subnet_ids

  tags = {
    Service     = "ml-sre"
    Terraform   = "true"
    Environment = var.env
  }
}
