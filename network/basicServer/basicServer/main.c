//
//  main.c
//  basicServer
//
//  Created by daminyang on 18/01/2017.
//  Copyright © 2017 daminyang. All rights reserved.
//

#include <stdio.h>
#include <netdb.h>

int main(int argc, const char * argv[]) {
    printf("Hello, World!\n");
    
    struct hostent* hostentcontent=gethostbyname("www.baidu.com");
    
    printf("host name=%s\n",hostentcontent->h_name);
    
    
    return 0;
}
