# mancala

the game of mancala (kalah variant) in your terminal! i don't think it's unbeatable,
but it is hard to beat (for me at least)

i wanted to see if i could come up with a heuristic to win lots of mancala games
without creating a database of perfect play, and i think i came up with a pretty
good one

## installation:
clone this directory, then inside the directory, run 
```
go install .
```
the executable will be installed to your `GOPATH`

## gameplay: 
the board looks like this: 
```
        4	4	4	4	4	4
cpu: 0 |-	-	-	-	-	-| player: 0
        4	4	4	4	4	4	
>> 
```

the player's cups are on the bottom row. from left to right, they are numbered 6 to 1
```
        4	4	4	4	4	4
cpu: 0 |-	-	-	-	-	-| player: 0
        4	4	4	4	4	4	
//      6	5	4	3	2	1	
>>
```

to move a piece, enter the number corresponding to the cup you wish to move.

```
        4	4	4	4	4	4
cpu: 0 |-	-	-	-	-	-| player: 0
        4	4	4	4	4	4	
>> 4
```

because this player entered `4`, the fourth cup from the player's store cup will be
moved.

```
        4	4	4	4	4	4
cpu: 0 |-	-	-	-	-	-| player: 1
        4	4	0	5	5	5	
>> 
```

because this player's last stone landed in their store cup, they get to play again.

```
        4	4	4	4	4	4
cpu: 0 |-	-	-	-	-	-| player: 1
        4	4	0	5	5	5	
>> 2
```

```
        4	4	5	5	5	5
cpu: 0 |-	-	-	-	-	-| player: 2
        4	4	0	5	5	0	
>> 
```

since they player's last stone did not land in their store cup, their turn is over.

now the cpu plays. and it's the player's turn again.
```
        5	5	0	5	5	5
cpu: 1 |-	-	-	-	-	-| player: 2
        5	4	1	5	5	0	
>> 5
```

since the player's last stone lands in an empty cup, and the opposite side of the
board has pieces, the player captures the pieces on both sides of the board, and
gets to go again.
```
        5	5	0	5	5	0
cpu: 1 |-	-	-	-	-	-| player: 8
        5	0	2	6	6	0	
>> 
```

gameplay continues until one side of the board is empty. the player with the most
stones in their store cup wins!
