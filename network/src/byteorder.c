#include <stdio.h>

void byteorder()
{
	union
	{
       short value;
       char union_bytes[sizeof(short)];
	}test;
    //测试数据，一个高位是1，一个低位是2
	test.value=0x0102;
	if((test.union_bytes[0]==1)&&(test.union_bytes[1]==2))
	{
        printf("big endian\n");
	}
	else if((test.union_bytes[0]==2)&&(test.union_bytes[1]==1))
	{
		printf("little endian\n");
	}
	else
	{
		printf("unknow...\n");
	}
}

int main(int argc, char const *argv[])
{
	
	byteorder();
	return 0;
}

