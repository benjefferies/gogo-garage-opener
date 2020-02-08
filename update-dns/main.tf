provider "aws" {
  # access_key = var.access_key
  # secret_key = var.secret_key
  region     = "eu-west-1"
  # allowed_account_ids = [
  #   var.account_id,
  # ]
}

terraform {
  backend "s3" {
    bucket         = "gogo-garage-door-terraform-state" // Will need to change, bucket needs to be unique
    key            = "domain/terraform.tfstate"
    region         = "eu-west-1"
    encrypt        = true
  }
}

resource "aws_route53_zone" "primary" {
  name = var.domain
}

resource "aws_route53_record" "a_record" {
  zone_id = aws_route53_zone.primary.zone_id
  name    = var.a_record_domain
  type    = "A"
  ttl     = "300"
  records = [var.ip_address]
}