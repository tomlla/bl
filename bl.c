#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <linux/limits.h>
#include <string.h>

#define BUFF_SIZE 8

/**
 * Device filee path (e.g. "/sys/class/backlight/intel_backlight/brightness")
 */
static char devFilePath[PATH_MAX];

void setDevFilePath(const char *newDevFilePath) {
    strcpy(devFilePath, newDevFilePath);
}

/**
 * This function reads value from devFilePath, and stores to fetchBuff.
 * On success, it returns `0`.
 * On error, it returns non-zero;
 */
int readDevFile(char *fetchBuff, size_t buffSize) {
    if (devFilePath == NULL) {
        fprintf(stderr, "device file path is not set.\n");
        return -1;
    }
    if (access(devFilePath, F_OK) == -1) {
        fprintf(stderr, "the device file not exists. (%s)\n", devFilePath);
        return -1;
    }

    FILE *devfp = fopen(devFilePath, "r");
    if (devfp == NULL) {
        fprintf(stderr, "Couldn't read dev file (%s)\n", devFilePath);
        fclose(devfp);
        return -1;
    }
    char *got = fgets(fetchBuff, buffSize, devfp);
    if (got == NULL) {
        fprintf(stderr, "Couldn't read dev file (%s)\n", devFilePath);
        fclose(devfp);
        return -1;
    }
    fclose(devfp);
    return 0;
}

// TODO: still writing...
// unsigned int getBrightnessLevel() {
//     uintmax_t num = strtoumax(s, NULL, 10);
//     if (num == UINTMAX_MAX && errno == ERANGE)
//             /* Could not convert. */
// }

int main(int argc, char const* argv[])
{
    // I've not decided about whether device file be constant value or CLI argument
    setDevFilePath("/sys/class/backlight/intel_backlight/brightness");

    char buff[BUFF_SIZE];
    int ret = readDevFile(buff, BUFF_SIZE);
    if (ret != 0) {
        exit(1);
    }
    printf("%s\n", buff);
    exit(0);
}