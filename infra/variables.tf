variable "aws_region" {
  description = "The AWS region to deploy to"
  default     = "us-east-1"
}

variable "lambda_function_name" {
  description = "The name of the Lambda function"
  default     = "mongodb-backup"
}

variable "lambda_zip_path" {
  description = "The path to the Lambda deployment package"
  default     = "./lambda.zip"
}

variable "mongodb_uri" {
  description = "MongoDB URI"
}

variable "mongodb_parallel" {
  description = "MongoDB parallel setting"
  default     = "10"
}
