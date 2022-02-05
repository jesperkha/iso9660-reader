<br>
<div align="center">
  <h1 align="center">iso9666-reader</h1>
  <h3 align="center">Small CLI program for traversing a disk image and extracting files.</h3>
</div>
<br>

## About

ISO 9660 is the standard file system for CD-ROMs. It is also widely used on DVD and BD media and may as well be present on USB sticks or hard disks. This program allows you to traverse the file system with a simple command line interface. You can also extract files from the disk.

[ISO 9660 Wiki](https://wiki.osdev.org/ISO_9660)

<br>

<div align="center">
  <img src=".github/example.gif" alt="Example GIF" width="700">
</div>

<br>

## Installation and use

**Installation**

If you want to try it out for yourself you can [get the newest version.](https://github.com/jesperkha/iso9660-reader/releases/tag/v1.0.0)

Then run:

```console
$ ./isoread [file path]
```

<br>

**Commands**

- `exit`

  Exits the program

- `help`

  View a list of commands

- `cd`

  You can navigate around using the standard `cd` command followed by a directory name. If you leave the name empty you will go back to the root directory.

- `ls`

  To list the contents of the current directory use `ls`.

- `open`

  To view a file as raw text you can use the `open` command followed by the files name. This is only really useful for text files.

- `get`

  Using the `get` command followed by a filename extracts the file into the directory you ran the program from. It keeps the files name.
