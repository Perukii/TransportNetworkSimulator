
#include "png.h"
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>

// PNGファイルを読み込む。
// 読み込んだデータは image_dataに書き込まれる。
int read_process();

// メモリ解放を行う。
int free_process();

// グローバル変数

FILE          *image_file;
png_structp   image;
png_infop     info;
unsigned char **image_data;


int read_process(int image_pixel_w, int image_pixel_h){

    int bit_depth, color_type, interlace_type;
	int image_w_ret, image_h_ret;

	image = png_create_read_struct(PNG_LIBPNG_VER_STRING, NULL, NULL, NULL);
	info = png_create_info_struct(image);
	png_init_io(image, image_file);
	png_read_info(image, info);

	png_get_IHDR(image, info, &image_w_ret, &image_h_ret,
	                &bit_depth, &color_type, &interlace_type,
	                NULL, NULL);

	image_data = (png_bytepp)malloc(image_pixel_h * sizeof(png_bytep));

	for (int i=0; i<image_pixel_h; i++)
	    image_data[i] = (png_bytep)malloc(png_get_rowbytes(image, info));

	png_read_image(image, image_data);

    return 1;

}

int free_process(int image_pixel_w, int image_pixel_h){

	for (int i=0; i<image_pixel_h; i++)
        free(image_data[i]);

    free(image_data);

	png_destroy_read_struct(
	        &image, &info, (png_infopp)NULL);

    fclose(image_file);

    return 1;

}

