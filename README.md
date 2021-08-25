# SAMUrig
A fourth-order Runge-Kutta simulation program to find rigidity of particles passing through SAMURAI magnet

Before the first run, modify the magnetic field map files to be used. main.go: line 211. The program will generate a "bmap.bin" binary field map and load it in next runs. Remove bmap.bin if you want a different magnetic field. The default is 1.44 T (1.40+(1.45-1.40) * (400/500)). Find the field map from [SAMURAI magnet page](https://ribf.riken.jp/SAMURAI/index.php?Magnet).

Modify rig/work.go if your input file has a different format.

Modify the main.go file to change the distances and angles in your setup. The default values are shown in the figure:

![image](fig.png)
