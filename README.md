## Counter Queue
Imitate processing task with workers in graceful-way, means processed task must be completed before allowing shut down 

### Build
```
docker build -t counter-queue .
```

### Run 
```
docker run -it counter-queue bin/counter-queue
```