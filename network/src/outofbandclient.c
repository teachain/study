#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <assert.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <errno.h>

int main(int argc, char const *argv[])
{
	if(argc<=2)
	{
		printf("Usage:%s ip_address port _number\n", argv[0]);
		return 1;
	}
	const char* ip=argv[1];

	int port=atoi(argv[2]);

	struct sockaddr_in server_address;

	bzero(&server_address,sizeof(server_address));

	server_address.sin_family=AF_INET;

	inet_pton(AF_INET,ip,&server_address.sin_addr);

	server_address.sin_port=htons(port);

	int sockfd=socket(PF_INET,SOCK_STREAM,0);

	assert(sockfd>=0);

	int ret =connect(sockfd,(struct sockaddr*)&server_address,sizeof(server_address));

	if(ret<0)
	{
		printf("connection failed ip=%s ;port=%d\n",ip,port);

		printf("connection failed ret=%d ;errno=%d\n",ret,errno);
	}
	else
	{
		const char * oob_data="abc";
		const char* normal_data="123";
		send(sockfd,normal_data,strlen(normal_data),0);
		send(sockfd,oob_data,strlen(oob_data),MSG_OOB);
		send(sockfd,normal_data,strlen(normal_data),0);
	}

	int code=0;

	scanf("%d",&code);

	printf("code=%d\n", code);

	close(sockfd);
	return 0;
}