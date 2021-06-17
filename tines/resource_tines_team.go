package tines

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tuckner/go-tines/tines"
)

func resourceTinesTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTinesTeamCreate,
		Read:   resourceTinesTeamRead,
		Update: resourceTinesTeamUpdate,
		Delete: resourceTinesTeamDelete,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTinesTeamCreate(d *schema.ResourceData, meta interface{}) error {

	name := d.Get("name").(string)

	tinesClient := meta.(*tines.Client)

	n := tines.Team{
		Name: name,
	}

	team, _, err := tinesClient.Team.Create(&n)
	if err != nil {
		return err
	}

	stid := strconv.Itoa(team.ID)

	d.SetId(stid)

	return resourceTinesTeamRead(d, meta)
}

func resourceTinesTeamRead(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	tid, _ := strconv.ParseInt(d.Id(), 10, 32)
	team, _, err := tinesClient.Team.Get(int(tid))
	if err != nil {
		return err
	}

	stid := strconv.Itoa(team.ID)

	d.SetId(stid)
	d.Set("team_id", team.ID)
	d.Set("name", team.Name)

	return nil
}

func resourceTinesTeamDelete(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)
	tid, _ := strconv.ParseInt(d.Id(), 10, 32)
	_, err := tinesClient.Team.Delete(int(tid))
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceTinesTeamUpdate(d *schema.ResourceData, meta interface{}) error {

	tinesClient := meta.(*tines.Client)

	name := d.Get("name").(string)
	tid, _ := strconv.ParseInt(d.Id(), 10, 32)

	n := tines.Team{
		Name: name,
	}

	team, _, err := tinesClient.Team.Update(int(tid), &n)
	if err != nil {
		return err
	}

	stid := strconv.Itoa(team.ID)

	d.SetId(stid)

	return resourceTinesTeamRead(d, meta)
}
