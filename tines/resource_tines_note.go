package tines

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tuckner/go-tines/tines"
)

func resourceTinesNote() *schema.Resource {
	return &schema.Resource{
		Create: resourceTinesNoteCreate,
		Read:   resourceTinesNoteRead,
		Update: resourceTinesNoteUpdate,
		Delete: resourceTinesNoteDelete,

		Schema: map[string]*schema.Schema{
			"note_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"position": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"story_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceTinesNoteCreate(d *schema.ResourceData, meta interface{}) error {

	position := d.Get("position").(string)
	content := d.Get("content").(string)
	storyID := d.Get("story_id").(int)

	tinesClient := meta.(*tines.Client)

	n := tines.Note{
		StoryID:  storyID,
		Content:  content,
		Position: position,
	}

	note, _, err := tinesClient.Note.Create(&n)
	if err != nil {
		return err
	}

	snid := strconv.Itoa(note.ID)

	d.SetId(snid)

	return resourceTinesNoteRead(d, meta)
}

func resourceTinesNoteRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	nid, _ := strconv.ParseInt(d.Id(), 10, 32)
	note, _, err := tinesClient.Note.Get(int(nid))
	if err != nil {
		return err
	}

	snid := strconv.Itoa(note.ID)

	d.SetId(snid)
	d.Set("story_id", note.StoryID)
	d.Set("position", note.Position)
	d.Set("content", note.Content)
	d.Set("note_id", note.ID)

	return nil
}

func resourceTinesNoteDelete(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	nid, _ := strconv.ParseInt(d.Id(), 10, 32)
	_, err := tinesClient.Note.Delete(int(nid))
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceTinesNoteUpdate(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	position := d.Get("position").(string)
	content := d.Get("content").(string)
	storyID := d.Get("story_id").(int)
	nid, _ := strconv.ParseInt(d.Id(), 10, 32)

	n := tines.Note{
		StoryID:  storyID,
		Content:  content,
		Position: position,
	}

	note, _, err := tinesClient.Note.Update(int(nid), &n)
	if err != nil {
		return err
	}

	snid := strconv.Itoa(note.ID)

	d.SetId(snid)
	d.Set("story_id", note.StoryID)
	d.Set("position", note.Position)
	d.Set("content", note.Content)
	d.Set("note_id", note.ID)

	return resourceTinesNoteRead(d, meta)
}
