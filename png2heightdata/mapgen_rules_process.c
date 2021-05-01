
typedef struct {
    int size;
    int height;
    int height_difference;
    int image_pixel_w, image_pixel_h;
} mapgen_panel_container;

typedef struct {
    int r, g, b;
    void (*func)(mapgen_panel_container*, int, int);
} mapgen_rule;

char hexchar(int n){
    return n + (n<10 ? '0':'a'-10);
}

void format_panel(char *target, mapgen_panel_container *cont, int adress, int data_digit){
    
    int point=0, value;
    
    value = cont[adress].height;
    for(int j=0; j<data_digit; j++){
        int n = value%10;
        target[point++] = hexchar(n);
        value /= 10;
    }
}

void panel_container_insert(mapgen_panel_container *cont, int x, int y, int image){
    
    int ad = y*cont->image_pixel_w+x;
    
    cont[ad].height=image;
    cont[ad].size++;
}

bool panel_equal(mapgen_panel_container *cont, int ax, int ay, int bx, int by){
    
    if(bx < 0 || bx >= cont->image_pixel_w || by < 0 || by >= cont->image_pixel_h) return 0;
    
    bool result = 1;
    for(int i=0; i<3; i++)
        result = result && (image_data[ay][ax*4+i] == image_data[by][bx*4+i]);
    return result;
}

int panel_equal_9(mapgen_panel_container *cont, int x, int y){
    int result = 0;
    for(int iy=-1; iy<=1; iy++){
        for(int ix=-1; ix<=1; ix++){
            result += ((int)panel_equal(cont, x, y, x+ix, y+iy)) << ((iy+1)*3+(ix+1));
        }
    }
    return result;
}

int panel_container_insert_9(
    mapgen_panel_container *cont, int x, int y, int depth, bool collision,
    int ts, int tc,   int te,
    int cs, int base, int ce,
    int bs, int bc,   int be,
    int pbe, int pbs,
    int pte, int pts){

    int equals = panel_equal_9(cont, x, y);

    //                0,0|0,1|1,0
    int ts_panel[] = {ts , cs, tc}; int ts_exists = (equals>>0)%2;
    int tc_panel[] = {tc ,pte,pts}; int tc_exists = (equals>>1)%2;
    int te_panel[] = {te , ce, tc}; int te_exists = (equals>>2)%2;
    int cs_panel[] = {cs ,pbs,pts}; int cs_exists = (equals>>3)%2;
    int ce_panel[] = {ce ,pbe,pte}; int ce_exists = (equals>>5)%2;
    int bs_panel[] = {bs , bc, cs}; int bs_exists = (equals>>6)%2;
    int bc_panel[] = {bc ,pbe,pbs}; int bc_exists = (equals>>7)%2;
    int be_panel[] = {be , bc, ce}; int be_exists = (equals>>8)%2;

    if(!ts_exists && !(tc_exists || cs_exists) )
        panel_container_insert(cont, x-1, y-1, ts_panel[cs_exists*2+tc_exists]);
    if(!tc_exists)
        panel_container_insert(cont, x  , y-1, tc_panel[te_exists*2+ts_exists]);
    if(!te_exists && !(tc_exists || ce_exists) )
        panel_container_insert(cont, x+1, y-1, te_panel[ce_exists*2+tc_exists]);
    if(!cs_exists)
        panel_container_insert(cont, x-1, y  , cs_panel[bs_exists*2+ts_exists]);
    if(!ce_exists)
        panel_container_insert(cont, x+1, y  , ce_panel[be_exists*2+te_exists]);
    if(!bs_exists && !(cs_exists || bc_exists) )
        panel_container_insert(cont, x-1, y+1, bs_panel[bc_exists*2+cs_exists]);
    if(!bc_exists)
        panel_container_insert(cont, x  , y+1, bc_panel[bs_exists*2+be_exists]);
    if(!be_exists && !(ce_exists || bc_exists) )
        panel_container_insert(cont, x+1, y+1, be_panel[bc_exists*2+ce_exists]);

    panel_container_insert(cont, x, y, base);

}


// quoted from : https://it-ojisan.tokyo/c-str-reverse/
void reverse(char* str){
	int size = strlen(str);
	int i,j;
	char tmp = {0};
	
	for(i = 0, j = size - 1; i < size / 2; i++, j--){
		tmp = str[i];
		str[i] = str[j];
		str[j] = tmp;
	}
	return;	
}