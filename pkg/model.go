package rewards

import (
    "errors"
    "time"
    "github.com/google/btree"
)

type Deduction struct {
    Item btree.Item
    Deducted int64
}

type Reward struct {
    TimeStamp time.Time `json:"timestamp"`
    Points int64 `json:"points"`
    Payer string `json:"payer"`
    Used int64 `json:"-"`
}

func (r Reward) Less(than btree.Item) bool {
    return r.TimeStamp.Unix() < than.(Reward).TimeStamp.Unix()
}

type RewardStore struct {
    Rewards *btree.BTree
}

func NewRewardStore() *RewardStore {
    return &RewardStore{
        Rewards: btree.New(2),
    }
}

func (s *RewardStore) AddReward(timeStamp string, points int64, payer string) error {
    ts, err := time.Parse(time.RFC3339, timeStamp)
    if err != nil {
        return errors.New("Invalid timestamp")
    }

    if points < 0  {
        balances := s.CheckBalance()

        if payerBalance, ok := balances[payer]; !ok || payerBalance + points < 0 {
            return errors.New("Insufficient points to apply transaction")
        }
    }

    s.Rewards.ReplaceOrInsert(Reward{
        TimeStamp: ts,
        Points: points,
        Payer: payer,
        Used: 0,
    })

    return nil
}

func (s *RewardStore) CheckBalance() map[string]int64 {
    balances := make(map[string]int64, 0)

    // Ascend tree in increasing time order and sum Points per Payer
    s.Rewards.Ascend(func (i btree.Item) bool {
        r := i.(Reward)
        if _, ok := balances[r.Payer]; ok {
            balances[r.Payer] += r.Points - r.Used
        } else {
            balances[r.Payer] = r.Points - r.Used
        }

        return true
    })

    return balances
}

func (s *RewardStore) UsePoints(requested int64) (map[string]int64, error) {
    remaining := requested
    totals := make(map[string]int64, 0)
    var deductions []Deduction

    s.Rewards.Ascend(func (i btree.Item) bool {
        r := i.(Reward)
        var d Deduction
        if r.Points == r.Used {
            return true
        } else if r.Points - r.Used > remaining {
            d = Deduction{
                Item: i,
                Deducted: remaining,
            }
            remaining = 0
        } else {
            d = Deduction{
                Item: i,
                Deducted: r.Points - r.Used,
            }
            remaining -= (r.Points - r.Used)
        }

        deductions = append(deductions, d)
        if remaining == 0 {
            return false
        }

        return true
    })

    if remaining > 0 {
        return nil, errors.New("Not enough points")
    }

    for _, d := range deductions {
        r := d.Item.(Reward)
        r.Used += d.Deducted
        s.Rewards.ReplaceOrInsert(r)

        if _, ok := totals[r.Payer]; ok {
            totals[r.Payer] -= int64(d.Deducted)
        } else {
            totals[r.Payer] = -1 * int64(d.Deducted)
        }
    }

    return totals, nil
}
