---
layout: "aws"
page_title: "AWS: aws_api_gateway_integration_response"
sidebar_current: "docs-aws-resource-api-gateway-integration-response"
description: |-
  Provides an HTTP Method Integration Response for an API Gateway Resource.
---

# aws\_api\_gateway\_integration\_response

Provides an HTTP Method Integration Response for an API Gateway Resource.

## Example Usage

```
resource "aws_api_gateway_rest_api" "MyDemoAPI" {
  name = "MyDemoAPI"
  description = "This is my API for demonstration purposes"
}

resource "aws_api_gateway_resource" "MyDemoResource" {
  rest_api_id = "${aws_api_gateway_rest_api.MyDemoAPI.id}"
  parent_id = "${aws_api_gateway_rest_api.MyDemoAPI.root_resource_id}"
  path_part = "mydemoresource"
}

resource "aws_api_gateway_method" "MyDemoMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.MyDemoAPI.id}"
  resource_id = "${aws_api_gateway_resource.MyDemoResource.id}"
  http_method = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "MyDemoIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.MyDemoAPI.id}"
  resource_id = "${aws_api_gateway_resource.MyDemoResource.id}"
  http_method = "${aws_api_gateway_method.MyDemoMethod.http_method}"
  type = "MOCK"
}

resource "aws_api_gateway_method_response" "200" {
  rest_api_id = "${aws_api_gateway_rest_api.MyDemoAPI.id}"
  resource_id = "${aws_api_gateway_resource.MyDemoResource.id}"
  http_method = "${aws_api_gateway_method.MyDemoMethod.http_method}"
  status_code = "200"
}

resource "aws_api_gateway_integration_response" "MyDemoIntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.MyDemoAPI.id}"
  resource_id = "${aws_api_gateway_resource.MyDemoResource.id}"
  http_method = "${aws_api_gateway_method.MyDemoMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.200.status_code}"
}
```

## Argument Reference

The following arguments are supported:

* `rest_api_id` - (Required) API Gateway ID
* `resource_id` - (Required) API Gateway Resource ID
* `http_method` - (Required) HTTP Method (GET, POST, PUT, DELETE, HEAD, OPTION)
* `status_code` - (Required) Specify the HTTP status code.
* `response_models` - (Optional) Specifies the  Model resources used for the response's content type
* `response_parameters` - (Optional) Represents response parameters that can be sent back to the caller
