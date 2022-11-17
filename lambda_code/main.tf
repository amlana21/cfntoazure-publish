resource "aws_lambda_function" "cust_resource_lambda" {
  function_name    = "custom_resource_lambda"
  filename         = "main.zip"
  handler          = "main"
  source_code_hash = filebase64sha256("main.zip")
  role             = "${var.role_arn}"
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = 30
  environment {
    variables = {
      directory_id = ""
      client_id = ""
      client_secret = ""
      subscription_id = ""
    }
  }
}



