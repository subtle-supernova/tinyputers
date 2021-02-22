package cmd

import (
  "log"
  "fmt"
  "time"
  "github.com/schollz/peerdiscovery"
  "github.com/spf13/cobra"
  bolt "go.etcd.io/bbolt"
)

var (
  name string
  interval int
  cycles int
  dbFile string

  discoverCmd = &cobra.Command{
    Use: "discover",
    Short: "discover other machines",
    Long: "Discover other machines and slap em in redis I think",
    Run: discoverLoop,
  }
)

func init() {
  discoverCmd.PersistentFlags().StringVar(&name, "name", "UNKNOWN", "The name of the machine")
  discoverCmd.PersistentFlags().IntVar(&interval, "interval", 5, "Interval to scan")
  discoverCmd.PersistentFlags().IntVar(&cycles, "cycles", 5, "Number of searches. -1 for keep searching.")
  discoverCmd.PersistentFlags().StringVar(&dbFile, "db", "hosts", "filename for db store")
}

func discoverLoop(cmd *cobra.Command, args []string) {
  fmt.Println("scanning")

  db, err := bolt.Open(dbFile + ".db", 0600, nil)
  if err != nil {
    log.Panicf("Could not open %s for db %v", dbFile, err)
  }
  defer db.Close()

  for {
    discoveries, err := peerdiscovery.Discover(peerdiscovery.Settings{
      Limit: cycles,
      Payload: []byte(name),
    })

    if err != nil {
      log.Printf("Error discovering: %v", err)
      time.Sleep(time.Duration(interval * 2) * time.Second)
    } else {
      if len(discoveries) > 0 {
        fmt.Println("Found some stuff: ")
        err := db.Update(func(tx *bolt.Tx) error {
          b, err := tx.CreateBucketIfNotExists([]byte("tinyputers"))
          if err != nil {
            return err
          }
          for i, d := range discoveries {
            fmt.Printf("%d) '%s' with payload '%s'\n", i, d.Address, d.Payload)
            err = b.Put([]byte(d.Address), []byte(d.Payload))
            if err != nil {
              return err
            }
          }
          return nil
        })
        if err != nil {
          log.Printf("Error writing to bolt: %v", err)
        }
      } else {
        log.Println("We didn't find anything")
      }
      time.Sleep(time.Duration(interval) * time.Second)
    }
  }
}

func Execute() error {
  return discoverCmd.Execute()
}
