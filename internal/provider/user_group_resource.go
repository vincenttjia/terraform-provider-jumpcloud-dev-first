// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserGroupResource{}
var _ resource.ResourceWithImportState = &UserGroupResource{}

func NewUserGroupResource() resource.Resource {
	return &UserGroupResource{}
}

// UserGroupResource defines the resource implementation.
type UserGroupResource struct {
	client *jcapiv2.APIClient
}

// UserGroupResourceModel describes the resource data model.
type UserGroupResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r *UserGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (r *UserGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Jumpcloud User Group Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "User group ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "User group name",
				Required:            true,
			},
		},
	}
}

func (r *UserGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *UserGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	UserGroups, _, err := r.client.UserGroupsApi.GroupsUserPost(ctx, "application/json", "application/json", map[string]interface{}{
		"body": jcapiv2.UserGroupPost{
			Name: data.Name.ValueString(),
		},
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error Create", err.Error())
		return
	}

	data.Id = types.StringValue(UserGroups.Id)
	data.Name = types.StringValue(UserGroups.Name)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	UserGroup, _, err := r.client.UserGroupsApi.GroupsUserGet(ctx, data.Id.ValueString(), "application/json", "application/json", nil)

	data.Id = types.StringValue(UserGroup.Id)
	data.Name = types.StringValue(UserGroup.Name)

	if err != nil {
		resp.Diagnostics.AddError("Client Error Read", err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := r.client.UserGroupsApi.GroupsUserPut(ctx, data.Id.ValueString(), "application/json", "application/json", map[string]interface{}{
		"body": jcapiv2.UserGroupPut{
			Name: data.Name.ValueString(),
		},
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error Update", err.Error())
		return
	}

	// Read updated data from the API
	UserGroup, _, err := r.client.UserGroupsApi.GroupsUserGet(ctx, data.Id.ValueString(), "application/json", "application/json", nil)

	data.Id = types.StringValue(UserGroup.Id)
	data.Name = types.StringValue(UserGroup.Name)

	if err != nil {
		resp.Diagnostics.AddError("Client Error Read", err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UserGroupsApi.GroupsUserDelete(ctx, data.Id.ValueString(), "application/json", "application/json", nil)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

}

func (r *UserGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
