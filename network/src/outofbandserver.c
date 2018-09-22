#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <assert.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <errno.h>

#define BUF_SIZE 1024

int main(int argc, char const *argv[])
{
	if(argc<=2)
	{
		printf("usage :%s ip_address port_number\n",argv[0]);

		return 1;
	}

	const char* ip=argv[1];

	int port=atoi(argv[2]);

	struct sockaddr_in address;

	bzero(&address,sizeof(address));

	address.sin_family=AF_INET;

	inet_pton(AF_INET,ip,&address.sin_addr);

	address.sin_port=htons(port);

	int sock=socket(PF_INET,SOCK_STREAM,0);

	assert(sock>=0);

	//设置socket 选项，重用地址
	int reuse =1;
    //经过setsockopt的设置之后，即是sock处于TIME_WAIT状态，与之绑定的socket地址也可以立即被重用。
	setsockopt(sock,SOL_SOCKET,SO_REUSEADDR,&reuse,sizeof(reuse));

	int recvBuffer=0;

	int len=sizeof(recvBuffer);

	getsockopt(sock,SOL_SOCKET,SO_SNDBUF,&recvBuffer,(socklen_t*)&len);

	printf("recvBuffer=%d\n",recvBuffer);

	int ret=bind(sock,(struct sockaddr*)&address,sizeof(address));

	assert(ret!=-1);

	ret=listen(sock,5);

	assert(ret!=-1);

	struct sockaddr_in client;

	socklen_t client_addrlength=sizeof(client);

	printf("bind ip=%s,port=%d\n", ip,port);

	int connfd=accept(sock,(struct sockaddr*)&client,&client_addrlength);

	printf("someone got.....\n");

	if (connfd<0){

		printf("errno is %d\n",errno);

	}else{

		char buffer[BUF_SIZE];

		memset(buffer,'\0',BUF_SIZE);

		ret=recv(connfd,buffer,BUF_SIZE-1,0);

		printf("got %d bytes of normal data '%s'\n", ret,buffer);

		memset(buffer,'\0',BUF_SIZE);

		ret=recv(connfd,buffer,BUF_SIZE-1,MSG_OOB);

		printf("got %d bytes of oob data '%s'\n", ret,buffer);

		memset(buffer,'\0',BUF_SIZE);

		ret=recv(connfd,buffer,BUF_SIZE-1,0);

		printf("got %d bytes of normal data '%s'\n", ret,buffer);

		close(connfd);
	}
    
    close(sock);

	return 0;
}

