package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/slcjordan/poc"
	"github.com/slcjordan/poc/logger"
)

type StartGameEncoding interface {
	EncodeStartGame(poc.StartGame) ([]byte, error)
	DecodeStartGame([]byte) (poc.StartGame, error)
}

type StartGame struct {
	Encoding StartGameEncoding
	Command  poc.StartGameCaller
}

func (s StartGame) CallBytes(ctx context.Context, b []byte) ([]byte, error) {
	game, err := s.Encoding.DecodeStartGame(b)
	if err != nil {
		return nil, poc.Error{Actual: fmt.Errorf("could not decode request: %w", err), Category: poc.MalformedError}
	}
	game, err = s.Command.CallStartGame(ctx, game)
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

type PerformMoveEncoding interface {
	EncodePerformMove(poc.PerformMove) ([]byte, error)
	DecodePerformMove([]byte) (poc.PerformMove, error)
}

type PerformMove struct {
	Encoding PerformMoveEncoding
	Command  poc.PerformMoveCaller
}

func (p PerformMove) CallBytes(ctx context.Context, b []byte) ([]byte, error) {
	move, err := p.Encoding.DecodePerformMove(b)
	if err != nil {
		return nil, poc.Error{Actual: fmt.Errorf("could not decode request: %w", err), Category: poc.MalformedError}
	}
	move, err = p.Command.CallPerformMove(ctx, move)
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

type ListGamesEncoding interface {
	EncodeListGames(poc.ListGames) ([]byte, error)
	DecodeListGames([]byte) (poc.ListGames, error)
}

type ListGames struct {
	Encoding ListGamesEncoding
	Command  poc.ListGamesCaller
}

func (l ListGames) CallBytes(ctx context.Context, b []byte) ([]byte, error) {
	list, err := l.Encoding.DecodeListGames(b)
	if err != nil {
		return nil, poc.Error{Actual: fmt.Errorf("could not decode request: %w", err), Category: poc.MalformedError}
	}
	list, err = l.Command.CallListGames(ctx, list)
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
