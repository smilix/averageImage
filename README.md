# averageImage

Calculates one average image for other images.

This go program analyses all given images and calculates for every pixel and every single RGB value the arithmetical average.

# Usage

```sh
./averageImage                                                                                                                                                                Usage of ./averageImage:
  -filePattern string
        A reg exp pattern that must match the file. (default "(?i)\\.jpe?g$")
  -folder string
        A folder with images
  -height int
        the expected height
  -maxImages int
        max amount of images to read, -1 for no max value (default -1)
  -output string
        Output image (is always a jpg)
  -width int
        the expected with
```


## Examples

```sh
./averageImage -folder /some/folder -maxImages 20 -output out1.jpg -width 4592 -height 3448 -filePattern '(?i)^_.*\.jpe?g$'
```
Maximum 20 images, file name must start with `_` and end with `.jpg` or `.jpeg` (ignore case)

```sh
./averageImage -folder /Users/holger/tmp/test -maxImages 5 -output out2.jpg -width 4592 -height 3448 -filePattern '(?i)\.jpe?g$'
```
Maximum 5 images, file name must end with `.jpg` or `.jpeg` (ignore case)

Or see `averageResult_sample.jpg`


# Misc

Attention, this is my first I-learn-go-project.
