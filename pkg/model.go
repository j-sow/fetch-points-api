package rewards

import (
    "errors"
    "time"
    "github.com/google/btree"
)

type Deduction struct {
    Item btree.Item
    Deducted uint32
}

type Reward struct {
    TimeStamp time.Time
    Points uint32
    Payer string
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

func (s *RewardStore) AddReward(timeStamp string, points uint32, payer string) error {
    ts, err := time.Parse(time.RFC3339, timeStamp)
    if err != nil {
        return errors.New("Invalid timestamp")
    }

    s.Rewards.ReplaceOrInsert(Reward{
        TimeStamp: ts,
        Points: points,
        Payer: payer,
    })

    return nil
}

func (s *RewardStore) CheckBalance() map[string]uint32 {
    balances := make(map[string]uint32, 0)

    // Ascend tree in increasing time order and sum Points per Payer
    s.Rewards.Ascend(func (i btree.Item) bool {
        r := i.(Reward)
        if _, ok := balances[r.Payer]; ok {
            balances[r.Payer] += r.Points
        } else {
            balances[r.Payer] = r.Points
        }

        return true
    })

    return balances
}

func (s *RewardStore) UsePoints(requested uint32) (map[string]uint32, error) {
    remaining := requested
    totals := make(map[string]uint32, 0)
    var deductions []Deduction

    s.Rewards.Ascend(func (i btree.Item) bool {
        r := i.(Reward)
        var d Deduction
        if r.Points > remaining {
            d = Deduction{
                Item: i,
                Deducted: remaining,
            }
            remaining = 0
        } else {
            d = Deduction{
                Item: i,
                Deducted: r.Points,
            }
            remaining -= r.Points
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
        if d.Deducted != r.Points{
            r.Points -= d.Deducted
            s.Rewards.ReplaceOrInsert(r)
        } else {
            s.Rewards.Delete(d.Item)
        }

        if _, ok := totals[r.Payer]; ok {
            totals[r.Payer] += d.Deducted
        } else {
            totals[r.Payer] = d.Deducted
        }
    }

    return totals, nil
}
