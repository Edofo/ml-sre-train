module "vpc" {
  source = "../vpc"
  env    = var.env
}

module "eks" {
  source = "../eks"
  env    = var.env

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnet
}
