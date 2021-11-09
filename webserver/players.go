package webserver

import (
	"les-randoms/riotinterface"
	"les-randoms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlePlayersRoute(c *gin.Context) {
	session := getSession(c)

	data := playersData{}

	setupNavData(&data.LayoutData.NavData, session)

	selectedItemName := setupSubnavData(&data.LayoutData.SubnavData, c, "Player Analyser", []string{"Last Game"})

	setupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = selectedItemName

	match, err := riotinterface.GetMatchFromId("EUW1_5540830822")
	if err != nil {
		utils.LogError(err.Error())
	} else {
		for _, p := range match.Info.Participants {
			if p.TeamId == 100 {
				data.LolGameReviewData.BlueTeam.Players = append(data.LolGameReviewData.BlueTeam.Players, lolPlayerGameReviewData{p.ChampionName})
			} else { // p.TeamId == 200
				data.LolGameReviewData.RedTeam.Players = append(data.LolGameReviewData.RedTeam.Players, lolPlayerGameReviewData{p.ChampionName})
			}
		}
	}

	c.HTML(http.StatusOK, "players.tmpl.html", data)
}
