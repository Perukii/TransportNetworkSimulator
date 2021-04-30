#define PANEL_DATA_DIGIT 4 
#define PANEL_DATA_SIZE 50

#define PANEL_VOID 0x00

#include "png_process.c"
#include "mapgen_rules_process.c"
#include "mapgen_rules.c"

FILE *map_file;

// MAPDATAを作成/保存する。
int mapgen_process();

// メイン関数。
int main(int argc, char** argv){

    if(argc != 4){
        printf("Error : png2heightdata : Invalid arguments.\n");
        return 0;
    }
    

    image_file = fopen(argv[1], "r");
    map_file   = fopen(argv[2], "w+");

    int height_difference = (int)argv[3][0]-(int)('0');
    
    if(image_file == NULL || map_file == NULL){
        printf("Error : png2heightdata : Failed to open file.\n");
        return 0;
    }

    // mapdata.pngを読み込む。
    printf("png2heightdata : reading files...\n");
    read_process();
    printf("png2heightdata : processing...\n");
    mapgen_process(height_difference);
    printf("png2heightdata : freeing resources...\n");
    free_process();

    fclose(map_file);
}


int mapgen_process(int height_difference){

    int r, g, b;
    mapgen_panel_container* cont;
    cont = (mapgen_panel_container*)malloc(sizeof(mapgen_panel_container)*image_width*image_height);

    for (int i=0; i<image_width*image_height; i++){
        cont[i].size = 0;
        cont[i].height = 0;
        cont[i].height_difference = height_difference;
    }

    for (int iy=0; iy<image_height; iy++){
        for (int ix=0; ix<image_width; ix++){

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
    char code[PANEL_DATA_SIZE];

    for (int iy=0; iy<image_height; iy++){
        for (int ix=0; ix<image_width; ix++){
            format_panel(code, cont, iy*image_width+ix);
            reverse(code);
            fprintf(map_file, "%s,", code);
        }
        fprintf(map_file, "\n");
    }

    free(cont);
    
}

