#!/usr/bin/env python3
"""
Optimized approach: Use sets to track code appearances
"""

import gzip
import sys

def main():
    # Track which codes appear in which files using sets
    codes_by_file = [set(), set(), set()]

    files = ['couponbase1.gz', 'couponbase2.gz', 'couponbase3.gz']

    # Process each file
    for file_idx, file_path in enumerate(files):
        print(f"Processing {file_path}...", file=sys.stderr)
        line_count = 0

        with gzip.open(file_path, 'rt', encoding='utf-8', errors='ignore') as f:
            for line in f:
                line_count += 1
                if line_count % 5000000 == 0:
                    print(f"  {line_count:,} lines, {len(codes_by_file[file_idx]):,} valid codes so far", file=sys.stderr)

                code = line.strip()
                if 8 <= len(code) <= 10:
                    codes_by_file[file_idx].add(code)

        print(f"Completed {file_path}: {len(codes_by_file[file_idx]):,} unique valid codes", file=sys.stderr)

    # Find codes appearing in at least 2 files using set operations (much faster!)
    print("\nFinding codes in 2+ files using set intersections...", file=sys.stderr)

    # Codes in files 1 AND 2
    in_1_and_2 = codes_by_file[0] & codes_by_file[1]
    # Codes in files 1 AND 3
    in_1_and_3 = codes_by_file[0] & codes_by_file[2]
    # Codes in files 2 AND 3
    in_2_and_3 = codes_by_file[1] & codes_by_file[2]

    # Union of all intersections gives us codes in at least 2 files
    valid_codes_set = in_1_and_2 | in_1_and_3 | in_2_and_3

    print(f"Codes in files 1&2: {len(in_1_and_2):,}", file=sys.stderr)
    print(f"Codes in files 1&3: {len(in_1_and_3):,}", file=sys.stderr)
    print(f"Codes in files 2&3: {len(in_2_and_3):,}", file=sys.stderr)

    # Sort and output
    valid_codes = sorted(valid_codes_set)

    print(f"\nTotal valid codes (in 2+ files): {len(valid_codes):,}", file=sys.stderr)
    print(f"Writing to valid_coupons.txt...", file=sys.stderr)

    with open('valid_coupons.txt', 'w') as out:
        for code in valid_codes:
            out.write(code + '\n')

    print("Done!", file=sys.stderr)

if __name__ == '__main__':
    main()
