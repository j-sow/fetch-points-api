package rewards

import (
    "time"
    "testing"
    
    "github.com/google/btree"
)

func parseTimeStamp(timeStamp string) time.Time {
    ts, _ := time.Parse(time.RFC3339, timeStamp)
    return ts
}

func TestAdd(t *testing.T) {
    store := NewRewardStore()

    rewards := []struct{
        ts string
        points int64
        payer string
    } {
        {"2020-11-02T14:00:00Z", 1000, "DANNON"},
        {"2020-10-31T11:00:00Z", 200, "UNILEVER"},
        {"2020-10-31T15:00:00Z", -200, "DANNON"},
        {"2020-11-01T14:00:00Z", 10000, "MILLER COORS"},
        {"2020-10-31T10:00:00Z", 300, "DANNON"},
    }

    for _, reward := range rewards {
        store.AddReward(reward.ts, reward.points, reward.payer)
    }

    var last int64 = 0
    store.Rewards.Ascend(func (i btree.Item) bool {
        t.Logf("%v", i.(Reward))
        now := i.(Reward).TimeStamp.Unix()
        if now < last {
            t.Errorf("TimeStamps not in order! %s", time.Time(i.(Reward).TimeStamp).Format(time.RFC3339))
            return false
        }

        last = now
        return true
    })
}

func TestBalance(t *testing.T) {
    store := NewRewardStore()

    for i := 0; i < 50; i++ {
        direction := 1
        payer := "PAYER_A"
        if i % 2 == 1 {
            direction = -1
            payer = "PAYER_B"
        }
        ts := time.Now().Add(time.Duration(direction * i) * time.Second) 
        store.AddReward(ts.Format(time.RFC3339), 100, payer)
    } 

    for payer, points := range store.CheckBalance() {
        if points != 2500 {
            t.Errorf("%d points from %s; expected 2500", points, payer)
        }
    }
}

func TestUse(t *testing.T) {
    store := NewRewardStore()

    rewards := []struct{
        ts string
        points int64
        payer string
    } {
        {"2020-11-02T14:00:00Z", 1000, "DANNON"},
        {"2020-10-31T11:00:00Z", 200, "UNILEVER"},
        {"2020-10-31T15:00:00Z", -200, "DANNON"},
        {"2020-11-01T14:00:00Z", 10000, "MILLER COORS"},
        {"2020-10-31T10:00:00Z", 300, "DANNON"},
    }

    for _, reward := range rewards {
        store.AddReward(reward.ts, reward.points, reward.payer)
    }

    deductions, err := store.UsePoints(5000)
    if err != nil {
        t.Errorf("Did not expect error; got %s", err)
    }

    numPayers := 0
    for payer, points := range deductions {
        if payer == "DANNON" {
            numPayers += 1
            if points != -100 {
                t.Errorf("Expected deduction of 100 points from DANNON: got %d", points)
            }
        } else if payer == "UNILEVER" {
            numPayers += 1
            if points != -200 {
                t.Errorf("Expected deduction of 200 points from UNILEVER: got %d", points)
            }
        } else if payer == "MILLER COORS" {
            numPayers += 1
            if points != -4700 {
                t.Errorf("Expected deduction of 200 points from MILLER COORS: got %d", points)
            }
        }
    }

    if numPayers != 3 {
        t.Errorf("Expected number of payers to be 3: got %d", numPayers)
    }

    balances := store.CheckBalance()
    for payer, points := range balances {
        t.Logf("%s: %d", payer, points)
    }
}
