
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

void rule_otherwise(RULE_ARG(cont, x, y), int r, int g, int b){
    panel_container_insert(cont, x, y, ((255-g)+b+r)*cont->height_difference);
}

