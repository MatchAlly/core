package rating

import (
	"math"
)

type MatchResult struct {
	OpponentRating    float64
	OpponentDeviation float64

	Result float64
}

func ApplyInactiveRatingPeriods(r Rating, inactivePeriods float64) Rating {
	r.Deviation = math.Sqrt(math.Pow(r.Deviation, 2) + math.Pow(r.Volatility, 2)*inactivePeriods)

	r = applyBounds(r)

	return r
}

func ApplyActiveRatingPeriod(r Rating, matchResults []MatchResult) Rating {
	if len(matchResults) == 0 {
		return ApplyInactiveRatingPeriods(r, 1.0)
	}

	variance := r.Deviation * r.Deviation
	invVarianceEstimate := 0.0

	// Compute the estimated variance of the player's rating based on game outcomes.
	// Compute the estimated improvement in rating, delta, by comparing the pre-period
	// rating to the performance rating based on game outcomes.
	deltaScale := 0.0
	invSqrPi := 0.10132118364233777

	for _, matchResult := range matchResults {
		oppDev := matchResult.OpponentDeviation
		oppRating := matchResult.OpponentRating

		g := 1.0 / math.Sqrt(1.0+3.0*oppDev*oppDev*invSqrPi)
		e := 1.0 / (1.0 + math.Exp(-g*(r.Value-oppRating)))

		invVarianceEstimate += g * g * e * (1.0 - e)
		deltaScale += g * (matchResult.Result - e)
	}

	varianceEstimate := 1.0 / invVarianceEstimate
	delta := varianceEstimate * deltaScale

	// Compute the new volatility
	a := math.Log(r.Volatility * r.Volatility)
	deltaSqr := delta * delta
	epsilon := 0.000001

	A := a
	B := deltaSqr - variance - varianceEstimate
	if B > 0.0 {
		B = math.Log(B)
	} else {
		B = a - tau
		for f(B, deltaSqr, variance, varianceEstimate, a) < 0.0 {
			B -= tau
		}
	}

	// Compute new volatility with numerical iteration using the Illinois algorithm
	// modification of the regula falsi method.
	fA := f(A, deltaSqr, variance, varianceEstimate, a)
	fB := f(B, deltaSqr, variance, varianceEstimate, a)
	for math.Abs(B-A) > epsilon {
		C := A + (A-B)*fA/(fB-fA)
		fC := f(C, deltaSqr, variance, varianceEstimate, a)

		if fC*fB < 0.0 {
			A = B
			fA = fB
		} else {
			fA *= 0.5
		}

		B = C
		fB = fC
	}

	newVolatility := math.Exp(A * 0.5)

	// Update the new rating deviation based on one period's worth of time elapsing
	newDeviation := math.Sqrt(variance + newVolatility*newVolatility)

	// Update the rating and rating deviation according to the match results
	newDeviation = 1.0 / math.Sqrt(1.0/(newDeviation*newDeviation)+invVarianceEstimate)
	newRating := r.Value + newDeviation*newDeviation*deltaScale

	r.Value = newRating
	r.Deviation = newDeviation
	r.Volatility = newVolatility

	r = applyBounds(r)

	return r
}

func applyBounds(r Rating) Rating {
	r.Value = math.Min(math.Max(r.Value, minRating), maxRating)
	r.Deviation = math.Min(math.Max(r.Deviation, minDeviation), maxDeviation)
	r.Volatility = math.Min(math.Max(r.Volatility, minVolatility), maxVolatility)

	return r
}

func f(x, deltaSqr, variance, varianceEstimate, a float64) float64 {
	eX := math.Exp(x)
	temp := variance + varianceEstimate + eX
	return eX*(deltaSqr-temp)/(2.0*temp*temp) - (x-a)/(tau*tau)
}
