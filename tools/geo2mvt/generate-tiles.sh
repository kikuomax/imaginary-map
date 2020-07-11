#!/bin/bash

ZOOM=4
INPUT_JSON=../../../covid19-research/backend/render/papers.json
OUTPUT_BASE=../../docs/tiles
OUTPUT_SUFFIX=-papers.pbf.gz
EDGE_SIZE=$((2 ** $ZOOM))

for x in $(seq 0 $(($EDGE_SIZE-1))); do
	for y in $(seq 0 $(($EDGE_SIZE - 1))); do
		echo $x $y
		OUTPUT_DIR=$OUTPUT_BASE/$ZOOM/$x
		mkdir -p $OUTPUT_DIR
		OUTPUT_JSON=$OUTPUT_DIR/$y$OUTPUT_SUFFIX
		./geo2mvt -z $ZOOM -x $x -y $y $INPUT_JSON $OUTPUT_JSON
	done
done
