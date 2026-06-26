data "aws_caller_identity" "current" {}

resource "aws_iam_role" "bedrock_direct_call_sonnet" {
  name = "bedrock-direct-call-sonnet-4-6"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  inline_policy {
    name = "bedrock-invoke-policy"

    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Effect = "Allow"
          Action = [
            "bedrock:InvokeModel",
            "bedrock:InvokeModelWithResponseStream"
          ]
          Resource = [
            "arn:aws:bedrock:eu-central-1:${data.aws_caller_identity.current.account_id}:inference-profile/eu.anthropic.claude-sonnet-4-6",
            "arn:aws:bedrock:eu-central-1:${data.aws_caller_identity.current.account_id}:foundation-model/anthropic.claude-sonnet-4-6",
            "arn:aws:bedrock:eu-central-1::foundation-model/anthropic.claude-sonnet-4-6",
            "arn:aws:bedrock:eu-central-2::foundation-model/anthropic.claude-sonnet-4-6",
            "arn:aws:bedrock:eu-north-1::foundation-model/anthropic.claude-sonnet-4-6",
            "arn:aws:bedrock:eu-north-1::foundation-model/zai.glm-5",
            "arn:aws:bedrock:eu-south-1::foundation-model/anthropic.claude-sonnet-4-6",
            "arn:aws:bedrock:eu-south-2::foundation-model/anthropic.claude-sonnet-4-6",
            "arn:aws:bedrock:eu-west-1::foundation-model/anthropic.claude-sonnet-4-6",
            "arn:aws:bedrock:eu-west-2::foundation-model/anthropic.claude-sonnet-4-6",
            "arn:aws:bedrock:eu-west-3::foundation-model/anthropic.claude-sonnet-4-6"
          ]
        }
      ]
    })
  }

  tags = {
    Name        = "bedrock-direct-call-sonnet-4-6"
    ManagedBy   = "Terraform"
    Environment = "production"
  }
}

output "role_arn" {
  description = "ARN of the Bedrock IAM role"
  value       = aws_iam_role.bedrock_direct_call_sonnet.arn
}

output "role_name" {
  description = "Name of the Bedrock IAM role"
  value       = aws_iam_role.bedrock_direct_call_sonnet.name
}
