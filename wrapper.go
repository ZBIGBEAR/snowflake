package snowflake

func GetID() (int64, error) {
	if snowFlake == nil {
		initSnowFlake(0)
	}

	return snowFlake.GetID()
}

func GetIDs(count int) ([]int64, error) {
	if snowFlake == nil {
		initSnowFlake(0)
	}

	if count <= NB4095 {
		return snowFlake.GetIDs(count)
	}

	split := count / NB4095
	left := count % NB4095

	var ids []int64
	for i := 0; i < split; i++ {
		splitIDs, err := snowFlake.GetIDs(NB4095)
		if err != nil {
			return nil, err
		}

		ids = append(ids, splitIDs...)
	}

	if left > 0 {
		leftIDs, err := snowFlake.GetIDs(left)
		if err != nil {
			return nil, err
		}

		ids = append(ids, leftIDs...)
	}

	return ids, nil
}
