#include <stdio.h>
#include <unistd.h>
#include <stdbool.h>

int main() {
	while (true) {
		printf("Hello enclave-cc!\n");
		sleep(2);
	}

	return true;
}
