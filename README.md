A simple Terraform module deploying SES recipient rules and AWS Lambda forwarding received messages to Slack channel.
Note: message formatting is minimal, in particular MIME messages (e.g. e-mails with attachments) are delivered without any additional processing.

Example usage:
```
resource "aws_ses_receipt_rule_set" "main" {
  rule_set_name = "primary-rules"
}

module "mail2slack" {
  source         = "github.com/pmzuk/mail2slack"
  name           = "hello-example"
  webhook        = "https://hooks.slack.com/services/XXXXX/YYYYY/ZZZZZ"
  rule_set_name  = aws_ses_receipt_rule_set.main.id
  recipients     = ["hello@example.com"]
}

```
Currently messages are delivered via Slack's Incoming Webhooks. Configuration is described on [Slack's website](https://api.slack.com/messaging/webhooks).

Additional requirements:
- Python >= 3.6 (required by underlying `terraform-aws-modules/lambda/aws` module)
- Go 1.16 - currently required to build Lambda code



