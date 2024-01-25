# Fiber testdata

## Blocks
- Install [ethdo](https://github.com/wealdtech/ethdo)
- Command:
```bash
ethdo block info --blockid $SLOT_NUMBER --ssz | xxd -r -p > block.bin.ssz
```

Or run the `backfill.sh` script in this directory:
```bash
./backfill.sh $SLOT_NUMBER $NUM_SAMPLES
```
