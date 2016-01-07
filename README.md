shapef
======
shapef is a simple program for converting shape definitions from a
JSON file to an [ESRI Shapefile](https://www.esri.com/library/whitepapers/pdfs/shapefile.pdf).
It uses the [go-shp module](https://github.com/jonas-p/go-shp) and is currently
just used for generating test files for [openLCA](http://openlca.org).

It has the following limitations:
* It only generates single-line polygon features
* Only numbers and strings are allowed as attribute values of features
* No validation of user input and no user friendly error handling
* ...

(Feel free to send pull-requists to improve this)

### Usage
The shapef executable expects a JSON file as input:

    shapef <file prefix>.json

This will generate the following files if everything is correct:
* `<file prefix>.shp` the main file of the shapefile
* `<file prefix>.shx` the index file of the features in the shapefile
* `<file prefix>.dbf` the dBASE file with the attribute data of the features

### Input format
The JSON format of the input file is as follows:
* The file should contain an array of feature objects.
* The coordinates of the features are stored as array of points in the
  `points` field.
* Each point is a simple array with two numbers: the x- and y-coordinate.
* The attributes of the features are stored as key-value pairs in the `data`
  field.

```json
[
    {
        "points": [[0, 0], [10, 0], [10, 10], [0, 10]],
        "data": {
            "label": "Shape 1",
            "value": 42,
            "code": "AB"
        }
    },
    {
        "points": [[10, 0], [20, 0], [20, 10], [10, 10]],
        "data": {
            "label": "Shape 2",
            "value": 24,
            "code": "CD"
        }
    }
]
```
