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
}

func (r Reward) Less(than btree.Item) bool {
    return r.TimeStamp.Unix() < than.(Reward).TimeStamp.Unix()
}

type RewardStore struct {
    Rewards *btree.BTree
    Balances map[string]int64
    UnusedDeductions map[string]int64
}

func NewRewardStore() *RewardStore {
    return &RewardStore{
        Rewards: btree.New(2),
        Balances: make(map[string]int64, 0),
        UnusedDeductions: make(map[string]int64, 0),
    }
}

func (s *RewardStore) AddReward(timeStamp string, points int64, payer string) error {
    ts, err := time.Parse(time.RFC3339, timeStamp)
    if err != nil {
        return errors.New("Invalid timestamp")
    }

    if points < 0 {
        if _, ok := s.UnusedDeductions[payer]; !ok {
            s.UnusedDeductions[payer] = -1 * points
        } else {
            s.UnusedDeductions[payer] -= points
        }
    } else {
        s.Rewards.ReplaceOrInsert(Reward{
            TimeStamp: ts,
            Points: points,
            Payer: payer,
        })
    }

    if _, ok := s.Balances[payer]; !ok {
        s.Balances[payer] = points
    } else {
        s.Balances[payer] += points
    }
    

    return nil
}

func (s *RewardStore) CheckBalance() map[string]int64 {
    return s.Balances
}

func (s *RewardStore) UsePoints(requested int64) (map[string]int64, error) {
    remaining := requested
    totals := make(map[string]int64, 0)
    var deductions []Deduction
    
    // Realize unused deductions
    s.Rewards.Ascend(func (i btree.Item) bool {
        r := i.(Reward)
        if unused, ok := s.UnusedDeductions[r.Payer]; ok {
            var d Deduction
            if r.Points > unused {
                d = Deduction{
                    Item: i,
                    Deducted: unused,
                }
                delete(s.UnusedDeductions, r.Payer)
            } else {
                d = Deduction{
                    Item: i,
                    Deducted: r.Points,
                }
                s.UnusedDeductions[r.Payer] -= r.Points
            }

            deductions = append(deductions, d)
        }

        if len(s.UnusedDeductions) == 0 {
            return false
        }

        return true
    })

    for _, d := range deductions {
        r := d.Item.(Reward)
        r.Points -= d.Deducted
        if r.Points == 0 {
            s.Rewards.Delete(r)
        } else {
            s.Rewards.ReplaceOrInsert(r)
        }
    }

    deductions = nil
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
        r.Points -= d.Deducted
        if r.Points == 0 {
            s.Rewards.Delete(r)
        } else {
            s.Rewards.ReplaceOrInsert(r)
        }

        if _, ok := totals[r.Payer]; ok {
            totals[r.Payer] -= d.Deducted
        } else {
            totals[r.Payer] = -1 * d.Deducted
        }

        s.Balances[r.Payer] -= d.Deducted
    }

    return totals, nil
}
