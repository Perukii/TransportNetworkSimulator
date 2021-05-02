

PROJECT_NAME="Tohoku"

RESOURCE_DIR="$PWD/projects/$PROJECT_NAME"
IMAGE_FILE="$RESOURCE_DIR/image.png"
CITY_LIST_FILE="$RESOURCE_DIR/citylist.txt"
PROPERTY="$RESOURCE_DIR/property.sh"

DATA_DIR="$PWD/data/$PROJECT_NAME"
HEIGHT_DATA_FILE="$DATA_DIR/heightdata.txt"
HEIGHT_DATA_DIGIT="5"
CITY_DATA_FILE="$DATA_DIR/citydata.txt"
PATH_DATA_FILE="$DATA_DIR/pathdata.txt"

COMMON_ARG="$HEIGHT_DATA_FILE $IMAGE_WIDTH \
$IMAGE_HEIGHT $HEIGHT_DATA_DIGIT $CITY_DATA_FILE \
$LONGITUDE_START $LONGITUDE_END $LATITUDE_START $LATITUDE_END \
$PATH_DATA_FILE"

mkdir -p $DATA_DIR

. $PROPERTY

#(cd png2heightdata;. png2heightdata.sh $IMAGE_FILE $HEIGHT_DATA_FILE \
#$HEIGHT_DIFFELENCE $IMAGE_WIDTH $IMAGE_HEIGHT $HEIGHT_DATA_DIGIT)

#(cd citylist2data;. citylist2data.sh $CITY_LIST_FILE $CITY_DATA_FILE)
(cd simpath;. simpath.sh $COMMON_ARG)

#(cd heightdata_viewer;. heightdata_viewer.sh )



