#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MB 1024 * 1024

int main(int argc, char **argv) {
  int memoryLimit;
  sscanf(argv[1], "%d", &memoryLimit);

  void *mem = malloc(memoryLimit * MB);
  if(mem == NULL) {
    exit(1);
  }
  memset(mem, 0, memoryLimit * MB);
  printf("using %d MB of memory\n", memoryLimit);

  printf("press any character to exit...\n");
  char c = getchar();

  return 0;
}
