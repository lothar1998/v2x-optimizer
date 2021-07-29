int V = ...;
int N = ...; 

range vRange = 1..V;
range nRange = 1..N;

float R[vRange][nRange] = ...;
int MRB[nRange] = ...;

dvar boolean x[nRange];
dvar boolean y[vRange][nRange];

minimize sum(n in nRange) x[n];

subject to {
  forall(n in nRange)
    MRB[n] >= sum(v in vRange) y[v][n] * R[v][n];
    
  forall(v in vRange)
    sum(n in nRange) y[v][n] == 1;
    
  forall(v in vRange)
    forall(n in nRange)
      x[n] >= y[v][n];
}

execute {
  writeln("RRH = ", x)
}
