#include <stdio.h>
#include "cairo_studio/cairo_studio.c"

int image_width, image_height;

int main(int argc, char** argv){
    if(argc != 4){
        printf("Error : heightdata_viewer : Invalid arguments.\n");
        return 0;
    }

    printf("heightdata_viewer : processing...\n");

    const char* heightdata_src = argv[1];
    int image_width = atoi(argv[2]);
    int image_height = atoi(argv[3]);

    FILE* heightdata_file = fopen(heightdata_src, "r");
    if(heightdata_file == NULL){
        printf("Error : heightdata_viewer : Failed to open file.\n");
        return 0;
    }

    int** heightdata;
    heightdata = (int**)malloc(sizeof(int*)*image_height);
    char* line = NULL;
    size_t len = 0;
    ssize_t read;
    const int data_segment_size = 5;
    
    for(int row = 0; row < image_height; row++) {
        if(getline(&line, &len, heightdata_file) == -1)break;
        heightdata[row] = (int*)malloc(sizeof(int)*image_width);

        for(int column = 0; column < image_width; column++){
            char segment[data_segment_size];
            for(int i=0; i<data_segment_size; i++){
                segment[i] = line[column*data_segment_size+i];
            }
            if(row == 0) printf("%s ", segment);
            heightdata[row][column] = atoi(segment);
        }
        
    }
    /*
    for(int column = 0; column < image_width; column++){
        printf("%d ", heightdata[1][column]);
    }

    */
    printf("finished\n");

    fclose(heightdata_file);
    free(heightdata);
    
}