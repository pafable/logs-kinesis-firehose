resource "aws_kinesis_firehose_delivery_stream" "mystream" {
  name        = "${var.ami-name}-delivery-stream-${var.region}"
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn        = aws_iam_role.firehose_role.arn
    bucket_arn      = aws_s3_bucket.bucket.arn
    buffer_size     = 5
    buffer_interval = 100
  }
}

resource "aws_s3_bucket" "bucket" {
  bucket        = "${var.ami-name}-bucket-${var.region}"
  acl           = "private"
  force_destroy = true
}

resource "aws_iam_role" "firehose_role" {
  name = "firehose_test_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "firehose.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}