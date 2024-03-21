package modules

import (
	"encoding/json"
	"log"
	"os"
)

var (
	WowLogChanId   string
	WowChanId      string
	BotChanId      string
	RoleId         string
	TriggerCommand string
	RoleBoost      string
	PendaRole      string
	PendaGoldRole  string
	SWCommand      string
	SWCRoleId      string

	config *configFile
)

type configFile struct {
	WowLogChanId   string `json:"wow_log_chan_id"`
	WowChanId      string `json:"wow_chan_id"`
	BotChanId      string `json:"bot_chan_id"`
	RoleId         string `json:"role_id"`
	TriggerCommand string `json:"trigger_command"`
	RoleBoost      string `json:"role_boost"`
	PendaRole      string `json:"penda_role"`
	PendaGoldRole  string `json:"penda_gold_role"`
	SWCommand      string `json:"sw_command"`
	SWCRoleId      string `json:"swc_role_id"`
}

func ReadConfig(configFile string) error {
	log.Println("Reading config file...")

	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Println("Error reading file IO: ", err.Error())
		return err
	}

	log.Println("JSON file : ", string(file))

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Println("Error reading JSON file: ", err.Error())
		return err
	}

	WowLogChanId = config.WowLogChanId
	WowChanId = config.WowChanId
	BotChanId = config.BotChanId
	RoleId = config.RoleId
	TriggerCommand = config.TriggerCommand
	RoleBoost = config.RoleBoost
	PendaRole = config.PendaRole
	PendaGoldRole = config.PendaGoldRole
	SWCommand = config.SWCommand
	SWCRoleId = config.SWCRoleId

	return nil
}
