
#define PANEL_VOID 0x00

#include "png_process.c"
#include "mapgen_rules_process.c"
int get_height(mapgen_panel_container*, int, int);
#include "mapgen_rules.c"

FILE *map_file;


// MAPDATAを作成/保存する。
int mapgen_process();

// メイン関数。
int main(int argc, char** argv){

    if(argc != 7){
        printf("Error : png2heightdata : Invalid arguments.\n");
        return 0;
    }
    
    image_file = fopen(argv[1], "r");
    map_file   = fopen(argv[2], "w+");

    int height_difference = (int)argv[3][0]-(int)('0');

    int image_pixel_w = atoi(argv[4]);
    int image_pixel_h = atoi(argv[5]);
    int data_digit = atoi(argv[6]);
    
    if(image_file == NULL || map_file == NULL){
        printf("Error : png2heightdata : Failed to open file.\n");
        return 0;
    }

    // mapdata.pngを読み込む。
    printf("png2heightdata : reading files...\n");
    read_process(image_pixel_w, image_pixel_h);
    printf("png2heightdata : processing...\n");
    mapgen_process(height_difference, image_pixel_w, image_pixel_h, data_digit);
    printf("png2heightdata : freeing resources...\n");
    free_process(image_pixel_w, image_pixel_h);

    fclose(map_file);
}

int get_height(mapgen_panel_container* cont, int ix, int iy){
    int r = image_data[iy][ix*4+0];
    int g = image_data[iy][ix*4+1];
    int b = image_data[iy][ix*4+2];
    return ((255-g)+b+r+1)*cont->height_difference;
}

int mapgen_process(int height_difference, int image_pixel_w, int image_pixel_h, int data_digit){

    int r, g, b;
    mapgen_panel_container* cont;
    cont = (mapgen_panel_container*)malloc(sizeof(mapgen_panel_container)*image_pixel_w*image_pixel_h);

    for (int i=0; i<image_pixel_w*image_pixel_h; i++){
        cont[i].size = 0;
        cont[i].height = 0;
        cont[i].height_difference = height_difference;
        cont[i].image_pixel_w = image_pixel_w;
        cont[i].image_pixel_h = image_pixel_h;
    }

    for (int iy=0; iy<image_pixel_h; iy++){
        for (int ix=0; ix<image_pixel_w; ix++){

            // image_data からrgb値を取得。
            r = image_data[iy][ix*4+0];
            g = image_data[iy][ix*4+1];
            b = image_data[iy][ix*4+2];

            // 全てのrulesを検索し、適合したらそれに応じたパネル。しなければrule_otherwiseで定義したパネル。
            for(int ir=0; ir<=RULES_NUM; ir++){
                if(
                    ir != RULES_NUM  &&
                    r == rules[ir].r &&
                    g == rules[ir].g &&
                    b == rules[ir].b
                    ){
                    rules[ir].func(cont, ix, iy);
                    break;
                }

                if(ir == RULES_NUM){
                    rule_otherwise(cont, ix, iy, r, g, b);
                    break;
                }
            }
        }
    }

    // 書き込む。
    char* code = (char*)malloc(sizeof(char)*data_digit);

    for (int iy=0; iy<image_pixel_h; iy++){
        for (int ix=0; ix<image_pixel_w; ix++){
            format_panel(code, cont, iy*(image_pixel_w+0.1666)+ix, data_digit);
            reverse(code);
            fprintf(map_file, "%s,", code);
        }
        fprintf(map_file, "\n");
    }

    free(code);

    free(cont);
    
}

