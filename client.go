package goroyale

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const baseURL = "https://api.royaleapi.com"

// Client allows you to easily interact with RoyaleAPI.
type Client struct {
	Token string

	client http.Client
}

// New creates a new RoyaleAPI client.
func New(token string) (c *Client, err error) {
	if token == "" {
		panic("Client requires token for authorization with the API.")
	}
	c.Token = token
	c.client = http.Client{Timeout: (10 * time.Second)}

	return
}

// Args represents special args to pass in the request.
// The API supports args for Field Filter https://docs.royaleapi.com/#/field_filter
// and Pagination https://docs.royaleapi.com/#/pagination.
type Args struct {
	Keys    []string
	Exclude []string
	Max     int
	Page    int
}

func argQuery(args Args) (q url.Values) {
	if args.Keys != nil {
		q.Add("keys", strings.Join(args.Keys, ","))
	}

	if args.Exclude != nil {
		q.Add("exclude", strings.Join(args.Keys, ","))
	}

	if args.Max != 0 {
		q.Add("max", string(args.Max))
	}

	if args.Page != 0 {
		q.Add("page", string(args.Page))
	}

	return
}

func (c Client) get(path string, args Args) (bytes []byte, err error) {
	path = baseURL + path
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return
	}
	req.Header.Add("auth", c.Token)
	req.URL.RawQuery = argQuery(args).Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bytes, err = ioutil.ReadAll(resp.Body)
	return
}

// GetAPIVersion requests the current version of the API.
// https://docs.royaleapi.com/#/endpoints/version
func (c Client) GetAPIVersion() (ver string, err error) {
	bytes, err := c.get("/version", Args{})
	if err != nil {
		return
	}
	ver = string(bytes)
	return
}

// TODO: Add Client.GetConstants for /constants

// GetPlayer retrieves a player by their tag.
// https://docs.royaleapi.com/#/endpoints/player
func (c Client) GetPlayer(tag string, args Args) (player Player, err error) {
	var b []byte
	path := "/player/" + tag
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &player)
	}
	return
}

// GetPlayers works like GetPlayer but can return multiple players.
// The API asks that you don't include more than 7 tags in this request.
// https://docs.royaleapi.com/#/endpoints/player?id=multiple-players
func (c Client) GetPlayers(tags []string, args Args) (players []Player, err error) {
	var b []byte
	path := "/player/" + strings.Join(tags, ",")
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &players)
	}
	return
}

// GetPlayerBattles gets battles a player participated in.
// https://docs.royaleapi.com/#/endpoints/player_battles
func (c Client) GetPlayerBattles(tag string, args Args) (battles []Battle, err error) {
	var b []byte
	path := "/player/" + tag + "/battles"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &battles)
	}
	return
}

// GetPlayersBattles works like GetPlayerBattles but can return battles from multiple players.
// https://docs.royaleapi.com/#/endpoints/player_battles?id=multiple-tags
func (c Client) GetPlayersBattles(tags []string, args Args) (battles [][]Battle, err error) {
	var b []byte
	path := "/player/" + strings.Join(tags, ",") + "/battles"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &battles)
	}
	return
}

// GetPlayerChests gets a player's upcoming chests.
// https://docs.royaleapi.com/#/endpoints/player_chests
func (c Client) GetPlayerChests(tag string, args Args) (chests PlayerChests, err error) {
	var b []byte
	path := "/player/" + tag + "/chests"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &chests)
	}
	return
}

// GetPlayersChests works like GetPlayerChests but can return chests for multiple players.
// https://docs.royaleapi.com/#/endpoints/player_chests?id=multiple-players
func (c Client) GetPlayersChests(tags []string, args Args) (chests []PlayerChests, err error) {
	var b []byte
	path := "/player/" + strings.Join(tags, ",") + "/chests"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &chests)
	}
	return
}

// TODO: ClanSearch (https://docs.royaleapi.com/#/endpoints/clan_search)

// GetClan returns info about a specific clan.
// https://docs.royaleapi.com/#/endpoints/clan
func (c Client) GetClan(tag string, args Args) (clan Clan, err error) {
	var b []byte
	path := "/clan/" + tag
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &clan)
	}
	return
}

// GetClans works like GetClan but can return multiple clans.
// https://docs.royaleapi.com/#/endpoints/clan?id=multiple-clans
func (c Client) GetClans(tags []string, args Args) (clans []Clan, err error) {
	var b []byte
	path := "/clan/" + strings.Join(tags, ",")
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &clans)
	}
	return
}

// GetClanBattles returns battles played by people in the specified clan.
// https://docs.royaleapi.com/#/endpoints/clan_battles
func (c Client) GetClanBattles(tag string, args Args) (battles []Battle, err error) {
	var b []byte
	path := "/clan/" + tag + "/battles"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &battles)
	}
	return
}

// GetClanWar returns data about the current clan war.
// https://docs.royaleapi.com/#/endpoints/clan_war
func (c Client) GetClanWar(tag string, args Args) (war ClanWar, err error) {
	var b []byte
	path := "/clan/" + tag + "/war"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &war)
	}
	return
}

// GetClanWarLog returns data about past clan wars.
// https://docs.royaleapi.com/#/endpoints/clan_warlog
func (c Client) GetClanWarLog(tag string, args Args) (warlog []ClanWarLogEntry, err error) {
	var b []byte
	path := "/clan/" + tag + "/warlog"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &warlog)
	}
	return
}

