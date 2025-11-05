# Promo Code Validation

## Problem
Find valid promo codes that:
- Are 8-10 characters long
- Appear in at least 2 of 3 files (~600MB each, ~313M total lines)

## Solution

```bash
python3 process_coupons.py > valid_coupons.txt
```

**Expected runtime:** ~3 hours on M3 Pro (observed ~5M lines/min processing)
**Memory required:** ~15-18GB peak (stores 313M unique codes across 3 sets)

## How it works
1. Streams through each gzipped file line-by-line
2. Stores valid-length codes (8-10 chars) in memory using Python sets (~100M unique codes per file)
3. Performs set intersections to find codes in 2+ files
4. Outputs sorted results

## Results
8 valid promo codes found (see `valid_coupons.txt`)

## Files
- `process_coupons.py` - Main script (optimized with set operations)
- `couponbase1.gz`, `couponbase2.gz`, `couponbase3.gz` - Input data
- `valid_coupons.txt` - Output results
