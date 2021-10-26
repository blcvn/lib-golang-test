# CPU Profiling
go test -run=none -bench=ClientServerParallel4 -cpuprofile=cprof net/http
go tool pprof --text http.test cprof

# Memory profiling 

 go test --memprofile 
  with net/http/pprof via http://myserver:6060:/debug/pprof/heap 
  
  or by calling runtime/pprof.WriteHeapProfile.


# Blocking Profiler

 'go test --blockprofile', with net/http/pprof via http://myserver:6060:/debug/pprof/block 
 or by calling runtime/pprof.Lookup("block").WriteTo.

#Goroutine Profiler

collect the profile with net/http/pprof via http://myserver:6060:/debug/pprof/goroutine, 
and visualize it to svg/pdf or by calling runtime/pprof.Lookup("goroutine").WriteTo.

#Garbage Collector Trace

 GODEBUG=gctrace=1 ./myserver

 #Scheduler Trace
GODEBUG=schedtrace=1000


net/http/pprof at the bottom of  http://myserver:6060/debug/pprof/heap?debug=1
