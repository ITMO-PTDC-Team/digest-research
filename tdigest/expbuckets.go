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
    return NewWithCompression(100)
}

func (t *ExpBuckets) Add(x, w float64) {
    if w <= 0 || math.IsNaN(x) || math.IsNaN(w) {
        return
    }

    if x < t.Min {
        t.Min = x
    }
    if x > t.Max {
        t.Max = x
    }

    if t.Min == t.Max {
        t.Buckets[0].Mean = x
        t.Buckets[0].Weight += w
        t.Count += w
        return
    }

    scale := float64(t.Compression-1) / math.Log2(t.Max-t.Min+1)
    bucketIndex := t.Compression - 1 - int(math.Log2(t.Max-x+1)*scale)

    if bucketIndex < 0 {
        bucketIndex = 0
    } else if bucketIndex >= t.Compression {
        bucketIndex = t.Compression - 1
    }

    bucket := &t.Buckets[bucketIndex]
    if bucket.Weight == 0 {
        bucket.Mean = x
    } else {
        bucket.Mean += (x - bucket.Mean) * (w / (bucket.Weight + w))
    }
    bucket.Weight += w
    t.Count += w
}

func (t *ExpBuckets) Quantile(q float64) float64 {
    if q < 0 || q > 1 || len(t.Buckets) == 0 {
        return math.NaN()
    }
    sort.Slice(t.Buckets, func(i, j int) bool {
        return t.Buckets[i].Mean < t.Buckets[j].Mean
    })
    totalWeight := 0.0
    for _, b := range t.Buckets {
        totalWeight += b.Weight
    }
    targetWeight := q * totalWeight
    cumulativeWeight := 0.0
    for i := 0; i < len(t.Buckets); i++ {
        if cumulativeWeight+t.Buckets[i].Weight >= targetWeight {
            if i == 0 {
                return t.Buckets[i].Mean
            }
            prevWeight := cumulativeWeight
            prevMean := t.Buckets[i-1].Mean
            currWeight := cumulativeWeight + t.Buckets[i].Weight
            currMean := t.Buckets[i].Mean

            ratio := (targetWeight - prevWeight) / (currWeight - prevWeight)
            return prevMean + ratio*(currMean-prevMean)
        }
        cumulativeWeight += t.Buckets[i].Weight
    }

    return t.Buckets[len(t.Buckets)-1].Mean
}
