#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <linux/limits.h>
#include <string.h>
#include <stdint.h>
#include <limits.h>
#include <errno.h>

#define BUFF_SIZE 8

/**
 * Device filee path (e.g. "/sys/class/backlight/intel_backlight/brightness")
 */
static char devfile_path[PATH_MAX];

void set_devfile_path(const char *new_devfile_path) {
    strcpy(devfile_path, new_devfile_path);
}

/**
 * This function reads value from devfile_path, and stores to fetch_buff.
 * On success, it returns `0`.
 * On error, it returns non-zero;
 */
int read_devfile(char *fetch_buff, size_t buff_size) {
    if (access(devfile_path, F_OK) == -1) {
        fprintf(stderr, "the device file not exists. (%s)\n", devfile_path);
        return -1;
    }

    FILE *devfp = fopen(devfile_path, "r");
    if (devfp == NULL) {
        fprintf(stderr, "Couldn't read dev file (%s)\n", devfile_path);
        fclose(devfp);
        return -1;
    }
    char *got = fgets(fetch_buff, buff_size, devfp);
    if (got == NULL) {
        fprintf(stderr, "Couldn't read dev file (%s)\n", devfile_path);
        fclose(devfp);
        return -1;
    }
    fclose(devfp);
    return 0;
}

int main()
{
    // I've not decided about whether device file be constant value or CLI argument
    set_devfile_path("/sys/class/backlight/intel_backlight/brightness");

    char buff[BUFF_SIZE];
    int ret = read_devfile(buff, BUFF_SIZE);
    if (ret != 0) {
        exit(1);
    }
    printf("%s\n", buff);
    exit(0);
}