// GetClanHistory returns a time series of member stats.
// This will only work with clans that have enabled stat tracking.
// https://docs.royaleapi.com/#/endpoints/clan_history
func (c Client) GetClanHistory(tag string, args Args) (history []ClanHistoryEntry, err error) {
	var b []byte
	path := "/clan/" + tag + "/history"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &history)
	}
	return
}

// GetClanWeeklyHistory works like GetClanHistory but returns weekly stats.
// https://docs.royaleapi.com/#/endpoints/clan_history_weekly
func (c Client) GetClanWeeklyHistory(tag string, args Args) (history []ClanHistoryEntry, err error) {
	var b []byte
	path := "/clan/" + tag + "/history/weekly"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &history)
	}
	return
}

// GetClanTracking returns basic data on whether a clan is tracked.
// https://docs.royaleapi.com/#/endpoints/clan_tracking
func (c Client) GetClanTracking(tag string, args Args) (tracking ClanTracking, err error) {
	var b []byte
	path := "/clan/" + tag + "/tracking"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &tracking)
	}
	return
}

// GetOpenTournaments returns a slice of open tournaments.
// https://docs.royaleapi.com/#/endpoints/tournaments_open
func (c Client) GetOpenTournaments(args Args) (tournaments []OpenTournament, err error) {
	var b []byte
	path := "/tournaments/open"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &tournaments)
	}
	return
}

// GetKnownTournaments returns a slice of tournaments people have searched for.
// https://docs.royaleapi.com/#/endpoints/tournaments_known
func (c Client) GetKnownTournaments(args Args) (tournaments []KnownTournament, err error) {
	var b []byte
	path := "/tournaments/known"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &tournaments)
	}
	return
}

// TODO: TournamentsSearch (https://docs.royaleapi.com/#/endpoints/tournaments_search)

// GetTournament returns the specified Tournament by tag.
// https://docs.royaleapi.com/#/endpoints/tournaments
func (c Client) GetTournament(tag string, args Args) (tournament Tournament, err error) {
	var b []byte
	path := "/tournaments/" + tag
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &tournament)
	}
	return
}

// GetTournaments works like GetTournament but can return multiple Tournaments.
// https://docs.royaleapi.com/#/endpoints/tournaments?id=multiple-tournaments
func (c Client) GetTournaments(tags []string, args Args) (tournaments []Tournament, err error) {
	var b []byte
	path := "/tournaments/" + strings.Join(tags, ",")
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &tournaments)
	}
	return
}

// GetTopClans returns the top 200 clans of a location/global leaderboard.
// https://docs.royaleapi.com/#/endpoints/top_clans
func (c Client) GetTopClans(location string, args Args) (topClans []TopClan, err error) {
	var b []byte
	path := "/top/clans/" + location
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &topClans)
	}
	return
}

// GetTopPlayers returns the top 200 players of a location/global leaderboard.
// https://docs.royaleapi.com/#/endpoints/top_players
func (c Client) GetTopPlayers(location string, args Args) (topPlayers []TopPlayer, err error) {
	var b []byte
	path := "/top/players/" + location
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &topPlayers)
	}
	return
}

// GetPopularClans returns stats on how often a clan's data has been requested from the API.
// https://docs.royaleapi.com/#/endpoints/popular_clans
func (c Client) GetPopularClans(args Args) (popularClans []PopularClan, err error) {
	var b []byte
	path := "/popular/clans"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &popularClans)
	}
	return
}

// GetPopularPlayers returns stats on how often a player's data has been request from the API.
// https://docs.royaleapi.com/#/endpoints/popular_players
func (c Client) GetPopularPlayers(args Args) (popularPlayers []PopularPlayer, err error) {
	var b []byte
	path := "/popular/players"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &popularPlayers)
	}
	return
}

// GetPopularTournaments returns stats on how often a tournament's data has been request from the API.
// https://docs.royaleapi.com/#/endpoints/popular_tournaments
func (c Client) GetPopularTournaments(args Args) (popularTournaments []PopularTournament, err error) {
	var b []byte
	path := "/popular/tournaments"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &popularTournaments)
	}
	return
}

// GetPopularDecks returns stats on how often a deck's data has been requested from the API.
// https://docs.royaleapi.com/#/endpoints/popular_decks
func (c Client) GetPopularDecks(args Args) (popularDecks []PopularDeck, err error) {
	var b []byte
	path := "/popular/decks"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &popularDecks)
	}
	return
}

// GetAPIKeyStats returns information about the currently authenticated token.
// https://docs.royaleapi.com/#/endpoints/auth_stats
func (c Client) GetAPIKeyStats(args Args) (keyStats APIKeyStats, err error) {
	var b []byte
	path := "/auth/stats"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &keyStats)
	}
	return
}

// GetEndpoints returns all the available endpoints for the API.
// It does not have any special incorporation with this wrapper and is simply included for completion's sake.
// https://docs.royaleapi.com/#/endpoints/endpoints
func (c Client) GetEndpoints(args Args) (endpoints []string, err error) {
	var b []byte
	path := "/endpoints"
	if b, err = c.get(path, args); err == nil {
		err = json.Unmarshal(b, &endpoints)
	}
	return
}
