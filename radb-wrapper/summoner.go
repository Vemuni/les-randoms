package radbwrapper

import (
	"les-randoms/database"
	"les-randoms/riotinterface"
	"les-randoms/utils"
	"time"
)

func riotSummonerToDBSummoner(summoner riotinterface.Summoner) database.Summoner {
	return database.Summoner_Construct(
		summoner.Id,
		0,
		summoner.AccountId,
		summoner.Puuid,
		summoner.Name,
		summoner.ProfileIconId,
		summoner.SummonerLevel,
		summoner.RevisionDate,
	)
}

func GetSummonerFromName(name string) (database.Summoner, error) {
	summoner, err := database.Summoner_SelectFirst("WHERE name=" + utils.Esc(name))
	if err == nil {
		if time.Now().Sub(summoner.LastUpdated).Hours() > 1 {
			riotSummoner, err := riotinterface.GetSummonerFromName(name)
			if err != nil { // In this case we return the last informations we have in the DB even if they are not the most recents possible
				return summoner, err
			}
			summoner = riotSummonerToDBSummoner(*riotSummoner)
			database.Summoner_Update(summoner)
		}
	} else if err.Error() == database.SummonerErrors.SummonerMissingInDB {
		riotSummoner, err := riotinterface.GetSummonerFromName(name)
		if err != nil {
			return database.Summoner{}, err
		}
		summoner = riotSummonerToDBSummoner(*riotSummoner)
		_, err = database.Summoner_Insert(summoner)
		if err != nil {
			return database.Summoner{}, err
		}
	} else {
		return database.Summoner{}, err
	}
	return summoner, nil
}
