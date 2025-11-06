# Promo Code Validation

## Problem
Find valid promo codes that:
- Are 8-10 characters long
- Appear in at least 2 of 3 files (~600MB each, ~313M total lines)

## Valid Coupons

The database is already pre-populated with these 8 valid coupons:

`BIRTHDAY`, `BUYGETON`, `FIFTYOFF`, `FREEZAAA`, `GNULINUX`, `HAPPYHRS`, `OVER9000`, `SIXTYOFF`

## Reproducing the Processing (Optional)

To regenerate the coupon list from scratch:

### 1. Download Coupon Files

```bash
# Download files (~2GB gzipped, ~600MB each uncompressed)
curl -O https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase{1,2,3}.gz
```

### 2. Run Processing Script

```bash
python3 process_coupons.py > valid_coupons.txt
```

**⚠️ WARNING:** ~3 hours runtime, 15-18GB RAM required

**Performance:** ~5M lines/min throughput on an Apple Mac M3 Pro

## How It Works

1. Streams through each gzipped file line-by-line
2. Stores valid-length codes (8-10 chars) in memory using Python sets (~100M unique codes per file)
3. Performs set intersections to find codes appearing in 2+ files
4. Outputs sorted results

## Files

- `process_coupons.py` - Main script (optimized with set operations)
- `couponbase1.gz`, `couponbase2.gz`, `couponbase3.gz` - Input data (download separately)
- `valid_coupons.txt` - Output results (8 valid codes)
