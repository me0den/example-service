package enum

import (
	"errors"
	"fmt"
	"strconv"
)

type BattleResult int

const (
	BattleResultLose BattleResult = -1
	BattleResultTie  BattleResult = 0
	BattleResultWin  BattleResult = 1
)

const (
	BattleResultLoseName = "lose"
	BattleResultTieName  = "tie"
	BattleResultWinName  = "win"
)

var (
	BattleResultName2Value = map[string]BattleResult{
		BattleResultLoseName: BattleResultLose,
		BattleResultTieName:  BattleResultTie,
		BattleResultWinName:  BattleResultWin,
	}
	BattleResultValue2Name = map[BattleResult]string{
		BattleResultLose: BattleResultLoseName,
		BattleResultTie:  BattleResultTieName,
		BattleResultWin:  BattleResultWinName,
	}
	BattleResultValues = map[BattleResult]struct{}{
		BattleResultLose: {},
		BattleResultTie:  {},
		BattleResultWin:  {},
	}
)

func (e BattleResult) String() string {
	return BattleResultValue2Name[e]
}

func (e *BattleResult) UnmarshalText(v []byte) error {
	var ok bool
	*e, ok = BattleResultName2Value[string(v)]
	if !ok {
		return fmt.Errorf("%s is not a valid BattleResult", string(v))
	}

	return nil
}

func (e BattleResult) MarshalText() ([]byte, error) {
	return []byte(BattleResultValue2Name[e]), nil
}

func (e *BattleResult) UnmarshalJSON(v []byte) error {
	str, err := strconv.Unquote(string(v))
	if err != nil {
		return errors.New("enums must be strings")
	}

	var ok bool
	*e, ok = BattleResultName2Value[str]
	if !ok {
		return fmt.Errorf("%s is not a valid BattleResult", str)
	}

	return nil
}

func (e BattleResult) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(BattleResultValue2Name[e])), nil
}

func (e *BattleResult) UnmarshalInt(v int) error {
	var ok bool
	_, ok = BattleResultValues[BattleResult(v)]
	if !ok {
		return fmt.Errorf("%d is not a valid BattleResult", v)
	}

	*e = BattleResult(v)
	return nil
}
