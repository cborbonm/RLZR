package ibin

import (
    "time"
    //"context"
	"runtime/pprof"
	"sync"
	"os"
	"log"
	//"github.com/stanford-esrg/lzr"
	"rlzr"
	"fmt"
)


func LZRMain() {
    // create a context that can be cancelled
    //ctx, cancel := context.WithCancel(context.Background())

	start := time.Now()

    //read in config 
    options, ok := rlzr.Parse()
	if !ok {
		fmt.Fprintln(os.Stderr,"Failed to parse command line options, exiting.")
		return
	}

	//For CPUProfiling
	if options.CPUProfile != "" {
		f, err := os.Create(options.CPUProfile)
		if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
	}

	//initalize
	ipMeta := rlzr.ConstructPacketStateMap( options )
    f := rlzr.InitFile( options.Filename )
	rlzr.InitParams()

    writingQueue := rlzr.ConstructWritingQueue( options.Workers )
    pcapIncoming := rlzr.ConstructPcapRoutine( options.Workers )
	timeoutQueue := rlzr.ConstructTimeoutQueue( options.Workers )
    retransmitQueue := rlzr.ConstructRetransmitQueue( options.Workers )
    timeoutIncoming := rlzr.PollTimeoutRoutine(
        &ipMeta,timeoutQueue, retransmitQueue, options.Workers, options.Timeout, options.RetransmitSec )
	incoming := rlzr.ConstructIncomingRoutine( options.Workers )
	var incomingDone sync.WaitGroup
	incomingDone.Add(options.Workers)
    done := false
	writing := false


    // record to file
    go func() {
        for {
            select {
                case input := <-writingQueue:
					writing = true
                    f.Record( input, options.Handshakes )
					writing = false
                }
        }
    }()
    //start all workers

    //read from zmap

	for i := 0; i < options.Workers; i ++ {
        go func( i int ) {
	        for input := range incoming {
				if rlzr.ReadZMap() {
					toACK := true
					toPUSH := false
					rlzr.SendAck( options, input, &ipMeta, timeoutQueue,
						retransmitQueue, writingQueue, toACK, toPUSH, rlzr.ACK)
				} else {
					 rlzr.SendSyn( input, &ipMeta, timeoutQueue )
				}
				ipMeta.FinishProcessing( input )
            }
            //ExitCondition: incoming channel closed
			if (i == options.Workers - 1) {
				for {
					if ipMeta.IsEmpty() {
						done=true
						break
					}
					//slow down to prevent CPU busy looping
					time.Sleep(1*time.Second)
					fmt.Fprintln(os.Stderr,"Finishing Last:", ipMeta.Count())
				}
			}
			incomingDone.Done()
			return
        }(i)
	}


    //read from pcap
    for i := 0; i < options.Workers; i ++ {
        go func( i int ) {
            for input := range pcapIncoming {
						//fmt.Println("pcap incoming")
						//fmt.Println(input)
                        inMap, startProcessing := ipMeta.IsStartProcessing( input )
                        //if not in map, return
                        if !inMap {
                            continue
                        }
                        //if another thread is processing, put input back
                        if !startProcessing {
                            pcapIncoming <- input
                            continue
                        }
				        rlzr.HandlePcap(options, input, &ipMeta, timeoutQueue,
							retransmitQueue, writingQueue )
                        ipMeta.FinishProcessing( input )
						//fmt.Println("finished pcap:")
						//fmt.Println(input)
            }
        }(i)
    }

    //read from timeout
    go func() {

        for input := range timeoutIncoming {
                    inMap, startProcessing := ipMeta.IsStartProcessing( input )
                    //if another thread is processing, put input back
                    //if not in map, return
                    if !inMap {
                        continue
                    }
                    if !startProcessing {
                        timeoutIncoming <- input
                        continue
                    }
                    rlzr.HandleTimeout( options, input, &ipMeta, timeoutQueue, retransmitQueue, writingQueue )
                    ipMeta.FinishProcessing( input )
		    }
    }()

    //exit gracefully when done
	incomingDone.Wait()

    for {
       if done && len(writingQueue) == 0 && !writing {
				if options.MemProfile != "" {
					f, err := os.Create(options.MemProfile)
					if err != nil {
						log.Fatal(err)
					}
					pprof.WriteHeapProfile(f)
					f.Close()
				}
			//closing file
			f.F.Flush()
			t := time.Now()
			elapsed := t.Sub(start)
			rlzr.Summarize( elapsed )
            return
       }
    }



} //end of main
