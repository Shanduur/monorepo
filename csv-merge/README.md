# CSV-Merge

## Data

Main *data* folder must be placed in the same directory as the script. It's structure should look like this:

```
./data
├── F0
│   ├── MF
│   │   ├── p0.csv
│   │   └── p1.csv
│   └── UX
│       ├── p0.csv
│       └── p1.csv
└── F1
    ├── MF
    │   ├── p0.csv
    │   └── p1.csv
    └── UX
        ├── p0.csv
        └── p1.csv
```

It will result in creation of files *F0_MF.xlsx*, *F0_UX.xlsx*, *F1_MF.xlsx* and *F1_UX.xlsx*, each one with two worksheets named *p0* and *p1*