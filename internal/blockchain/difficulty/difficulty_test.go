package difficulty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Difficulty(t *testing.T) {
	t.Run("update", func(t *testing.T) {
		for name, test := range map[string]struct {
			providedDifficulty Difficulty
			providedNBReTry    uint
			expectedDifficulty Difficulty
		}{
			"with the mining reference value": {
				providedDifficulty: Difficulty(1),
				providedNBReTry:    nbTryWaiting,
				expectedDifficulty: Difficulty(1),
			},
			"with twice of the mining reference value": {
				providedDifficulty: Difficulty(1),
				providedNBReTry:    nbTryWaiting / 2,
				expectedDifficulty: Difficulty(2),
			},
			"with 50% of the mining reference value": {
				providedDifficulty: Difficulty(1),
				providedNBReTry:    nbTryWaiting * 2,
				expectedDifficulty: Difficulty(1),
			},
			"with twice of the mining reference value but current difficulty is 1": {
				providedDifficulty: Difficulty(1),
				providedNBReTry:    nbTryWaiting * 2,
				expectedDifficulty: Difficulty(1),
			},
		} {
			t.Run(name, func(t *testing.T) {
				test.providedDifficulty.Update(test.providedNBReTry)
				assert.Equal(t, test.expectedDifficulty, test.providedDifficulty)
			})
		}
	})

	t.Run("percent", func(t *testing.T) {
		providedDifficulty := Difficulty(1)

		for name, test := range map[string]struct {
			providedNBReTry    uint
			expectedDifficulty uint
		}{
			"with the mining reference value": {
				providedNBReTry:    nbTryWaiting,
				expectedDifficulty: 100,
			},
			"with 50% of the mining reference value": {
				providedNBReTry:    nbTryWaiting / 2,
				expectedDifficulty: 50,
			},
			"with twice of the mining reference value": {
				providedNBReTry:    nbTryWaiting * 2,
				expectedDifficulty: 200,
			},
			"with twice of the mining reference value but current difficulty is 1": {
				providedNBReTry:    0,
				expectedDifficulty: 0,
			},
		} {
			t.Run(name, func(t *testing.T) {
				percent := providedDifficulty.percent(test.providedNBReTry)
				assert.Equal(t, int(test.expectedDifficulty), int(percent))
			})
		}
	})

	t.Run("Int", func(t *testing.T) {
		difficulty := Difficulty(5)
		assert.Equal(t, 5, difficulty.Int())
	})

	t.Run("Save", func(t *testing.T) {
		difficulty := Difficulty(1)
		difficulty.Save(8)
		assert.Equal(t, Difficulty(8), difficulty)

		t.Run("must do nothing oif difficulty is less than one", func(t *testing.T) {
			difficulty.Save(0)
			assert.Equal(t, Difficulty(8), difficulty)
		})
	})
}
