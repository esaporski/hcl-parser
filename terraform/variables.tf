/* Variables */
variable "aws_region" {
  type        = string
  description = "This is the AWS region."
  default     = "us-east-2"
}

variable "environment_name" {
  type        = string
  description = "Application environment name."
}

variable "app_name" {
  type        = string
  description = "Application name."
}

variable "owner" {
  type        = string
  description = "Your name."
}

variable "instance_type" {
  type        = string
  description = "Type of instance to start."
}
