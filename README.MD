# Open-IM-Demo-Register Use Tutorial

## First step:

Log in to this site: https://www.aliyun.com/product/sms

Buy Short Message Service

You will get accessKeyId and accessKeySecret

According to the developer guide

Apply for your SignName and TemplateCode

The variable in the template should be ${code}

Now you have four strings of characters: 

accessKeyId

accessKeySecret

SignName,

TemplateCode

## Second step

Open file: Open-IM-Demo-Register/config/config.xml

Find these four key-value pairs

Fill in the value you got before

| config.xml                   | aliyun.com                   |
| ---------------------------- | ---------------------------- |
| accessKeyId                  | accessKeyId                  |
| accessKeySecret              | accessKeySecret              |
| SignName                     | SignName                     |
| VerificationCodeTemplateCode | VerificationCodeTemplateCode |

## Third step

Configure your Mysql, redis, OpenIM-Server, log in config.xml

It is worth noting that superCode

When you have not purchased the verification code service, or do not want to use the verification code service

You can use this key-value pair

Enter the superCode directly in the verification code box









