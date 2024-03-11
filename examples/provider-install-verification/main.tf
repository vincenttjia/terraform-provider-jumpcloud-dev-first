terraform {
  required_providers {
    jumpcloud = {
      source = "hashicorp.com/vincenttjia/jumpcloud"
    }
  }
}

provider "jumpcloud" {
    api_key = ""
    organization_id = ""
}

data "jumpcloud_user_group" "this" {
    name = "AWS"
}

resource "jumpcloud_user_group" "this" {
  name = "Test Provider"
}
