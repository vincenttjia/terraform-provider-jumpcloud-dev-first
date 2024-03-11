// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UserGroupDataSource{}

func NewUserGroupDataSource() datasource.DataSource {
	return &UserGroupDataSource{}
}

// UserGroupDataSource defines the data source implementation.
type UserGroupDataSource struct {
	client *jcapiv2.APIClient
}

// UserGroupDataSourceModel describes the data source data model.
type UserGroupDataSourceModel struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
}

func (d *UserGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (d *UserGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Jumpcloud User Group Data Source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "User group ID",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "User group name",
				Required:            true,
			},
		},
	}
}

func (d *UserGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jcapiv2.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jcapiv2.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *UserGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	UserGroup, _, err := d.client.UserGroupsApi.GroupsUserList(ctx, "application/json", "application/json", map[string]interface{}{
		"filter": []string{"name:eq:" + data.Name.ValueString(), "type:eq:user_group"},
	})

	data.Id = types.StringValue(UserGroup[0].Id)
	data.Name = types.StringValue(UserGroup[0].Name)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
