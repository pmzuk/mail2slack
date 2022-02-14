data "aws_caller_identity" "current" {}

resource "aws_sns_topic" "sns" {
  name_prefix = var.name
}

resource "aws_sns_topic_policy" "sns" {
  arn    = aws_sns_topic.sns.arn
  policy = jsonencode({
    Version = "2008-10-17"
    Id = "__default_policy_ID"
    Statement = [{
    Sid = "__default_statement_ID"
    Effect = "Allow"
    Principal = {
      Service = "ses.amazonaws.com"
    },
    Action = [
      "SNS:Publish"
    ],
    Resource = aws_sns_topic.sns.arn
    Condition = {
      StringEquals = {
        "AWS:SourceOwner": data.aws_caller_identity.current.account_id
      }
    }
    }]
  })
}

module "lambda" {
  source                = "terraform-aws-modules/lambda/aws"
  version               = "2.34.0"
  function_name         = var.name
  runtime               = "go1.x"
  handler               =  "mail2slack"
  environment_variables = {
    SLACK_URL = var.webhook
  }
  source_path           = [{
    path        = "${path.module}/lambda",
    commands    = ["CGO_ENABLED=0 go build", ":zip"],
  }]


  allowed_triggers = {
    sns = {
      principal  = "sns.amazonaws.com"
      source_arn = aws_sns_topic.sns.arn
    }
  }

  create_current_version_allowed_triggers = false
}

resource "aws_sns_topic_subscription" "lambda" {
  topic_arn = aws_sns_topic.sns.arn
  protocol  = "lambda"
  endpoint  = module.lambda.lambda_function_arn
}


resource "aws_ses_receipt_rule" "this" {
  name          = var.name
  rule_set_name = var.rule_set_name
  recipients    = var.recipients
  enabled       = true
  scan_enabled  = true

  sns_action {
    topic_arn = aws_sns_topic.sns.arn
    position  = 1
    encoding  = "Base64"
  }
}
