package webserver

import (
	"les-randoms/riotinterface"
	"les-randoms/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handlePlayersRoute(c *gin.Context) {
	session := getSession(c)

	if c.Param("subNavItem") == "" {
		c.Redirect(http.StatusFound, "/players/LastGame")
	}

	data := playersData{}

	setupNavData(&data.LayoutData.NavData, session)

	selectedItemName := setupSubnavData(&data.LayoutData.SubnavData, c, "Player Analyser", []string{"LastGame"}, map[string]string{"LastGame": "Last Game"})

	setupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = selectedItemName

	data.LastGameParameters.SummonerName = c.Param("param1")
	err := setupLolGameReviewData(&data)
	if err != nil {
		c.Redirect(http.StatusFound, "/players/LastGame")
	}

	c.HTML(http.StatusOK, "players.tmpl.html", data)
}

func setupLolGameReviewData(data *playersData) error {
	if data.LastGameParameters.SummonerName != "" {
		puuid, err := riotinterface.GetPuuidFromSummonerName(data.LastGameParameters.SummonerName)
		if err != nil {
			utils.LogError(err.Error())
			return err
		}
		matchId, err := riotinterface.GetLastMatchIdFromPuuid(puuid)
		if err != nil {
			utils.LogError(err.Error())
			return err
		}
		match, err := riotinterface.GetMatchFromId(matchId)
		if err != nil {
			utils.LogError(err.Error())
			return err
		}
		data.LolGameReviewData.GameDuration = strconv.Itoa(int(match.Info.GameDuration)/60) + ":" + strconv.Itoa(int(match.Info.GameDuration)%60)
		data.LolGameReviewData.GameMode = riotinterface.ParseGameMode(match.Info.GameMode)

		for _, p := range match.Info.Participants {
			var kda string
			if p.Deaths > 0 {
				tmp := (float32(p.Kills) + float32(p.Assists)) / float32(p.Deaths)
				kda = strconv.Itoa(int(tmp)) + "." + strconv.Itoa(int(tmp*100)%100)
			} else {
				kda = "Perfect"
			}
			trinket := riotinterface.GetItemsFromInt(p.Item6).Image.Full
			items := make([]string, 0)
			items = append(items, riotinterface.GetItemsFromInt(p.Item0).Image.Full)
			items = append(items, riotinterface.GetItemsFromInt(p.Item1).Image.Full)
			items = append(items, riotinterface.GetItemsFromInt(p.Item2).Image.Full)
			items = append(items, riotinterface.GetItemsFromInt(p.Item3).Image.Full)
			items = append(items, riotinterface.GetItemsFromInt(p.Item4).Image.Full)
			items = append(items, riotinterface.GetItemsFromInt(p.Item5).Image.Full)
			player := lolPlayerGameReviewData{
				"",
				p.SummonerName,
				p.ChampionName,
				riotinterface.GetSummonerSpellImageNameByKey(p.Summoner1Id),
				riotinterface.GetSummonerSpellImageNameByKey(p.Summoner2Id),
				p.TotalMinionsKilled,
				strconv.Itoa(p.GoldEarned/1000) + "." + strconv.Itoa((p.GoldEarned%1000)/100),
				p.Kills,
				p.Deaths,
				p.Assists,
				kda,
				p.WardsPlaced,
				p.WardsKilled,
				p.VisionWardsBoughtInGame,
				p.VisionScore,
				trinket,
				items,
			}
			player.Version, _ = riotinterface.GetLastVersionFromGameVersion(match.Info.GameVersion)
			if p.TeamId == 100 {
				data.LolGameReviewData.BlueTeam.Players = append(data.LolGameReviewData.BlueTeam.Players, player)
			} else { // p.TeamId == 200
				data.LolGameReviewData.RedTeam.Players = append(data.LolGameReviewData.RedTeam.Players, player)
			}
		}
	}
	return nil
}
