#include "go_in_c.h"
void (*table[1]) = {Start};


//gcc -shared -pthread -o x.dll scratch.c go_in_c.a -lWinMM -lntdll -lWS2_32