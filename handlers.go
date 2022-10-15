/*
	Copyright (C) 2021-2022  The YNOproject Developers

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

var (
	errLengthMismatch = errors.New("command length mismatch")
)

func (h *Hub) handleIdent(msg []string, sender *HubClient) (err error) {
	if len(msg) != 1 {
		return errLengthMismatch
	}

	sender.valid = true
	sender.sendMsg("ident") // tell client they're valid
	h.handleValidClient(sender)

	return nil
}

func (h *Hub) handleM(msg []string, sender *HubClient) (err error) {
	if len(msg) != 3 {
		return errLengthMismatch
	}
	// check if the coordinates are valid
	x, errconv := strconv.Atoi(msg[1])
	if errconv != nil || x < 0 {
		return errconv
	}
	y, errconv := strconv.Atoi(msg[2])
	if errconv != nil || y < 0 {
		return errconv
	}
	sender.x = x
	sender.y = y

	if msg[0] == "m" {
		if sender.syncCoords {
			checkHubConditions(h, sender, "coords", "")
		}
		h.broadcast(sender, sender, "m", sender.sClient.id, msg[1:]) // user %id% moved to x y
	} else {
		checkHubConditions(h, sender, "teleport", "")
	}

	return nil
}

func (h *Hub) handleF(msg []string, sender *HubClient) (err error) {
	if len(msg) != 2 {
		return errLengthMismatch
	}
	// check if direction is valid
	facing, errconv := strconv.Atoi(msg[1])
	if errconv != nil || facing < 0 || facing > 3 {
		return errconv
	}
	sender.facing = facing
	h.broadcast(sender, "f", sender.sClient.id, msg[1]) // user %id% facing changed to f

	return nil
}

func (h *Hub) handleSpd(msg []string, sender *HubClient) (err error) {
	if len(msg) != 2 {
		return errLengthMismatch
	}
	spd, errconv := strconv.Atoi(msg[1])
	if errconv != nil {
		return errconv
	}
	if spd < 0 || spd > 10 { // something's not right
		return errconv
	}
	sender.spd = spd
	h.broadcast(sender, "spd", sender.sClient.id, msg[1])

	return nil
}

func (h *Hub) handleSpr(msg []string, sender *HubClient) (err error) {
	if len(msg) != 3 {
		return errLengthMismatch
	}
	if !isValidSprite(msg[1]) {
		return err
	}
	if config.gameName == "2kki" && (!strings.Contains(msg[1], "syujinkou") && !strings.Contains(msg[1], "effect") && !strings.Contains(msg[1], "yukihitsuji_game") && !strings.Contains(msg[1], "zenmaigaharaten_kisekae") && !strings.Contains(msg[1], "主人公") && !strings.Contains(msg[1], "#null")) || strings.Contains(msg[1], "zenmaigaharaten_kisekae") && h.roomId != 176 {
		return err
	}
	index, errconv := strconv.Atoi(msg[2])
	if errconv != nil || index < 0 {
		return errconv
	}
	sender.sClient.spriteName = msg[1]
	sender.sClient.spriteIndex = index
	h.broadcast(sender, "spr", sender.sClient.id, msg[1:])

	return nil
}

func (h *Hub) handleFl(msg []string, sender *HubClient) (err error) {
	if len(msg) != 6 {
		return errLengthMismatch
	}
	red, errconv := strconv.Atoi(msg[1])
	if errconv != nil || red < 0 || red > 255 {
		return errconv
	}
	green, errconv := strconv.Atoi(msg[2])
	if errconv != nil || green < 0 || green > 255 {
		return errconv
	}
	blue, errconv := strconv.Atoi(msg[3])
	if errconv != nil || blue < 0 || blue > 255 {
		return errconv
	}
	power, errconv := strconv.Atoi(msg[4])
	if errconv != nil || power < 0 {
		return errconv
	}
	frames, errconv := strconv.Atoi(msg[5])
	if errconv != nil || frames < 0 {
		return errconv
	}
	if msg[0] == "rfl" {
		sender.flash[0] = red
		sender.flash[1] = green
		sender.flash[2] = blue
		sender.flash[3] = power
		sender.flash[4] = frames
		sender.repeatingFlash = true
	}
	h.broadcast(sender, msg[0], sender.sClient.id, msg[1:])

	return nil
}

func (h *Hub) handleRrfl(sender *HubClient) (err error) {
	sender.repeatingFlash = false
	for i := 0; i < 5; i++ {
		sender.flash[i] = 0
	}
	h.broadcast(sender, "rrfl", sender.sClient.id)

	return nil
}

func (h *Hub) handleH(msg []string, sender *HubClient) (err error) {
	if len(msg) != 2 {
		return errLengthMismatch
	}
	hiddenBin, errconv := strconv.Atoi(msg[1])
	if errconv != nil || hiddenBin < 0 || hiddenBin > 1 {
		return errconv
	}
	sender.hidden = hiddenBin == 1
	h.broadcast(sender, msg[0], sender.sClient.id, msg[1])

	return nil
}

func (h *Hub) handleSys(msg []string, sender *HubClient) (err error) {
	if len(msg) != 2 {
		return errLengthMismatch
	}
	if !isValidSystem(msg[1], false) {
		return err
	}
	sender.sClient.systemName = msg[1]
	h.broadcast(sender, "sys", sender.sClient.id, msg[1])

	return nil
}

func (h *Hub) handleSe(msg []string, sender *HubClient) (err error) {
	if len(msg) != 5 {
		return errLengthMismatch
	}
	if !isValidSound(msg[1]) {
		return err
	}
	volume, errconv := strconv.Atoi(msg[2])
	if errconv != nil || volume < 0 || volume > 100 {
		return errconv
	}
	tempo, errconv := strconv.Atoi(msg[3])
	if errconv != nil || tempo < 10 || tempo > 400 {
		return errconv
	}
	balance, errconv := strconv.Atoi(msg[4])
	if errconv != nil || balance < 0 || balance > 100 {
		return errconv
	}
	h.broadcast(sender, "se", sender.sClient.id, msg[1:])

	return nil
}

func (h *Hub) handleP(msg []string, sender *HubClient) (err error) {
	isShow := msg[0] == "ap"
	msgLength := 18
	if isShow {
		msgLength = 20
	}
	if len(msg) != msgLength {
		return errLengthMismatch
	}

	if isShow {
		checkHubConditions(h, sender, "picture", msg[17])
		if !isValidPicName(msg[17]) {
			return err
		}
	}

	picId, errconv := strconv.Atoi(msg[1])
	if errconv != nil || picId < 1 {
		return errconv
	}

	positionX, errconv := strconv.Atoi(msg[2])
	if errconv != nil {
		return errconv
	}
	positionY, errconv := strconv.Atoi(msg[3])
	if errconv != nil {
		return errconv
	}
	mapX, errconv := strconv.Atoi(msg[4])
	if errconv != nil {
		return errconv
	}
	mapY, errconv := strconv.Atoi(msg[5])
	if errconv != nil {
		return errconv
	}
	panX, errconv := strconv.Atoi(msg[6])
	if errconv != nil {
		return errconv
	}
	panY, errconv := strconv.Atoi(msg[7])
	if errconv != nil {
		return errconv
	}

	magnify, errconv := strconv.Atoi(msg[8])
	if errconv != nil || magnify < 0 {
		return errconv
	}
	topTrans, errconv := strconv.Atoi(msg[9])
	if errconv != nil || topTrans < 0 {
		return errconv
	}
	bottomTrans, errconv := strconv.Atoi(msg[10])
	if errconv != nil || bottomTrans < 0 {
		return errconv
	}

	red, errconv := strconv.Atoi(msg[11])
	if errconv != nil || red < 0 || red > 200 {
		return errconv
	}
	green, errconv := strconv.Atoi(msg[12])
	if errconv != nil || green < 0 || green > 200 {
		return errconv
	}
	blue, errconv := strconv.Atoi(msg[13])
	if errconv != nil || blue < 0 || blue > 200 {
		return errconv
	}
	saturation, errconv := strconv.Atoi(msg[14])
	if errconv != nil || saturation < 0 || saturation > 200 {
		return errconv
	}

	effectMode, errconv := strconv.Atoi(msg[15])
	if errconv != nil || effectMode < 0 {
		return errconv
	}
	effectPower, errconv := strconv.Atoi(msg[16])
	if errconv != nil {
		return errconv
	}

	var pic *Picture
	if isShow {
		picName := msg[17]
		if picName == "" {
			return err
		}

		useTransparentColorBin, errconv := strconv.Atoi(msg[18])
		if errconv != nil || useTransparentColorBin < 0 || useTransparentColorBin > 1 {
			return errconv
		}

		fixedToMapBin, errconv := strconv.Atoi(msg[19])
		if errconv != nil || fixedToMapBin < 0 || fixedToMapBin > 1 {
			return errconv
		}

		pic = &Picture{
			name:                picName,
			useTransparentColor: useTransparentColorBin == 1,
			fixedToMap:          fixedToMapBin == 1,
		}

		if _, found := sender.pictures[picId]; found {
			rpErr := h.processMsg("rp"+delim+msg[1], sender)
			if rpErr != nil {
				return rpErr
			}
		}
	} else {
		if _, found := sender.pictures[picId]; found {
			duration, errconv := strconv.Atoi(msg[17])
			if errconv != nil || duration < 0 {
				return errconv
			}

			pic = sender.pictures[picId]
		} else {
			return nil
		}
	}

	pic.positionX = positionX
	pic.positionY = positionY
	pic.mapX = mapX
	pic.mapY = mapY
	pic.panX = panX
	pic.panY = panY
	pic.magnify = magnify
	pic.topTrans = topTrans
	pic.bottomTrans = bottomTrans
	pic.red = red
	pic.blue = blue
	pic.green = green
	pic.saturation = saturation
	pic.effectMode = effectMode
	pic.effectPower = effectPower

	h.broadcast(sender, msg[0], sender.sClient.id, msg[1:])

	sender.pictures[picId] = pic

	return nil
}

func (h *Hub) handleRp(msg []string, sender *HubClient) (err error) {
	if len(msg) != 2 {
		return errLengthMismatch
	}
	picId, errconv := strconv.Atoi(msg[1])
	if errconv != nil || picId < 1 {
		return errconv
	}
	h.broadcast(sender, "rp", sender.sClient.id, msg[1])
	delete(sender.pictures, picId)

	return nil
}

func (h *Hub) handleSay(msg []string, sender *HubClient) (err error) {
	if sender.sClient.muted {
		return nil
	}

	if len(msg) != 2 {
		return errLengthMismatch
	}
	msgContents := strings.TrimSpace(msg[1])
	if sender.sClient.name == "" || sender.sClient.systemName == "" || msgContents == "" || len(msgContents) > 150 {
		return err
	}
	h.broadcast(sender, "say", sender.sClient.id, msgContents)

	return nil
}

func (h *Hub) handleSs(msg []string, sender *HubClient) (err error) {
	if len(msg) != 3 {
		return errLengthMismatch
	}
	switchId, errconv := strconv.Atoi(msg[1])
	if errconv != nil {
		return errconv
	}
	valueBin, errconv := strconv.Atoi(msg[2])
	if errconv != nil || valueBin < 0 || valueBin > 1 {
		return errconv
	}
	var value bool
	if valueBin == 1 {
		value = true
	}
	sender.switchCache[switchId] = value
	if switchId == 1430 && config.gameName == "2kki" { // time trial mode
		if value {
			sender.sendMsg("sv", "88", "0") // time elapsed
		}
	} else {
		if len(sender.hub.minigameConfigs) > 0 {
			for m, minigame := range sender.hub.minigameConfigs {
				if minigame.SwitchId == switchId && minigame.SwitchValue == value && sender.minigameScores[m] < sender.varCache[minigame.VarId] {
					tryWritePlayerMinigameScore(sender.sClient.uuid, minigame.MinigameId, sender.varCache[minigame.VarId])
				}
			}
		}

		for _, c := range append(globalConditions, h.conditions...) {
			validVars := !c.VarTrigger
			if c.VarTrigger {
				if c.VarId > 0 {
					if value, ok := sender.varCache[c.VarId]; ok {
						if validVar, _ := c.checkVar(c.VarId, value); validVar {
							validVars = true
						}
					}
				} else if len(c.VarIds) > 0 {
					validVars = true
					for _, vId := range c.VarIds {
						if value, ok := sender.varCache[vId]; ok {
							if validVar, _ := c.checkVar(vId, value); !validVar {
								validVars = false
								break
							}
						} else {
							validVars = false
							break
						}
					}
				} else {
					validVars = true
				}
			}

			if validVars {
				if switchId == c.SwitchId {
					if valid, _ := c.checkSwitch(switchId, value); valid {
						if c.VarTrigger || (c.VarId == 0 && len(c.VarIds) == 0) {
							if !c.TimeTrial {
								if checkConditionCoords(c, sender) {
									success, err := tryWritePlayerTag(sender.sClient.uuid, c.ConditionId)
									if err != nil {
										return err
									}
									if success {
										sender.sendMsg("b")
									}
								}
							} else if config.gameName == "2kki" {
								sender.sendMsg("ss", "1430", "0")
							}
						} else {
							varId := c.VarId
							if len(c.VarIds) > 0 {
								varId = c.VarIds[0]
							}
							sender.sendMsg("sv", varId, "0")
						}
					}
				} else if len(c.SwitchIds) > 0 {
					if valid, s := c.checkSwitch(switchId, value); valid {
						if s == len(c.SwitchIds)-1 {
							if c.VarTrigger || (c.VarId == 0 && len(c.VarIds) == 0) {
								if !c.TimeTrial {
									if checkConditionCoords(c, sender) {
										success, err := tryWritePlayerTag(sender.sClient.uuid, c.ConditionId)
										if err != nil {
											return err
										}
										if success {
											sender.sendMsg("b")
										}
									}
								} else if config.gameName == "2kki" {
									sender.sendMsg("ss", "1430", "0")
								}
							} else {
								varId := c.VarId
								if len(c.VarIds) > 0 {
									varId = c.VarIds[0]
								}
								sender.sendMsg("sv", varId, "0")
							}
						} else {
							sender.sendMsg("ss", c.SwitchIds[s+1], "0")
						}
					}
				}
			}
		}
	}

	return nil
}

func (h *Hub) handleSv(msg []string, sender *HubClient) (err error) {
	if len(msg) != 3 {
		return errLengthMismatch
	}
	varId, errconv := strconv.Atoi(msg[1])
	if errconv != nil {
		return errconv
	}
	value, errconv := strconv.Atoi(msg[2])
	if errconv != nil {
		return errconv
	}
	sender.varCache[varId] = value

	conditions := append(globalConditions, h.conditions...)

	if varId == 88 && config.gameName == "2kki" {
		for _, c := range conditions {
			if c.TimeTrial && value < 3600 {
				if checkConditionCoords(c, sender) {
					success, err := tryWritePlayerTimeTrial(sender.sClient.uuid, h.roomId, value)
					if err != nil {
						return err
					}
					if success {
						sender.sendMsg("b")
					}
				}
			}
		}
	} else {
		if len(sender.hub.minigameConfigs) > 0 {
			for m, minigame := range sender.hub.minigameConfigs {
				if minigame.VarId == varId && sender.minigameScores[m] < value {
					if minigame.SwitchId > 0 {
						sender.sendMsg("ss", minigame.SwitchId, "0")
					} else {
						tryWritePlayerMinigameScore(sender.sClient.uuid, minigame.MinigameId, value)
					}
				}
			}
		}

		for _, c := range conditions {
			validSwitches := c.VarTrigger
			if !c.VarTrigger {
				if c.SwitchId > 0 {
					if value, ok := sender.switchCache[c.SwitchId]; ok {
						if validSwitch, _ := c.checkSwitch(c.SwitchId, value); validSwitch {
							validSwitches = true
						}
					}
				} else if len(c.SwitchIds) > 0 {
					validSwitches = true
					for _, sId := range c.SwitchIds {
						if value, ok := sender.switchCache[sId]; ok {
							if validSwitch, _ := c.checkSwitch(sId, value); !validSwitch {
								validSwitches = false
								break
							}
						} else {
							validSwitches = false
							break
						}
					}
				} else {
					validSwitches = true
				}
			}

			if validSwitches {
				if varId == c.VarId {
					if valid, _ := c.checkVar(varId, value); valid {
						if !c.VarTrigger || (c.SwitchId == 0 && len(c.SwitchIds) == 0) {
							if !c.TimeTrial {
								if checkConditionCoords(c, sender) {
									success, err := tryWritePlayerTag(sender.sClient.uuid, c.ConditionId)
									if err != nil {
										return err
									}
									if success {
										sender.sendMsg("b")
									}
								}
							} else if config.gameName == "2kki" {
								sender.sendMsg("ss", "1430", "0")
							}
						} else {
							switchId := c.SwitchId
							if len(c.SwitchIds) > 0 {
								switchId = c.SwitchIds[0]
							}
							sender.sendMsg("ss", switchId, "0")
						}
					}
				} else if len(c.VarIds) > 0 {
					if valid, v := c.checkVar(varId, value); valid {
						if v == len(c.VarIds)-1 {
							if !c.VarTrigger || (c.SwitchId == 0 && len(c.SwitchIds) == 0) {
								if !c.TimeTrial {
									if checkConditionCoords(c, sender) {
										success, err := tryWritePlayerTag(sender.sClient.uuid, c.ConditionId)
										if err != nil {
											return err
										}
										if success {
											sender.sendMsg("b")
										}
									}
								} else if config.gameName == "2kki" {
									sender.sendMsg("ss", "1430", "0")
								}
							} else {
								switchId := c.SwitchId
								if len(c.SwitchIds) > 0 {
									switchId = c.SwitchIds[0]
								}
								sender.sendMsg("ss", switchId, "0")
							}
						} else {
							sender.sendMsg("sv", c.VarIds[v+1], "0")
						}
					}
				}
			}
		}
	}

	return nil
}

func (h *Hub) handleSev(msg []string, sender *HubClient) (err error) {
	if len(msg) != 3 {
		return errLengthMismatch
	}
	actionBin, errconv := strconv.Atoi(msg[2])
	if errconv != nil || actionBin < 0 || actionBin > 1 {
		return errconv
	}
	triggerType := "event"
	if actionBin == 1 {
		triggerType = "eventAction"
	}
	checkHubConditions(h, sender, triggerType, msg[1])

	if sender.hub.roomId != currentEventVmMapId {
		return err
	}

	eventIdInt, err := strconv.Atoi(msg[1])
	if err != nil {
		return err
	}

	if currentEventVmEventId != eventIdInt {
		return err
	}

	exp, err := tryCompleteEventVm(currentEventPeriodId, sender.sClient.uuid, currentEventVmMapId, currentEventVmEventId)
	if err != nil {
		return err
	}
	if exp > -1 {
		sender.sClient.sendMsg("vm", exp)
	}

	return nil
}

// SESSION

func (s *Session) handleI(sender *SessionClient) (err error) {
	badgeSlotRows, badgeSlotCols := getPlayerBadgeSlotCounts(sender.name)
	playerInfo := PlayerInfo{
		Uuid:          sender.uuid,
		Name:          sender.name,
		Rank:          sender.rank,
		Badge:         sender.badge,
		BadgeSlotRows: badgeSlotRows,
		BadgeSlotCols: badgeSlotCols,
	}
	playerInfoJson, err := json.Marshal(playerInfo)
	if err != nil {
		return err
	}

	sender.sendMsg("i", playerInfoJson)

	return nil
}

func (s *Session) handleName(msg []string, sender *SessionClient) (err error) {
	if sender.hClient == nil {
		return err
	}

	if len(msg) != 2 {
		return errLengthMismatch
	}

	if sender.name != "" || !isOkString(msg[1]) || len(msg[1]) > 12 {
		return err
	}
	sender.name = msg[1]

	sender.hClient.hub.broadcast(sender.hClient, "name", sender.id, sender.name) // broadcast name change to hub if client is in one

	return nil
}

func (s *Session) handlePloc(msg []string, sender *SessionClient) (err error) {
	if sender.hClient == nil {
		return err
	}

	if len(msg) != 3 {
		return errLengthMismatch
	}

	if len(msg[1]) != 4 {
		return errors.New("invalid prev map ID")
	}

	sender.hClient.prevMapId = msg[1]
	sender.hClient.prevLocations = msg[2]
	checkHubConditions(sender.hClient.hub, sender.hClient, "prevMap", sender.hClient.prevMapId)

	return nil
}

func (s *Session) handleGSay(msg []string, sender *SessionClient) (err error) {
	if sender.hClient == nil {
		return err
	}

	if sender.muted {
		return nil
	}

	if len(msg) != 3 {
		return errLengthMismatch
	}
	msgContents := strings.TrimSpace(msg[1])
	if sender.name == "" || sender.systemName == "" {
		return errors.New("invalid client")
	} else if msgContents == "" || len(msgContents) > 150 {
		return errors.New("invalid message")
	}

	enableLocBin, errconv := strconv.Atoi(msg[2])
	if errconv != nil || enableLocBin < 0 || enableLocBin > 1 {
		return errconv
	}

	mapId := "0000"
	prevMapId := "0000"
	var prevLocations string
	x := -1
	y := -1

	mapId = sender.hClient.mapId
	prevMapId = sender.hClient.prevMapId
	prevLocations = sender.hClient.prevLocations
	x = sender.hClient.x
	y = sender.hClient.y

	session.broadcast("p", sender.uuid, sender.name, sender.systemName, sender.rank, sender.account, sender.badge)
	session.broadcast("gsay", sender.uuid, mapId, prevMapId, prevLocations, x, y, msgContents)

	return nil
}

func (s *Session) handlePSay(msg []string, sender *SessionClient) (err error) {
	if sender.muted {
		return nil
	}

	if len(msg) != 2 {
		return errLengthMismatch
	}
	msgContents := strings.TrimSpace(msg[1])
	if sender.name == "" || sender.systemName == "" {
		return errors.New("invalid client")
	} else if msgContents == "" || len(msgContents) > 150 {
		return errors.New("invalid message")
	}

	partyId, err := getPlayerPartyId(sender.uuid)
	if err != nil {
		return err
	}
	if partyId == 0 {
		return errors.New("player not in a party")
	}
	partyMemberUuids, err := getPartyMemberUuids(partyId)
	if err != nil {
		return err
	}
	for _, uuid := range partyMemberUuids {
		if client, ok := clients.Load(uuid); ok {
			client.(*SessionClient).sendMsg("psay", sender.uuid, msgContents)
		}
	}

	return nil
}

func (s *Session) handlePt(sender *SessionClient) (err error) {
	partyId, err := getPlayerPartyId(sender.uuid)
	if err != nil {
		return err
	}
	if partyId == 0 {
		return errors.New("player not in a party")
	}
	partyData, err := getPartyData(sender.uuid)
	if err != nil {
		return err
	}
	if sender.uuid != partyData.OwnerUuid {
		partyData.Pass = ""
	}
	partyDataJson, err := json.Marshal(partyData)
	if err != nil {
		return err
	}
	sender.sendMsg("pt", partyDataJson)

	return nil
}

func (s *Session) handleEp(sender *SessionClient) (err error) {
	period, err := getCurrentEventPeriodData()
	if err != nil {
		return err
	}
	periodJson, err := json.Marshal(period)
	if err != nil {
		return err
	}
	sender.sendMsg("ep", periodJson)

	return nil
}

func (s *Session) handleE(sender *SessionClient) (err error) {
	periodId, err := getCurrentEventPeriodId()
	if err != nil {
		return err
	}
	currentEventLocationsData, err := getCurrentPlayerEventLocationsData(periodId, sender.uuid)
	if err != nil {
		return err
	}
	var hasIncompleteEvent bool
	for _, currentEventLocation := range currentEventLocationsData {
		if !currentEventLocation.Complete {
			hasIncompleteEvent = true
			break
		}
	}
	if !hasIncompleteEvent && config.gameName == "2kki" {
		addPlayer2kkiEventLocation(-1, 2, 0, 0, sender.uuid)
		currentEventLocationsData, err = getCurrentPlayerEventLocationsData(periodId, sender.uuid)
		if err != nil {
			return err
		}
	}

	currentEventVmsData, err := getCurrentPlayerEventVmsData(periodId, sender.uuid)
	if err != nil {
		return err
	}

	eventsData := &EventsData{
		Locations: currentEventLocationsData,
		Vms:       currentEventVmsData,
	}

	eventsDataJson, err := json.Marshal(eventsData)
	if err != nil {
		return err
	}
	sender.sendMsg("e", eventsDataJson)

	return nil
}
