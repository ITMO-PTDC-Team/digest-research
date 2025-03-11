package main
import "math"

func sinIntegratedQ(k , compression float64) float64 {
	return (math.Sin(math.Min(k, compression)*math.Pi/compression-math.Pi/2.0) + 1.0) / 2.0
}

func sinIntegratedLocation(q ,compression float64) float64 {
	return compression * (math.Asin(2.0*q-1.0) + math.Pi/2.0) / math.Pi
}


func expIntegratedQ(k , compression float64) float64{
    q := math.Min(k, compression) / compression

    c := 4.0
    return (math.Exp(c * q) - 1.0) / (math.Exp(c) - 1.0)
}

func expIntegratedLocation(q ,compression float64) float64 {
    c := 4.0
    return compression * math.Log(1.0 + (math.Exp(c) - 1.0) * q) / c
}

func powIntegratedQ(k, compression float64) float64 {
    q := math.Min(k, compression) / compression
    pow := 4.0
    return math.Pow(q, pow)
}

func powIntegratedLocation(q, compression float64) float64 {
    pow := 4.0
    return compression * math.Pow(q, 1.0/pow)
}

func sqrtSinIntegratedQ(k , compression float64) float64 {
	return math.Sqrt((math.Sin(math.Min(k, compression)*math.Pi/compression-math.Pi/2.0) + 1.0) / 2.0)
}

func sqrtSinIntegratedLocation(q ,compression float64) float64 {
	return compression * (math.Asin(2.0*q*q-1.0) + math.Pi/2.0) / math.Pi
}

var Tdgest_names =[]string{"tdgest_sinus", "tdgest_exp", "tdgest_pow","tdgest_sin_to_sqr"}
var Tdgest_functions =[]func(float64,float64)float64{sinIntegratedQ,sinIntegratedLocation,expIntegratedQ,expIntegratedLocation,powIntegratedQ,powIntegratedLocation,sqrtSinIntegratedQ,sqrtSinIntegratedLocation}
