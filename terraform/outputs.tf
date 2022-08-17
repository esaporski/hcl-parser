/* Outputs */
output "instance_id" {
  description = "ID of the EC2 instance"
  value       = aws_instance.nginx.id
}

output "instance_private_ip" {
  description = "The private IP address assigned to the instance."
  value       = aws_instance.nginx.private_ip
}
