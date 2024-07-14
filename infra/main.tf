provider "aws" {
  region = var.aws_region
}

resource "aws_lambda_function" "mongodb_backup" {
  function_name = var.lambda_function_name
  role          = aws_iam_role.lambda_exec.arn
  handler       = "main"
  runtime       = "provided.al2"
  filename      = var.lambda_zip_path

  environment {
    variables = {
      AWS_REGION       = var.aws_region
      MONGODB_URI      = var.mongodb_uri
      MONGODB_PARALLEL = var.mongodb_parallel
    }
  }
}

resource "aws_iam_role" "lambda_exec" {
  name = "lambda_exec_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action    = "sts:AssumeRole",
      Effect    = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
