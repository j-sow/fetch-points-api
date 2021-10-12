package rewards

import (
    "time"
    "testing"
    
    "github.com/google/btree"
)

func TestAdd(t *testing.T) {
    store := NewRewardStore()

    for i := 0; i < 50; i++ {
        direction := 1
        if i % 2 == 1 {
            direction = -1
        }
        ts := time.Now().Add(time.Duration(direction * i) * time.Second) 
        store.AddReward(ts.Format(time.RFC3339), 100, "STEVE_DOT_COM")
    } 

    var last int64 = 0
    store.Rewards.Ascend(func (i btree.Item) bool {
        now := i.(Reward).TimeStamp.Unix()
        if now < last {
            t.Errorf("TimeStamps not in order! %s", i.(Reward).TimeStamp.Format(time.RFC3339))
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

    ts1 := time.Now()
    ts2 := ts1.Add(time.Duration(10) * time.Second) 
    ts3 := ts1.Add(time.Duration(20) * time.Second)
    store.AddReward(ts1.Format(time.RFC3339), 100, "PAYER_A")
    store.AddReward(ts2.Format(time.RFC3339), 100, "PAYER_B")
    store.AddReward(ts3.Format(time.RFC3339), 100, "PAYER_A")

    deductions, err := store.UsePoints(150)
    if err != nil {
        t.Errorf("Did not expect error; got %s", err)
    }

    numPayers := 0
    for payer, points := range deductions {
        if payer == "PAYER_A" {
            numPayers += 1
            if points != 100 {
                t.Errorf("Expected deduction of 100 points from PAYER_A: got %d", points)
            }
        } else if payer == "PAYER_B" {
            numPayers += 1
            if points != 50 {
                t.Errorf("Expected deduction of 100 points from PAYER_B: got %d", points)
            }
        }
    }
    rewardsRemaining := store.Rewards.Len()
    if rewardsRemaining != 2 {
        t.Errorf("Expected 2 rewards remaining; got %d", rewardsRemaining)
    }

    if numPayers != 2 {
        t.Errorf("Expected number of payers to be 2: got %d", numPayers)
    }

    _, err = store.UsePoints(3000)
    if err == nil {
        t.Errorf("Exepected error; got nil")
    }

    rewardsRemaining = store.Rewards.Len()
    if rewardsRemaining != 2 {
        t.Errorf("Expected 2 rewards remaining; got %d", rewardsRemaining)
    }
}
