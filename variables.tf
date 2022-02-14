variable "webhook" {
  description = "Webhook address"
  type        = string
  default     = ""
}

variable "name" {
  description = "Lambda name"
  type        = string
  default     = "mail2slack"
}

variable "rule_set_name" {
  description = "SES rule set to attach"
  type        = string
  default     = ""
}

variable "recipients" {
  description = "Matching addresses"
  type        = list(string)
  default     = []
}
