## gridhash
GridHash is a highly experimental password hashing tool intended to create high-strength hashes, that are resistant to high-speed cracking.

#### Why?
The need for a better way to hash passwords is clear, I created this to test some ideas I've had on how to use existing primitives to create hashes that are far stronger than what's in use today.

#### So, is this production ready?
No.

Seriously, this is an experiment, it's unlikely that I'd ever call it production ready. While it does have some security guarantees, those are weaker than some other options available today.

#### How does it work?
GridHash creates a grid of hashes, that are generated based on prior hashes. This increases both memory and CPU requirements to calculate a hash. It looks a little like this:

```
+---+------+------+------+------+------+
| X |  0   |  1   |  2   |  3   |  4   |
+---+------+------+------+------+------+
| 0 | HMAC | H    | H    | H    | H    |
| 1 | H    | HMAC | H    | H    | H    |
| 2 | H    | H    | HMAC | H    | H    |
| 3 | H    | H    | H    | HMAC | H    |
| 4 | H    | H    | H    | H    | HMAC |
+---+------+------+------+------+------+
```

Where `H` is a hash function, and `HMAC` is, big surprise, an `HMAC` function.

To get the process started, cells `0:0` and `0:1` are pre-populated, then the others are filled in, by row, from top down, then left to right. When `X == Y` the `HMAC` is used, instead on the hash function used elsewhere.

To calculate a cells value, values are appended from the bottom up, then right to left.

For example, on cell `4:4`, the value is calculated as:

```
HMAC(3:4 + 2:4 + 1:4 + 0:4 + 4:3 + 4:2 + 4:1 + 4:0 + extra_data + salt)
```

Here, `extra_data` is a block of data deterministically generated from a seeded PRNG, this size of the blocks are configurable (as is most things), so it can be tuned to the specific environment.

The idea here is to require as much data in memory as possible, and as much CPU time as possible, to make life harder for those attempting to crack these hashes.

#### Configuration
There are several ways that this can be tuned for a specific environment:

* Grid Size - This has the largest performance impact, it sets the size (X & Y) of the grid, so a grid size of 64, results in a grid with 64^2 cells.
* Hash Iterations - Each time `H` is called, it is looped the specified number of times.
* HMAC Iterations - Each time `HMAC` is called, it is looped the specified number of times.
* Extra Bytes - An array of `grid_size` elements, one element is appended (which element is based on column in the grid), prior to the salt, before the hash or HMAC operation.

#### What primitives are used?
Any combination of hash and HMAC functions may be used, the code here currently uses SHA-256 and SHA-256-HMAC, though future versions will allow selection of certain primitives from the command line.

#### Current Status
Rough. Incomplete. Requires code changes just to change the password it's hashing.

I've been focusing on the hashing process and testing with fixed settings, so currently the application isn't very usable, but the algorithm itself is in-place and working.

My next steps are to clean up the code, improve efficiency, and add a useful interface.

#### But I really want to use this in production!
Then send in some pull requests to improve the code & documentation, and get it to the point that it can be presented for a formal review.

But don't use it till the review is complete.

#### License
This software is released under the MIT license.
