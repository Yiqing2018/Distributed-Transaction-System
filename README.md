# Distributed Transaction System

## Overview
This is a [course project](https://courses.engr.illinois.edu/ece428/sp2019//mps/mp3.html) of CS425 in UIUC.  

Implemented a distributed transaction system, supporting transactions that read and write to distributed objects while ensuring full ACI properties.  

## Running Instruction

step1. run coordinator

```
$./coordinator
```
step2. run servers A.B.C.D.E

```
$./server
```
step3. run client(s)

```
$./client
```
step4. GET/SET operations

```
$ GET A.x
$ SET A.x 10
```
the Addresses of five servers are hard-codes in client program, you may want to modify that :)

***

## Design

### Two-Phrase Commit
![](http://ww2.sinaimg.cn/large/006tNc79ly1g5fq37x55ej30n90cwacb.jpg)

At the very beginning, client would send **"openTransaction"** to the coordinator, and once received TID(transaction ID) message from coordinator, it would print out **"OK"** on the screen, which means ready to compute transaction!  

While client receives operation command from user input, like  **"SET A.X 999"**, it would send the operation msg along with TID(transaction ID) to every server.

![](http://ww1.sinaimg.cn/large/006tNc79ly1g5frb038c5j30l00ewtbh.jpg)

**A.Commit Operation**  

**Two-Phrase Commit**  is implemented in this way:  

(1) Client send **COMMIT** to five servers and the coordinator  

(2) Once coordinator received COMMIT from client, coordinator send **canCommit** to five servers.  

(3) Unitl coordinator received 5 YES, coordinator would send **doCommit** to five servers and send **OK** to client together, else send **doAbort** to five servers and send **ABORTED** to client.

**B.Abort Operation**
When Client sends **ABORT** to the coordinator, or coordinator detects **DeadLock**, coordinator will send **doAbort** to five servers directly.

### Concurrency and isolation
Client may execute transactions concurrently, in our project, the **serializability** is guaranteed with **strict two Phase Locking protocol** (strict 2PL).  

There are two types of Locks: **Shared Lock(S Lock) & Exclusive Lock(X Lock)**, S lock is required for **GET** operation, and X lock is required for **SET** operation.  
All locks within one transaction would be released until the transaction is done.

In our implementation, we use **map structure** to store **Lock and holder** information on every server.

```
var serverLock = struct {
	sync.RWMutex
	//key:=object, value:=locktype
	lockInfo map[string]int
	//key:=object, value:=holder or holders
	holder map[string][]string
}{
	lockInfo: make(map[string]int),
	holder:   make(map[string][]string),
}
```

when server received **operation command**:  
**IF** it is ok to grab the required lock:  
send return value to client  
update serverLock   
**ELSE**: 
put the operation into buffer
 
when server received **doCommit(T)/doAbort(T)** command:  
**release** all the locks Transaction T is holding  
 check the **buffer**, may execute some operations that are waiting for the released lock

For example:  

<p align="center">

  <img width="300" src="http://ww3.sinaimg.cn/large/006tNc79ly1g5fqnezkwpj30ec0alwf4.jpg">
    <img width="300" src="http://ww2.sinaimg.cn/large/006tNc79ly1g5fqns6kmpj30ja08ddgk.jpg">
  
</p>

In this scenario, T2 has to wait until the XLock(B) is released by T2.
When serverB receives "SET B.X 7" message from T2, it would put this operation into buffer and send wait-for information to coordinator (for DeadLock detection)
When serverB receives "COMMIT" message from T1, it would release all the locks held by T1, and check the buffer, execute operation "SET B.X 7"

### Deadlock Detection

**Centralized Detection**  
Centralized Detection is used to prevent the deadlock in the distributed transaction system.
- A coordination plays a role in Two-Phase Commit Protocol, a reuse of coordinator will be highly efficient.
- To relieve bandwidth burden

**Server Side**  
To be specific, each server (A, B, C, D and E) will report the wait-for relationships to coordinator.  
For example, when a Transaction (TID = 1) is waiting for one object (A.x) that the Transaction (TID = 2) occupied, Server A will send **"2 → 1"**  message to coordinator, which means that Transaction 1 is waiting for Transaction 2.

**Coordinator Side**  
The coordinator will preserve a global waif-for graph. Whenever coordinator receives a wait-for message, it would add an 
**edge**, and check deadlock existence. For example, this coordinator receives wait-for relationship messages as belows:  

2 → 3, T3 is waiting for T2  
1 → 2, T2 is waiting for T1  
3 → 1, T1 is waiting for T3  

In this scenario, there is a circle in the global wait-for graph, in other words, these transactions caused a deadlock.

**Abort Strategy**  
coordinator would choose to abort the Transaction with the lowest TID. In the below figure, the transaction with color red will be aborted.
<p align="center">
    <img width="600" src="http://ww3.sinaimg.cn/large/006tNc79ly1g5fqvo2pdqj30fz09hdge.jpg">
</p>

*** 
