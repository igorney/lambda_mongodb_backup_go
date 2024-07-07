output "lambda_function_name" {
  value = aws_lambda_function.mongodb_backup.function_name
}

output "lambda_function_arn" {
  value = aws_lambda_function.mongodb_backup.arn
}
