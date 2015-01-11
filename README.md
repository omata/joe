# joe

![](http://i.imgur.com/y8g506n.png?1)

Joe generates `.gitignore` files from the command line.

## Features

- Written in known Python
- Easy to install
- Stupidly easy to use
- Supports att Github-supported `.gitignore` files
- Warns if target directory isn't a git repo

## Usage

### Basic usage

If there is an existing `.gitignore` file, appends the new file to it. Otherwise, makes a new `.gitignore` file.

```bash
$ joe java    # saves or appends the .gitignore file for java language
```

### Overwrite existing file

```bash
$ joe -o java    # saves a new .gitignore file for java language
```

### Multiple languages

```bash
$ joe java,javascript,python    # saves a new .gitignore file for multiple languages
```

### List all available files

```bash
$ joe ls

# OR

$ joe list
```

Output:

```bash
Output goes here
```

### View the output

```bash
$ joe -v java    # prints the .gitignore file for java language
```