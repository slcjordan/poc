package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/slcjordan/poc"
	"github.com/slcjordan/poc/logger"
)

// StartGameEncoding may deserialize a startGame input and serialize a
// startGame result.
type StartGameEncoding interface {
	EncodeStartGame(poc.StartGame) ([]byte, error)
	DecodeStartGame([]byte) (poc.StartGame, error)
}

// StartGame command turns a startGame command (usually a pipeline) into a []byte command.
type StartGame struct {
	Encoding StartGameEncoding
	Pipeline poc.StartGameCaller
}

// CallBytes forwards parsed bytes to the StartGame command.
func (s StartGame) CallBytes(ctx context.Context, b []byte) ([]byte, error) {
	game, err := s.Encoding.DecodeStartGame(b)
	if err != nil {
		return nil, poc.Error{Actual: fmt.Errorf("could not decode request: %w", err), Category: poc.MalformedError}
	}
	game, err = s.Pipeline.CallStartGame(ctx, game)
	if err != nil {
		return nil, err
	}
	result, err := s.Encoding.EncodeStartGame(game)
	if err != nil {
		logger.Errorf(ctx, "could not encode start game response %#v: %s", game, err)
		return nil, poc.Error{Actual: errors.New("could not encode response"), Category: poc.UnknownError}
	}
	return result, nil
}

// PerformMoveEncoding may deserialize a performMove input and serialize a
// performMove result.
type PerformMoveEncoding interface {
	EncodePerformMove(poc.PerformMove) ([]byte, error)
	DecodePerformMove([]byte) (poc.PerformMove, error)
}

// PerformMove command turns a performMove command (usually a pipeline) into a []byte command.
type PerformMove struct {
	Encoding PerformMoveEncoding
	Pipeline poc.PerformMoveCaller
}

// CallBytes forwards parsed bytes to the PerformMove command.
func (p PerformMove) CallBytes(ctx context.Context, b []byte) ([]byte, error) {
	move, err := p.Encoding.DecodePerformMove(b)
	if err != nil {
		return nil, poc.Error{Actual: fmt.Errorf("could not decode request: %w", err), Category: poc.MalformedError}
	}
	move, err = p.Pipeline.CallPerformMove(ctx, move)
	if err != nil {
		return nil, err
	}
	result, err := p.Encoding.EncodePerformMove(move)
	if err != nil {
		logger.Errorf(ctx, "could not encode perform move response %#v: %s", move, err)
		return nil, poc.Error{Actual: errors.New("could not encode response"), Category: poc.UnknownError}
	}
	return result, nil
}

// ListGamesEncoding may deserialize a listGames input and serialize a
// listGames result.
type ListGamesEncoding interface {
	EncodeListGames(poc.ListGames) ([]byte, error)
	DecodeListGames([]byte) (poc.ListGames, error)
}

//Pipelinemes command turns a listGames command (usually a pipeline) into a []byte command.
type ListGames struct {
	Encoding ListGamesEncoding
	Pipeline poc.ListGamesCaller
}

// CallBytes forwards parsed bytes to the ListGames command.
func (l ListGames) CallBytes(ctx context.Context, b []byte) ([]byte, error) {
	list, err := l.Encoding.DecodeListGames(b)
	if err != nil {
		return nil, poc.Error{Actual: fmt.Errorf("could not decode request: %w", err), Category: poc.MalformedError}
	}
	list, err = l.Pipeline.CallListGames(ctx, list)
	if err != nil {
		return nil, err
	}
	result, err := l.Encoding.EncodeListGames(list)
	if err != nil {
		logger.Errorf(ctx, "could not encode list games response %#v: %s", list, err)
		return nil, poc.Error{Actual: errors.New("could not encode response"), Category: poc.UnknownError}
	}
	return result, nil
}
