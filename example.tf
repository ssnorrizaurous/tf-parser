data "aws_vpc" "main" {
  filter {
    name   = "tag:Name"
    values = ["arene-spl-visualizer-stg-ap-northeast-1"]
  }
}

data "aws_subnets" "private" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.main.id]
  }
  tags = {
    Name = "*internal*"
  }
}

module "spl-visualizer-database" {
  source = "../../../../../../../modules/woven-nxt-gen/aws/neptune"

  cluster_name       = local.woven_app
  private_subnet_ids = data.aws_subnets.private.ids
  replication_count  = 1
  vpc_id             = data.aws_vpc.main.id
  trusted_cidrs      = local.vpc_cidrs
  woven_app          = local.woven_app
  woven_env          = local.env_name
  woven_org_code     = local.woven_org_code
}

module "expose-control-plane-infra-database" {
  source        = "../../../../../../../modules/woven-nxt-gen/aws/db-privatelink"
  rds_port      = module.spl-visualizer-database.neptune_port
  rds_dns       = module.spl-visualizer-database.writer_endpoint
  listener_port = module.spl-visualizer-database.neptune_port
  endpoint_allowed_principals_arns = [
    "arn:aws:iam::862662394037:root"
  ]
  private_subnet_ids_one_per_az = data.aws_subnets.private.ids
  service_name                  = "${local.woven_app}-neptune"
  vpc_id                        = data.aws_vpc.main.id
  woven_app                     = local.woven_app
  woven_env                     = local.env_name
  woven_org_code                = local.woven_org_code
  vpc_cidrs                     = local.vpc_cidrs
}
