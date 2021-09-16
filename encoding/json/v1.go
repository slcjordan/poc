package json

import (
	"encoding/json"

	"github.com/slcjordan/poc"
)

type v1Variant struct {
	MaxTimesThroughDeck int32 `json:"max_times_through_deck"`
}

type v1PositionedCard struct {
	Position []string `json:"position"`
	Suit     string   `json:"suit"`
	Index    string   `json:"index"`
}

type v1Board struct {
	Piles [13][]v1PositionedCard `json:"piles"`
	Score int32                  `json:"score"`
}

type v1Move struct {
	OldPileNum      int      `json:"old_pile_num"`
	OldPileIndex    int      `json:"old_pile_index"`
	OldPilePosition []string `json:"old_pile_position"`
	NewPileNum      int      `json:"new_pile_num"`
	NewPileIndex    int      `json:"new_pile_index"`
	NewPilePosition []string `json:"new_pile_position"`
}

type v1SavedGameDetail struct {
	GameID            int64      `json:"game_id"`
	Board             v1Board    `json:"board"`
	History           [][]v1Move `json:"history"`
	PossibleNextMoves [][]v1Move `json:"possible_moves"`
	Variant           v1Variant  `json:"variant"`
}

func v1LookupPosition(desc []string) poc.Position {
	var result poc.Position
	for _, d := range desc {
		result |= map[string]poc.Position{
			"face_up": poc.FaceUp,
		}[d]
	}
	return result
}

func toV1Moves(moves [][]poc.Move) [][]v1Move {
	result := make([][]v1Move, len(moves))
	lookupV1Move := map[poc.Position][]string{
		poc.FaceUp: {"face_up"},
	}
	for i := range moves {
		result[i] = make([]v1Move, len(moves[i]))
		for j := range moves[i] {
			result[i][j] = v1Move{
				OldPileNum:      moves[i][j].OldPileNum,
				OldPileIndex:    moves[i][j].OldPileIndex,
				OldPilePosition: lookupV1Move[moves[i][j].OldPilePosition],
				NewPileNum:      moves[i][j].NewPileNum,
				NewPileIndex:    moves[i][j].NewPileIndex,
				NewPilePosition: lookupV1Move[moves[i][j].NewPilePosition],
			}
		}
	}
	return result
}

func toV1Piles(cards [13][]poc.PositionedCard) [13][]v1PositionedCard {
	var result [13][]v1PositionedCard
	lookupV1Move := map[poc.Position][]string{
		poc.FaceUp: {"face_up"},
	}
	for idx := 0; idx < 13; idx++ {
		result[idx] = make([]v1PositionedCard, len(cards[idx]))
		for i := range cards[idx] {
			result[idx][i].Position = lookupV1Move[cards[idx][i].Position]
			result[idx][i].Suit = cards[idx][i].Card.Suit.String()
			result[idx][i].Index = cards[idx][i].Card.Index.String()
		}
	}
	return result
}

func toV1SavedGame(saved poc.SavedGameDetail) v1SavedGameDetail {
	return v1SavedGameDetail{
		GameID: saved.GameID,
		Board: v1Board{
			Piles: toV1Piles(saved.Board.Piles),
			Score: saved.Board.Score,
		},
		History:           toV1Moves(saved.History),
		PossibleNextMoves: toV1Moves(saved.PossibleNextMoves),
		Variant: v1Variant{
			MaxTimesThroughDeck: saved.Variant.MaxTimesThroughDeck,
		},
	}
}

// V1 json encoding.
type V1 struct{}

// DecodeStartGame unmarshals start game input.
func (v V1) DecodeStartGame(b []byte) (poc.StartGame, error) {
	var variant v1Variant
	err := json.Unmarshal(b, &variant)
	if err != nil {
		return poc.StartGame{}, err
	}
	return poc.StartGame{
		Input: poc.Variant{
			MaxTimesThroughDeck: variant.MaxTimesThroughDeck,
		},
	}, nil
}

// EncodeStartGame marshals start game result.
func (v V1) EncodeStartGame(started poc.StartGame) ([]byte, error) {
	result := toV1SavedGame(started.Result)
	return json.Marshal(result)
}

// DecodePerformMove unmarshals perform move input.
func (v V1) DecodePerformMove(b []byte) (poc.PerformMove, error) {
	var jsonMoves []v1Move
	err := json.Unmarshal(b, &jsonMoves)
	if err != nil {
		return poc.PerformMove{}, err
	}

	moves := make([]poc.Move, len(jsonMoves))
	for i, curr := range jsonMoves {
		moves[i].OldPileNum = curr.OldPileNum
		moves[i].OldPileIndex = curr.OldPileIndex
		moves[i].OldPilePosition = v1LookupPosition(curr.OldPilePosition)
		moves[i].NewPileNum = curr.NewPileNum
		moves[i].NewPileIndex = curr.NewPileIndex
		moves[i].NewPilePosition = v1LookupPosition(curr.NewPilePosition)
	}
	var result poc.PerformMove
	result.Input.Move = moves
	return result, nil
}

// EncodePerformMove marshals perform move result.
func (v V1) EncodePerformMove(performed poc.PerformMove) ([]byte, error) {
	result := toV1SavedGame(performed.Result)
	return json.Marshal(result)
}

// DecodeListGames unmarshals list games input.
func (v V1) DecodeListGames(b []byte) (poc.ListGames, error) {
	var result poc.ListGames
	return result, nil
}

// EncodeListGames marshals list games result.
func (v V1) EncodeListGames(list poc.ListGames) ([]byte, error) {
	result := make([]struct {
		GameID int64 `json:"game_id"`
		IsOver bool  `json:"is_over"`
		IsWon  bool  `json:"is_won"`
		Score  int32 `json:"score"`
	}, len(list.Result))
	for i := range list.Result {
		result[i].GameID = list.Result[i].GameID
		result[i].Score = list.Result[i].Score
	}
	return json.Marshal(result)
}
