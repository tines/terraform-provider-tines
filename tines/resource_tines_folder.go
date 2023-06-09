package tines

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tines/go-tines/tines"
)

func resourceTinesFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceTinesFolderCreate,
		Read:   resourceTinesFolderRead,
		Update: resourceTinesFolderUpdate,
		Delete: resourceTinesFolderDelete,

		Schema: map[string]*schema.Schema{
			"folder_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceTinesFolderCreate(d *schema.ResourceData, meta interface{}) error {

	name := d.Get("name").(string)
	contentType := d.Get("content_type").(string)
	teamID := d.Get("team_id").(int)

	tinesClient := meta.(*tines.Client)

	n := tines.Folder{
		Name:        name,
		ContentType: contentType,
		TeamID:      teamID,
	}

	folder, _, err := tinesClient.Folder.Create(&n)
	if err != nil {
		return err
	}

	sfid := strconv.Itoa(folder.ID)

	d.SetId(sfid)

	return resourceTinesFolderRead(d, meta)
}

func resourceTinesFolderRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	fid, _ := strconv.ParseInt(d.Id(), 10, 32)
	folder, _, err := tinesClient.Folder.Get(int(fid))
	if err != nil {
		log.Printf("[DEBUG] Error: %v", err)
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] Folder %v no longer exists", d.Id())
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	sfid := strconv.Itoa(folder.ID)

	d.SetId(sfid)
	d.Set("folder_id", folder.ID)
	d.Set("name", folder.Name)
	d.Set("team_id", folder.TeamID)
	d.Set("content_type", folder.ContentType)
	d.Set("size", folder.Size)

	return nil
}

func resourceTinesFolderDelete(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	fid, _ := strconv.ParseInt(d.Id(), 10, 32)
	_, err := tinesClient.Folder.Delete(int(fid))
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceTinesFolderUpdate(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	name := d.Get("name").(string)
	contentType := d.Get("content_type").(string)
	teamID := d.Get("team_id").(int)
	fid, _ := strconv.ParseInt(d.Id(), 10, 32)

	n := tines.Folder{
		Name:        name,
		ContentType: contentType,
		TeamID:      teamID,
	}

	folder, _, err := tinesClient.Folder.Update(int(fid), &n)
	if err != nil {
		return err
	}

	sfid := strconv.Itoa(folder.ID)

	d.SetId(sfid)

	return resourceTinesFolderRead(d, meta)
}
