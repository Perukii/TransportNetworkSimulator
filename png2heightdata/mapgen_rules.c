
#define RULE_ARG(C,X,Y) mapgen_panel_container* C, int X, int Y
#define RULE_ARG_P RULE_ARG(cont, x, y)
#define RULES_NUM 1

void rule_255_255_255(RULE_ARG_P); // 海
void rule_otherwise  (RULE_ARG_P, int r, int g, int b); // 上記以外

mapgen_rule rules[RULES_NUM] = {
    {255, 255, 255, rule_255_255_255},
};

void rule_255_255_255(RULE_ARG(cont, x, y)){
    panel_container_insert(cont, x, y, 0x000);
}

int get_max(int a, int b){
    if(a<b) return b;
    else    return a;
}

void rule_otherwise(RULE_ARG(cont, x, y), int r, int g, int b){

    double bias = 30.0;
    int target = get_height(cont, x, y);

    if(x >= 1 && x < cont->image_pixel_w-1){
        int lf = get_height(cont, x-1, y);
        int rg = get_height(cont, x+1, y);
        if(target > lf && target > rg){
            if(get_max(abs(target-lf), abs(target-rg)) > bias){
                target = get_max(lf, rg);
            }
        }
    }
    if(y >= 1 && y < cont->image_pixel_h-1){
        int up = get_height(cont, x, y-1);
        int dw = get_height(cont, x, y+1);
        if(target > up && target > dw){
            if(get_max(abs(target-up), abs(target-dw)) > bias){
                target = get_max(up, dw);
            }
        }
    }
    panel_container_insert(cont, x, y, target);
}

