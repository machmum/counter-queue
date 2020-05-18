package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/machmum/counter-queue/counter"
)

// instruction below
var _ = `
Bayangkan sebuah sistem antrian yang terdiri dari sejumlah loket.
Kemudian terdapat sejumlah orang yang mengantri di loket untuk dilayani.
Masing-masing loket malayani satu orang, dengan lama waktu layanan yang berbeda-beda.
Ketika tutup, maka loket tidak boleh mengambil antrian berikutnya, namun harus tetap menyelesaikan layanan yang sedang berjalan sebelum tutup.

Dari ilustrasi di atas, buatlah simulasi aplikasi dengan ketentuan berikut:
1. Jumlah loket yang buka bisa dipilih (1 - 5)
2. Antrian terdiri dari 10 orang dengan simulasi waktu pelayanan untuk masing-masing orang adalah sebagai berikut (dalam detik): 
   [1, 2, 4, 2, 3, 5, 2, 3, 1, 3]
3. Proses melakukan pelayanan dapat diwakili dengan time.Sleep
4. Ketika aplikasi dihentikan dengan Ctrl + C atau termination signal, maka lakukan simulasi loket tutup untuk semua loket.
5. Code disimpan di public github/gitlab repository dengan incremental commit. Kirimkan link github/gitlab repositoynya melalui email setelah selesai.

Optional:
+ Buat ke dalam dockerfile
`

var numOfCounters int

// scanCounter: get number of counter to work with
func scanCounter() {
	fmt.Println("Please input number of counter [1-5]:")
	var counterNum int
	if _, err := fmt.Scanf("%d", &counterNum); err != nil {
		log.Fatalf("Fail scan input: %v", err.Error())
	}
	numOfCounters = counterNum
}

func main() {
	for {
		scanCounter()
		if numOfCounters <= 5 {
			break
		}
	}

	// initialize counter
	c := counter.NewCounter()

	// start with simulation time
	simulationTime := []int{1, 2, 4, 2, 3, 5, 2, 3, 1, 3}

	go c.Start(simulationTime)

	// set up cancellation context and wait group
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// start counter with cancellation context passed
	go c.StartProcess(ctx)

	// start workers and add numberOfWorkers to waitGroup
	wg.Add(numOfCounters)

	// start [numOfCounters] workers
	for i := 1; i < numOfCounters+1; i++ {
		go c.LoketJobs(wg, i)
	}

	// handle sigterm and await termChan signal
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan // blocks here until either SIGINT or SIGTERM is received

	// handle shutdown
	fmt.Println("***\nshutdown signal received\n***")

	// signal cancellation to context.Context,
	// call the cancelFunc to notify the counter it's time to shut down
	cancelFunc()

	wg.Wait() // block here until worker in counter.counterJobs are done

	fmt.Println("counter done processing, shutting down!")

}
