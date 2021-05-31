#### fcache


#### Elimination strategy
data was saved in memory , because memories is limit, when memory was run out of,
how do we kick of data in memory?

1. FIFO(first in first out) 
2. LFU(least frequently used) 
   > LFU algorithm need a field to record data visited times
3. LRU(least recently used)
   > when data was visited, move it to front, when memory is full, kick out tail data, cuz
   > it was least recently used



