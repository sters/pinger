# pinger

Ping somethings for use warmup and standby or others.

```go
ctx := context.Background()

fooPinger := pinger.NewWorker(...)
barPinger := pinger.NewWorker(...)

eg, ctx := errgroup.WithContext(ctx)
eg.Go(func() error { return fooPinger.Run(ctx) })
eg.Go(func() error { return barPinger.Run(ctx) })

if err := eg.Wait(); err != nil {
    ...
}
```
