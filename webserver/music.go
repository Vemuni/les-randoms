package webserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"les-randoms/database"
	discordbot "les-randoms/discord-bot/web-exec"
	webexec "les-randoms/discord-bot/web-exec"
	"les-randoms/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func handleMusicRoute(c *gin.Context) {
	session := getSession(c)

	if c.Param("subNavItem") == "" {
		c.Redirect(http.StatusFound, "/discord-bot/music/playing")
	}

	if isNotAuthentified(session) {
		redirectToAuth(c)
		return
	}

	if getAccessStatus(session, "/discord-bot") <= database.RightTypes.Forbidden {
		redirectToIndex(c)
		return
	}

	data := discordBotMusicData{}

	setupNavData(&data.LayoutData.NavData, session)

	selectedItemName := setupSubnavData(&data.LayoutData.SubnavData, c, "Music Bot", []string{"Playing"}, map[string]string{"Playing": "Playing"})

	setupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = selectedItemName

	switch data.LayoutData.SubnavData.SelectedSubnavItemIndex {
	case 0:
		if c.Param("order") != "" {
			err := handlePlayingOrder(c)
			if err != nil {
				utils.LogError("Happened while executing discord bot playing order " + utils.Esc(c.Param("order")) + " : " + err.Error())
			} else {
				utils.LogClassic("Discord bot playing order " + utils.Esc(c.Param("order")) + " successfully executed")
				return
			}
		}
		if setupPlayingData(&data) != nil {
			c.Redirect(http.StatusFound, "/discord-bot/music/playing")
		}
	}

	c.HTML(http.StatusFound, "discord-bot-music.tmpl.html", data)
}

func handlePlayingOrder(c *gin.Context) error {
	order := c.Param("order")
	if order == "resume" {
		return discordbot.ExecuteMusicResume()
	} else if order == "pause" {
		return discordbot.ExecuteMusicPause()
	} else {
		return errors.New("Unknown order")
	}
}

func setupPlayingData(data *discordBotMusicData) error {
	data.DiscordBotMusicPlayData.CurrentPlayStatus = discordbot.GetPlayStatus()
	return nil
}

func handlePlayingWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.LogError(err.Error())
		return
	}
	defer conn.Close()
	conn.SetCloseHandler(func(code int, text string) error {
		utils.LogDebug("Connection closed : " + fmt.Sprint(code) + " : " + text)
		return nil
	})
	for {
		//mt, msg, err := conn.ReadMessage()
		//if err != nil {
		//	utils.LogError(err.Error())
		//}
		//if mt != websocket.TextMessage {
		//	utils.LogError(err.Error())
		//}
		//var v message
		//json.Unmarshal(msg, &v)
		time.Sleep(time.Second)
		for !webexec.GetPlayStatus() {
			utils.LogDebug("tick2")
			time.Sleep(time.Second)
			currentMusicDuration := webexec.GetCurrentTime()
			jsonContent, err := json.Marshal(struct {
				PlayStatus  bool
				CurrentTime string
			}{
				PlayStatus:  webexec.GetPlayStatus(),
				CurrentTime: fmt.Sprintf("%02d:%02d", int(currentMusicDuration.Minutes()), int(currentMusicDuration.Seconds())%60),
			})
			if err != nil {
				utils.LogError(err.Error())
				return
			}
			utils.LogDebug(string(jsonContent))
			err = conn.WriteMessage(websocket.TextMessage, jsonContent)
			if err != nil {
				utils.LogError(err.Error())
				return
			}
		}
	}
}
