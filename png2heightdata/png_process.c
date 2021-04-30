
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
int           image_width, image_height;
FILE          *image_file;
png_structp   image;
png_infop     info;
unsigned char **image_data;


int read_process(){

    int bit_depth, color_type, interlace_type;

	image = png_create_read_struct(PNG_LIBPNG_VER_STRING, NULL, NULL, NULL);
	info = png_create_info_struct(image);
	png_init_io(image, image_file);
	png_read_info(image, info);

	png_get_IHDR(image, info, &image_width, &image_height,
	                &bit_depth, &color_type, &interlace_type,
	                NULL, NULL);

	image_data = (png_bytepp)malloc(image_height * sizeof(png_bytep));

	for (int i=0; i<image_height; i++)
	    image_data[i] = (png_bytep)malloc(png_get_rowbytes(image, info));

	png_read_image(image, image_data);

    return 1;

}

int free_process(){

	for (int i=0; i<image_height; i++)
        free(image_data[i]);

    free(image_data);

	png_destroy_read_struct(
	        &image, &info, (png_infopp)NULL);

    fclose(image_file);

    return 1;

}

