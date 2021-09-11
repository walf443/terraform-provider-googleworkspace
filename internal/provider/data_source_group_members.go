package googleworkspace

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGroupMembers() *schema.Resource {
	// Generate datasource schema from resource
	myscheme := map[string]*schema.Schema{
		"group_id": {
			Description: "Identifies the group in the API request. The value can be the group's email address, " +
				"group alias, or the unique group ID.",
			Type:     schema.TypeString,
			Required: true,
		},
		"members": {
			Description: "member list",
			Type:        schema.TypeList,
			Required:    true,
		},
	}
	dsSchema := datasourceSchemaFromResourceSchema(myscheme)
	addRequiredFieldsToSchema(dsSchema, "group_id")

	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Group Members data source in the Terraform Googleworkspace provider.",

		ReadContext: dataSourceGroupMembersRead,

		Schema: dsSchema,
	}
}

func dataSourceGroupMembersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*apiClient)

	directoryService, diags := client.NewDirectoryService()
	if diags.HasError() {
		return diags
	}

	membersService, diags := GetMembersService(directoryService)
	if diags.HasError() {
		return diags
	}

	groupId := d.Get("group_id").(string)
	members, err := membersService.List(groupId).Do()
	if err != nil {
		return handleNotFoundError(err, d, d.Id())
	}

	d.Set("group_id", groupId)
	d.Set("members", members)

	return diags
}
