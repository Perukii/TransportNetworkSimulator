

PROJECT_NAME="Tohoku"

RESOURCE_DIR="$PWD/projects/$PROJECT_NAME"
IMAGE_FILE="$RESOURCE_DIR/image.png"
PROPERTY="$RESOURCE_DIR/property.sh"

DATA_DIR="$PWD/data/$PROJECT_NAME"
HEIGHT_DATA_FILE="$DATA_DIR/heightdata.txt"

mkdir -p $DATA_DIR
#echo "" > $HEIGHT_DATA_FILE

. $PROPERTY

(cd png2heightdata;. png2heightdata.sh $IMAGE_FILE $HEIGHT_DATA_FILE $HEIGHT_DIFFELENCE $IMAGE_WIDTH $IMAGE_HEIGHT)
(cd heightdata_viewer;. heightdata_viewer.sh $HEIGHT_DATA_FILE $IMAGE_WIDTH $IMAGE_HEIGHT)


