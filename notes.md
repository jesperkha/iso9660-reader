# Notes

### Test specs

folder/MYFILE.TXT
testvol

### Details

- Filenames are uppercase only + symbols
- Double endian is little followed by big
- Starts after 32KiB
- Sector is 2048 bytes

### File structure

32 kilobytes of boot info (not used)
n sectors of length 2048 bytes ends with a terminator of type 255
