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
    if (devfile_path == NULL) {
        fprintf(stderr, "device file path is not set.\n");
        return -1;
    }
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

#define T(NAME, R_T, ERR_T) struct Tuple_#NAME { #R_T result; #ERR_T err; };
T(A, int, (char *))

// struct Tuple {
//     uint16_t value;
//     const char const *err;
// };

uint16_t to_int(const char *str_num) {
    char *end;
    const int base = 10;

    errno = 0;
    long int converted_num = strtol(str_num, &end, base);
    if (errno == ERANGE) {
    }
    if (converted_num == LONG_MAX) {
        fprintf(stderr, "Couldn't read dev file (%s)\n", devfile_path);
    }
}


int main(int argc, char const* argv[])
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