#include <pthread.h>
#include <stdio.h>


int i;
pthread_mutex_t mutexmex = PTHREAD_MUTEX_INITIALIZER;

void* threadFunction1(){
	int x;
	for (x = 0 ; x < 1000000 ; x++){
		pthread_mutex_lock(&mutexmex);
		i ++;
		pthread_mutex_unlock(&mutexmex);
	}
	return;
}

void* threadFunction2(){
	int y;
	for (y =0 ; y < 1000001 ; y++){
		pthread_mutex_lock(&mutexmex);
		i --;
		pthread_mutex_unlock(&mutexmex);
	}
	return;
}

int main(){
	i = 0;
	pthread_t thread1;
	pthread_t thread2;
	pthread_create(&thread1, NULL, threadFunction1, NULL);
	pthread_create(&thread2, NULL, threadFunction2, NULL);

	pthread_join(thread1, NULL);
	pthread_join(thread2, NULL);

	printf("i: %d \n",i);
	return 0;
}	
