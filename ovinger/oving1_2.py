from threading import Thread
i = 0

def threadFunction1():
	global i
	for x in range (0, 1000000):
		i+=1

def threadFunction2():
	global i
	for y in range (0, 1000000):
		i-=1

def main():
	global i

	thread1 = Thread(target = threadFunction1, args = (),)

	thread1.start()

	thread2 = Thread(target = threadFunction2, args = (),)

	thread2.start()
	
	thread1.join()

	thread2.join()

	print(i)

main()

		 	
