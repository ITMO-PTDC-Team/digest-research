package main

import (
    "math"
    "sort"
)

type ExpBuckets struct {
    Compression int
    Buckets     []Bucket
    Min         float64
    Max         float64
    Count       float64
}

type Bucket struct {
    Mean   float64
    Weight float64
}

func NewWithCompression(compression int) *ExpBuckets {
    return &ExpBuckets{
        Compression: compression,
        Buckets:     make([]Bucket, compression),
        Min:         math.MaxFloat64,
        Max:         -math.MaxFloat64,
        Count:       0,
    }
}

func New() *ExpBuckets {
    return NewWithCompression(100000)
}


func (t *ExpBuckets) Add(x, w float64) {
    if w <= 0 || math.IsNaN(x) || math.IsNaN(w) {
        return
    }

    if x < t.Min {
        t.Min = x
    }

    wasMaxUpdated := false
    if x > t.Max {
        t.Max = x
        wasMaxUpdated = true
    }

    rangeSize := t.Max - t.Min
    if rangeSize == 0 {
        rangeSize = 1 
    }

    normalizedX := (x - t.Min) / rangeSize
    bucketIndex := int(math.Log1p(normalizedX*(math.E-1)) * float64(t.Compression-1))

    if bucketIndex < 0 {
        bucketIndex = 0
    } else if bucketIndex >= t.Compression {
        bucketIndex = t.Compression - 1
    }

    bucket := &t.Buckets[bucketIndex]

    if bucket.Weight == 0 {
        bucket.Mean = x
    } else {
        bucket.Mean = (bucket.Mean*bucket.Weight + x*w) / (bucket.Weight + w)
    }
    bucket.Weight += w
    t.Count += w

    if !wasMaxUpdated && bucketIndex < t.Compression-1 {
        t.Max -= (t.Max - x) * 0.05 
    }
}


func (t *ExpBuckets) Quantile(q float64) float64 {
    if q < 0 || q > 1 || len(t.Buckets) == 0 || t.Count == 0 {
        return math.NaN()
    }

    validBuckets := make([]Bucket, 0, len(t.Buckets))
    for _, b := range t.Buckets {
        if b.Weight > 0 {
            validBuckets = append(validBuckets, b)
        }
    }

    if len(validBuckets) == 0 {
        return math.NaN()
    }

    sort.Slice(validBuckets, func(i, j int) bool {
        return validBuckets[i].Mean < validBuckets[j].Mean
    })

    totalWeight := 0.0
    for _, b := range validBuckets {
        totalWeight += b.Weight
    }

    targetWeight := q * totalWeight
    cumulativeWeight := 0.0

    for i := 0; i < len(validBuckets); i++ {
        if cumulativeWeight+validBuckets[i].Weight >= targetWeight {
            if i == 0 {
                return validBuckets[i].Mean
            }

            prevWeight := cumulativeWeight
            prevMean := validBuckets[i-1].Mean
            currWeight := cumulativeWeight + validBuckets[i].Weight
            currMean := validBuckets[i].Mean

            if currWeight == prevWeight {
                return currMean
            }

            ratio := (targetWeight - prevWeight) / (currWeight - prevWeight)
            return prevMean + ratio*(currMean-prevMean)
        }
        cumulativeWeight += validBuckets[i].Weight
    }

    return validBuckets[len(validBuckets)-1].Mean
}
