package tines

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tuckner/go-tines/tines"
)

func resourceTinesAnnotation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTinesAnnotationCreate,
		Read:   resourceTinesAnnotationRead,
		Update: resourceTinesAnnotationUpdate,
		Delete: resourceTinesAnnotationDelete,

		Schema: map[string]*schema.Schema{
			"annotation_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"position": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"story_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceTinesAnnotationCreate(d *schema.ResourceData, meta interface{}) error {

	position := d.Get("position").(map[string]interface{})
	content := d.Get("content").(string)
	storyID := d.Get("story_id").(int)

	tinesClient := meta.(*tines.Client)

	n := tines.Annotation{
		StoryID:  storyID,
		Content:  content,
		Position: position,
	}

	annotation, _, err := tinesClient.Annotation.Create(&n)
	if err != nil {
		return err
	}

	snid := strconv.Itoa(annotation.ID)

	d.SetId(snid)

	return resourceTinesAnnotationRead(d, meta)
}

func resourceTinesAnnotationRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	nid, _ := strconv.ParseInt(d.Id(), 10, 32)
	annotation, _, err := tinesClient.Annotation.Get(int(nid))
	if err != nil {
		log.Printf("[DEBUG] Error: %v", err)
		if strings.Contains(err.Error(), "404") {
			log.Printf("[INFO] Annotation %v no longer exists", d.Id())
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	snid := strconv.Itoa(annotation.ID)

	d.SetId(snid)
	d.Set("story_id", annotation.StoryID)
	d.Set("position", annotation.Position)
	d.Set("content", annotation.Content)
	d.Set("annotation_id", annotation.ID)

	return nil
}

func resourceTinesAnnotationDelete(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	nid, _ := strconv.ParseInt(d.Id(), 10, 32)
	_, err := tinesClient.Annotation.Delete(int(nid))
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceTinesAnnotationUpdate(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	position := d.Get("position").(map[string]interface{})
	content := d.Get("content").(string)
	storyID := d.Get("story_id").(int)
	nid, _ := strconv.ParseInt(d.Id(), 10, 32)

	n := tines.Annotation{
		StoryID:  storyID,
		Content:  content,
		Position: position,
	}

	annotation, _, err := tinesClient.Annotation.Update(int(nid), &n)
	if err != nil {
		return err
	}

	snid := strconv.Itoa(annotation.ID)

	d.SetId(snid)

	return resourceTinesAnnotationRead(d, meta)
}
