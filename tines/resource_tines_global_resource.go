package tines

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tuckner/go-tines/tines"
)

func resourceTinesGlobalResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTinesGlobalResourceCreate,
		Read:   resourceTinesGlobalResourceRead,
		Update: resourceTinesGlobalResourceUpdate,
		Delete: resourceTinesGlobalResourceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"folder_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"global_resource_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceTinesGlobalResourceCreate(d *schema.ResourceData, meta interface{}) error {

	name := d.Get("name").(string)
	valueType := d.Get("value_type").(string)
	value := d.Get("value").(string)
	teamID := d.Get("team_id").(int)
	folderID := d.Get("folder_id").(int)

	tinesClient := meta.(*tines.Client)

	gr := tines.GlobalResource{
		Name:      name,
		ValueType: valueType,
		Value:     value,
		TeamID:    teamID,
		FolderID:  folderID,
	}

	globalresource, _, err := tinesClient.GlobalResource.Create(&gr)
	if err != nil {
		return err
	}

	sgrid := strconv.Itoa(globalresource.ID)

	d.SetId(sgrid)

	return resourceTinesGlobalResourceRead(d, meta)
}

func resourceTinesGlobalResourceRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	grid, _ := strconv.ParseInt(d.Id(), 10, 32)
	globalresource, _, err := tinesClient.GlobalResource.Get(int(grid))
	if err != nil {
		return err
	}

	sgrid := strconv.Itoa(globalresource.ID)

	d.SetId(sgrid)
	d.Set("name", globalresource.Name)
	d.Set("value", globalresource.Value)
	d.Set("value_type", globalresource.ValueType)
	d.Set("team_id", globalresource.TeamID)
	d.Set("folder_id", globalresource.FolderID)
	d.Set("global_resource_id", globalresource.ID)

	return nil
}

func resourceTinesGlobalResourceDelete(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	grid, _ := strconv.ParseInt(d.Id(), 10, 32)
	_, err := tinesClient.GlobalResource.Delete(int(grid))
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceTinesGlobalResourceUpdate(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	name := d.Get("name").(string)
	valueType := d.Get("value_type").(string)
	value := d.Get("value").(string)
	teamID := d.Get("team_id").(int)
	folderID := d.Get("folder_id").(int)
	grid, _ := strconv.ParseInt(d.Id(), 10, 32)

	gr := tines.GlobalResource{
		Name:      name,
		ValueType: valueType,
		Value:     value,
		TeamID:    teamID,
		FolderID:  folderID,
	}

	globalresource, _, err := tinesClient.GlobalResource.Update(int(grid), &gr)
	if err != nil {
		return err
	}

	sgrid := strconv.Itoa(globalresource.ID)

	d.SetId(sgrid)
	d.Set("name", globalresource.Name)
	d.Set("value", globalresource.Value)
	d.Set("value_type", globalresource.ValueType)
	d.Set("folder_id", globalresource.FolderID)
	d.Set("grid", globalresource.ID)

	return resourceTinesGlobalResourceRead(d, meta)
}
