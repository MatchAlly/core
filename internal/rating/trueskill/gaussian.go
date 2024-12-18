package trueskill

import (
	"math"
)

type GaussianDistribution struct {
	Mean              float64
	StandardDeviation float64
	Precision         float64
	PrecisionMean     float64
	Variance          float64
}

func newGaussianDistribution(mean, standardDeviation float64) *GaussianDistribution {
	g := &GaussianDistribution{
		Mean:              mean,
		StandardDeviation: standardDeviation,
	}
	g.updateInternalValues()
	return g
}

func fromPrecisionMean(precisionMean, precision float64) *GaussianDistribution {
	g := &GaussianDistribution{
		Precision:     precision,
		PrecisionMean: precisionMean,
	}
	g.Variance = 1.0 / precision
	g.StandardDeviation = math.Sqrt(g.Variance)
	g.Mean = g.PrecisionMean / g.Precision
	return g
}

func (g *GaussianDistribution) updateInternalValues() {
	g.Variance = g.StandardDeviation * g.StandardDeviation
	g.Precision = 1.0 / g.Variance
	g.PrecisionMean = g.Precision * g.Mean
}

func (g *GaussianDistribution) NormalizationConstant() float64 {
	return 1.0 / (math.Sqrt(2*math.Pi) * g.StandardDeviation)
}

func (g *GaussianDistribution) Clone() *GaussianDistribution {
	return &GaussianDistribution{
		Mean:              g.Mean,
		StandardDeviation: g.StandardDeviation,
		Variance:          g.Variance,
		Precision:         g.Precision,
		PrecisionMean:     g.PrecisionMean,
	}
}

func (left GaussianDistribution) Mul(right GaussianDistribution) *GaussianDistribution {
	return fromPrecisionMean(left.PrecisionMean+right.PrecisionMean, left.Precision+right.Precision)
}

func absoluteDifference(left, right GaussianDistribution) float64 {
	return math.Max(
		math.Abs(left.PrecisionMean-right.PrecisionMean),
		math.Sqrt(math.Abs(left.Precision-right.Precision)))
}

func (left GaussianDistribution) Sub(right GaussianDistribution) float64 {
	return absoluteDifference(left, right)
}

func (left GaussianDistribution) Div(right GaussianDistribution) *GaussianDistribution {
	return fromPrecisionMean(left.PrecisionMean-right.PrecisionMean, left.Precision-right.Precision)
}

func (g *GaussianDistribution) SimpleAt(x float64) float64 {
	return g.At(x, 0, 1)
}

func (g *GaussianDistribution) At(x, mean, standardDeviation float64) float64 {
	multiplier := 1.0 / (standardDeviation * math.Sqrt(2*math.Pi)) // See http://mathworld.wolfram.com/NormalDistribution.html
	expPart := math.Exp(-(math.Pow(x-mean, 2.0) / (2 * (standardDeviation * standardDeviation))))
	return multiplier * expPart
}

func (g *GaussianDistribution) SimpleCumulativeTo(x float64) float64 {
	return g.CumulativeTo(x, 0, 1)
}

func (g *GaussianDistribution) CumulativeTo(x, mean, standardDeviation float64) float64 {
	invsqrt2 := -0.707106781186547524400844362104
	return 0.5 * errorFunctionCumulativeTo(invsqrt2*x)
}

func square(x float64) float64 {
	return x * x
}

func logProductNormalization(left, right GaussianDistribution) float64 {
	if left.Precision == 0 || right.Precision == 0 {
		return 0
	}
	varianceSum := left.Variance + right.Variance
	meanDifference := left.Mean - right.Mean
	logSqrt2Pi := math.Log(math.Sqrt(2 * math.Pi))
	return -logSqrt2Pi - (math.Log(varianceSum) / 2.0) - (square(meanDifference) / (2.0 * varianceSum))
}

func logRatioNormalization(numerator, denominator GaussianDistribution) float64 {
	if numerator.Precision == 0 || denominator.Precision == 0 {
		return 0
	}
	varianceDifference := denominator.Variance - numerator.Variance
	meanDifference := numerator.Mean - denominator.Mean
	logSqrt2Pi := math.Log(math.Sqrt(2 * math.Pi))
	return math.Log(denominator.Variance) + logSqrt2Pi - math.Log(varianceDifference)/2.0 +
		square(meanDifference)/(2*varianceDifference)
}

func errorFunctionCumulativeTo(x float64) float64 {
	z := math.Abs(x)
	t := 2.0 / (2.0 + z)
	ty := 4*t - 2

	coefficients := []float64{
		-1.3026537197817094, 6.4196979235649026e-1,
		1.9476473204185836e-2, -9.561514786808631e-3, -9.46595344482036e-4,
		3.66839497852761e-4, 4.2523324806907e-5, -2.0278578112534e-5,
		-1.624290004647e-6, 1.303655835580e-6, 1.5626441722e-8, -8.5238095915e-8,
		6.529054439e-9, 5.059343495e-9, -9.91364156e-10, -2.27365122e-10,
		9.6467911e-11, 2.394038e-12, -6.886027e-12, 8.94487e-13, 3.13092e-13,
		-1.12708e-13, 3.81e-16, 7.106e-15, -1.523e-15, -9.4e-17, 1.21e-16, -2.8e-17,
	}

	ncof := len(coefficients)
	d := 0.0
	dd := 0.0

	for j := ncof - 1; j > 0; j-- {
		tmp := d
		d = ty*d - dd + coefficients[j]
		dd = tmp
	}

	ans := t * math.Exp(-z*z+0.5*(coefficients[0]+ty*d)-dd)
	if x >= 0.0 {
		return ans
	} else {
		return 2.0 - ans
	}
}
