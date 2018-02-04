#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <linux/limits.h>
#include <string.h>

#define BUFF_SIZE 8

/**
 * Device filee path (e.g. "/sys/class/backlight/intel_backlight/brightness")
 */
static char dev_file_path[PATH_MAX];

void set_dev_file_path(const char *new_dev_file_path) {
    strcpy(dev_file_path, new_dev_file_path);
}

/**
 * This function reads value from dev_file_path, and stores to fetch_buff.
 * On success, it returns `0`.
 * On error, it returns non-zero;
 */
int read_dev_file(char *fetch_buff, size_t buff_size) {
    if (dev_file_path == NULL) {
        fprintf(stderr, "device file path is not set.\n");
        return -1;
    }
    if (access(dev_file_path, F_OK) == -1) {
        fprintf(stderr, "the device file not exists. (%s)\n", dev_file_path);
        return -1;
    }

    FILE *devfp = fopen(dev_file_path, "r");
    if (devfp == NULL) {
        fprintf(stderr, "Couldn't read dev file (%s)\n", dev_file_path);
        fclose(devfp);
        return -1;
    }
    char *got = fgets(fetch_buff, buff_size, devfp);
    if (got == NULL) {
        fprintf(stderr, "Couldn't read dev file (%s)\n", dev_file_path);
        fclose(devfp);
        return -1;
    }
    fclose(devfp);
    return 0;
}

int main(int argc, char const* argv[])
{
    // I've not decided about whether device file be constant value or CLI argument
    set_dev_file_path("/sys/class/backlight/intel_backlight/brightness");

    char buff[BUFF_SIZE];
    int ret = read_dev_file(buff, BUFF_SIZE);
    if (ret != 0) {
        exit(1);
    }
    printf("%s\n", buff);
    exit(0);
}