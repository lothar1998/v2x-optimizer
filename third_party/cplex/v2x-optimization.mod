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
  writeln("N = ", N);
  writeln("V = ", V);
  writeln("RRH_COUNT = ", cplex.getObjValue());
  write("RRH_ENABLE = [");
  if (N > 0) {
    write(x[1]);
    for (var n=2; n <= N; n++) {
      write(" ", x[n]);
    }
  }
  writeln("]");
  write("VEHICLE_ASSIGNMENT = [");
  for (var n=1; n <= N; n++) {
    if (y[1][n] == 1) {
      write(n - 1);
      break;
    }
  }
  for (var v=2; v <= V; v++) {
    for (var n=1; n <= N; n++) {
      if (y[v][n] == 1) {
        write(" ", n - 1);
        break;
      }
    }
  }
  writeln("]");
}
