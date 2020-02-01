# greedy-guided-astar

## Purpose
This is a reasonably unique implementation of the A\* algorithm. A greedy
algorithm is used to find an initial solution, then A\* is 
used on successive segments to find a moderately optimal solution.

Golang was chosen as it has inherent parallel processing, and I've needed
an excuse to write go code.

### The Current Problem
The example problem used here is the NPuzzle. This is a sliding puzzle
as seen in your childhood, with n rows of n numbered tiles (missing one).

Pieces are slid into the empty tile to attempt to organize the numbers.
```
[214]    [123]
[673] -> [456]  
[_58]    [78_]
```

My original Python implementation of A\* could find a guaranteed optimal
solution for a 3x3 in a matter of seconds, but a 4x4 would fail after
several hours (out of memory).

The Python implementation of the greedy guided A\* could solve a 4x4 in 
around 3 minutes and find a reasonable solution.  Golang should be much quicker.

And... I've implemented simple A\* in golang, and it can solve a 4x4 
in around 15 seconds with a shuffle of 1000 and a Seed of 1.  Other seeds 
crash with the expected memory overrun.

```
StartState
[6 1 13 10]
[15 0 4 9]
[2 14 7 3]
[12 11 5 8]

Goal
[1 2 3 4]
[5 6 7 8]
[9 10 11 12]
[13 14 15 0]

Found solution in 21.32594146s time.
Found solution in 53 steps

Found solution in 23.866670582s time.
Found solution in 53 steps

```


### Abstraction
The SequentialInterface interface gives the functions needed by the solver.
Any sequential problem should have these functions defined, and should allow
this method to be used.