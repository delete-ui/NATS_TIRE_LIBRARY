package types

import "time"

type SportType string

const (
	SportCSGO     SportType = "counter-strike"
	SportDOTA2    SportType = "dota2"
	SportLoL      SportType = "league-of-legends"
	SportValorant SportType = "valorant"
	SportRainbow6 SportType = "rainbow-six"
)

type Bookmaker string

const (
	BookmakerParivivsion Bookmaker = "parivision"
	BookmakerFonbet      Bookmaker = "fonbet"
	BookmakerOlimpBet    Bookmaker = "olimp-bet"
	BookmakerBetBoom     Bookmaker = "bet-boom"
	BookmakerWinline     Bookmaker = "winline"
)

type MarketType string

const (
	MarketMatchWinner  MarketType = "match-winner"
	MarketTotalMaps    MarketType = "total-maps"
	MarketMainHandicap MarketType = "main-handicap"
	MarketMainTotal    MarketType = "main-total"
	MarketPainting     MarketType = "painting"
)

type MatchBundle struct {
	CorrelationID   int
	TeamNames       []string
	BookmakerBundle map[Bookmaker]string
}

// TODO: Добавить информацию о лигах\турнирах
type MatchMonitoring struct {
	CorrelationID   int
	SportType       SportType
	TeamNames       []string
	BookmakerBundle map[Bookmaker]string //match url
	Bets            map[Bookmaker]Bet
	Timestamp       time.Time
}

type Bet struct {
	BetMarket MarketType
	Less      float64
	More      float64
}

type Fork struct {
	CorrelationID   int
	TeamNames       []string
	BookmakerBundle map[Bookmaker]string //match url
	Timestamp       time.Time
}
