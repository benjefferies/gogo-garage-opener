output "ip_address" {
  value = aws_route53_record.a_record.records
}

output "domain" {
  value = aws_route53_record.a_record.name
}
